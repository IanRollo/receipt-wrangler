package wranglerasynq

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
)

// fakeQueueLister serves canned tasks per queue and state, recording the options
// each list call was made with.
type fakeQueueLister struct {
	tasks     map[string]map[asynq.TaskState][]*asynq.TaskInfo
	errByQ    map[string]error
	optionCounts []int
}

func (lister *fakeQueueLister) list(state asynq.TaskState, queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error) {
	if err := lister.errByQ[queue]; err != nil {
		return nil, err
	}

	lister.optionCounts = append(lister.optionCounts, len(opts))

	return lister.tasks[queue][state], nil
}

func (lister *fakeQueueLister) ListPendingTasks(q string, o ...asynq.ListOption) ([]*asynq.TaskInfo, error) {
	return lister.list(asynq.TaskStatePending, q, o...)
}
func (lister *fakeQueueLister) ListScheduledTasks(q string, o ...asynq.ListOption) ([]*asynq.TaskInfo, error) {
	return lister.list(asynq.TaskStateScheduled, q, o...)
}
func (lister *fakeQueueLister) ListActiveTasks(q string, o ...asynq.ListOption) ([]*asynq.TaskInfo, error) {
	return lister.list(asynq.TaskStateActive, q, o...)
}
func (lister *fakeQueueLister) ListRetryTasks(q string, o ...asynq.ListOption) ([]*asynq.TaskInfo, error) {
	return lister.list(asynq.TaskStateRetry, q, o...)
}
func (lister *fakeQueueLister) ListArchivedTasks(q string, o ...asynq.ListOption) ([]*asynq.TaskInfo, error) {
	return lister.list(asynq.TaskStateArchived, q, o...)
}

func quickScanTaskInfo(t *testing.T, id string, state asynq.TaskState, payload QuickScanTaskPayload) *asynq.TaskInfo {
	t.Helper()
	bytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	return &asynq.TaskInfo{ID: id, Queue: string(models.QuickScanQueue), Type: QuickScan, Payload: bytes, State: state, MaxRetry: 3}
}

func emailTaskInfo(t *testing.T, id string, state asynq.TaskState, payload EmailProcessTaskPayload) *asynq.TaskInfo {
	t.Helper()
	bytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	return &asynq.TaskInfo{ID: id, Queue: string(models.EmailReceiptProcessingQueue), Type: EmailProcess, Payload: bytes, State: state, MaxRetry: 3}
}

func resolveByPayload(payload RerunTaskPayload) (uint, error) {
	if payload.GroupSettingsId == 0 {
		return payload.GroupId, nil
	}
	if payload.GroupSettingsId == 99 {
		return 0, errors.New("no such group settings")
	}

	// Group settings N belongs to group N * 10.
	return payload.GroupSettingsId * 10, nil
}

func TestProcessingQueueState_MapsEveryAsynqState(t *testing.T) {
	tests := []struct {
		state    asynq.TaskState
		expected string
		listed   bool
	}{
		{asynq.TaskStatePending, structs.ProcessingQueueStateQueued, true},
		{asynq.TaskStateScheduled, structs.ProcessingQueueStateQueued, true},
		{asynq.TaskStateActive, structs.ProcessingQueueStateProcessing, true},
		{asynq.TaskStateRetry, structs.ProcessingQueueStateRetrying, true},
		{asynq.TaskStateArchived, structs.ProcessingQueueStateFailed, true},
		{asynq.TaskStateCompleted, "", false},
		{asynq.TaskStateAggregating, "", false},
	}

	for _, test := range tests {
		got, ok := processingQueueState(test.state)
		if ok != test.listed || got != test.expected {
			t.Errorf("%v: got (%q, %v), want (%q, %v)", test.state, got, ok, test.expected, test.listed)
		}
	}
}

func TestBuildProcessingQueueItem_QuickScan(t *testing.T) {
	directory := containedTempDir(t)
	file := writeTempFile(t, directory, "scan.jpg", 0)
	info := quickScanTaskInfo(t, "q1", asynq.TaskStateRetry, QuickScanTaskPayload{
		Token:            &structs.Claims{UserId: 7},
		GroupId:          3,
		TempPath:         file,
		OriginalFileName: "lunch.jpg",
	})
	info.Retried = 2
	info.LastErr = "boom"
	info.LastFailedAt = time.Now().Add(-time.Minute)
	info.NextProcessAt = time.Now().Add(time.Minute)

	item, ok := buildProcessingQueueItem(info, structs.ProcessingQueueTaskTypeQuickScan, resolveByPayload)

	if !ok {
		t.Fatal("expected the task to be listed")
	}
	if item.State != structs.ProcessingQueueStateRetrying || item.TaskType != structs.ProcessingQueueTaskTypeQuickScan {
		t.Errorf("state/type wrong: %s %s", item.State, item.TaskType)
	}
	if item.GroupId != 3 || item.UploaderUserId == nil || *item.UploaderUserId != 7 || item.FileName != "lunch.jpg" {
		t.Errorf("identity fields wrong: %+v", item)
	}
	if item.RetryCount != 2 || item.MaxRetries != 3 || item.LastError != "boom" {
		t.Errorf("retry fields wrong: %+v", item)
	}
	if item.NextRetryAt == nil || item.LastFailedAt == nil || item.UploadedAt == nil {
		t.Errorf("times missing: %+v", item)
	}
	if !item.HasSourceFile || item.CanBeRestarted || !item.CanRetryNow {
		t.Errorf("flags wrong for a retrying task with its file: %+v", item)
	}
}

func TestBuildProcessingQueueItem_ArchivedWithFileCanBeRestarted(t *testing.T) {
	directory := containedTempDir(t)
	file := writeTempFile(t, directory, "scan.jpg", 0)
	info := quickScanTaskInfo(t, "q1", asynq.TaskStateArchived, QuickScanTaskPayload{GroupId: 3, TempPath: file})

	item, ok := buildProcessingQueueItem(info, structs.ProcessingQueueTaskTypeQuickScan, resolveByPayload)

	if !ok || item.State != structs.ProcessingQueueStateFailed {
		t.Fatalf("unexpected result: %v %+v", ok, item)
	}
	if !item.CanBeRestarted || item.CanRetryNow || item.NextRetryAt != nil {
		t.Errorf("flags wrong: %+v", item)
	}
	if item.UploaderUserId != nil {
		t.Error("a payload with no token has no known uploader")
	}
}

func TestBuildProcessingQueueItem_EmailUploadResolvesGroupFromSettings(t *testing.T) {
	directory := containedTempDir(t)
	file := writeTempFile(t, directory, "receipt.pdf", 0)
	info := emailTaskInfo(t, "e1", asynq.TaskStateActive, EmailProcessTaskPayload{
		GroupSettingsId: 4,
		TempFilePath:    file,
		Attachment:      structs.Attachment{Filename: "receipt.pdf"},
		Metadata:        structs.EmailMetadata{Subject: "Your receipt", FromEmail: "shop@example.com"},
	})

	item, ok := buildProcessingQueueItem(info, structs.ProcessingQueueTaskTypeEmailUpload, resolveByPayload)

	if !ok {
		t.Fatal("expected the task to be listed")
	}
	if item.GroupId != 40 || item.State != structs.ProcessingQueueStateProcessing {
		t.Errorf("group/state wrong: %+v", item)
	}
	if item.UploaderUserId != nil || item.FileName != "receipt.pdf" || item.EmailSubject != "Your receipt" || item.EmailSender != "shop@example.com" {
		t.Errorf("email fields wrong: %+v", item)
	}
}

func TestBuildProcessingQueueItem_BodyOnlyEmailHasNoSourceFileButCanBeRestarted(t *testing.T) {
	info := emailTaskInfo(t, "e1", asynq.TaskStateArchived, EmailProcessTaskPayload{GroupSettingsId: 4})

	item, ok := buildProcessingQueueItem(info, structs.ProcessingQueueTaskTypeEmailUpload, resolveByPayload)

	if !ok || item.HasSourceFile || !item.CanBeRestarted {
		t.Fatalf("a body-only email stays rerunnable without a file: %v %+v", ok, item)
	}
}

func TestBuildProcessingQueueItem_SkipsWhatItCannotAuthorize(t *testing.T) {
	undecodable := &asynq.TaskInfo{ID: "bad", Type: QuickScan, Payload: []byte("{not json"), State: asynq.TaskStateRetry}
	unresolvable := emailTaskInfo(t, "lost", asynq.TaskStateRetry, EmailProcessTaskPayload{GroupSettingsId: 99})
	noGroup := quickScanTaskInfo(t, "nogroup", asynq.TaskStateRetry, QuickScanTaskPayload{})
	foreign := &asynq.TaskInfo{ID: "other", Type: "something:else", Payload: []byte("{}"), State: asynq.TaskStateRetry}
	completed := quickScanTaskInfo(t, "done", asynq.TaskStateCompleted, QuickScanTaskPayload{GroupId: 1})

	for _, info := range []*asynq.TaskInfo{undecodable, unresolvable, noGroup, foreign, completed} {
		if _, ok := buildProcessingQueueItem(info, structs.ProcessingQueueTaskTypeQuickScan, resolveByPayload); ok {
			t.Errorf("task %s must not be listed", info.ID)
		}
	}
}

func TestListProcessingQueueItems_ReadsEveryStateOfBothQueuesNewestFirst(t *testing.T) {
	directory := containedTempDir(t)
	older := writeTempFile(t, directory, "older.jpg", 2*time.Hour)
	newer := writeTempFile(t, directory, "newer.jpg", time.Minute)

	lister := &fakeQueueLister{tasks: map[string]map[asynq.TaskState][]*asynq.TaskInfo{
		string(models.QuickScanQueue): {
			asynq.TaskStatePending:  {quickScanTaskInfo(t, "older", asynq.TaskStatePending, QuickScanTaskPayload{GroupId: 1, TempPath: older})},
			asynq.TaskStateArchived: {quickScanTaskInfo(t, "newer", asynq.TaskStateArchived, QuickScanTaskPayload{GroupId: 1, TempPath: newer})},
		},
		string(models.EmailReceiptProcessingQueue): {
			asynq.TaskStateActive: {emailTaskInfo(t, "mail", asynq.TaskStateActive, EmailProcessTaskPayload{GroupSettingsId: 1})},
		},
	}}

	items, truncated, err := listProcessingQueueItems(lister, resolveByPayload)

	if err != nil || truncated {
		t.Fatalf("unexpected: %v truncated=%v", err, truncated)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
	if items[0].AsynqTaskId != "newer" || items[1].AsynqTaskId != "older" || items[2].AsynqTaskId != "mail" {
		t.Errorf("order wrong: %s %s %s", items[0].AsynqTaskId, items[1].AsynqTaskId, items[2].AsynqTaskId)
	}
}

func TestListProcessingQueueItems_CapsEachStateAndReportsTruncation(t *testing.T) {
	var many []*asynq.TaskInfo
	for i := 0; i < processingQueuePerStateCap; i++ {
		many = append(many, quickScanTaskInfo(t, fmt.Sprintf("t%d", i), asynq.TaskStatePending, QuickScanTaskPayload{GroupId: 1}))
	}
	lister := &fakeQueueLister{tasks: map[string]map[asynq.TaskState][]*asynq.TaskInfo{
		string(models.QuickScanQueue): {asynq.TaskStatePending: many},
	}}

	_, truncated, err := listProcessingQueueItems(lister, resolveByPayload)

	if err != nil {
		t.Fatal(err)
	}
	if !truncated {
		t.Error("a full page must report truncation")
	}
	if len(lister.optionCounts) == 0 {
		t.Fatal("no list call was made")
	}
	for _, count := range lister.optionCounts {
		// A page size and a page number: asynq options are opaque, so the cap
		// itself is pinned by the truncation check above.
		if count != 2 {
			t.Errorf("list call carried %d options, want page size and page", count)
		}
	}
}

func TestListProcessingQueueItems_UnknownQueueIsEmptyNotAnError(t *testing.T) {
	lister := &fakeQueueLister{errByQ: map[string]error{
		string(models.QuickScanQueue):                fmt.Errorf("asynq: %w", asynq.ErrQueueNotFound),
		string(models.EmailReceiptProcessingQueue): fmt.Errorf("asynq: %w", asynq.ErrQueueNotFound),
	}}

	items, _, err := listProcessingQueueItems(lister, resolveByPayload)

	if err != nil || items == nil || len(items) != 0 {
		t.Fatalf("expected an empty non-nil list, got %v %v", items, err)
	}
}

func TestListProcessingQueueItems_RedisFailureIsAnError(t *testing.T) {
	lister := &fakeQueueLister{errByQ: map[string]error{string(models.QuickScanQueue): errors.New("connection refused")}}

	if _, _, err := listProcessingQueueItems(lister, resolveByPayload); err == nil {
		t.Fatal("expected the error to surface")
	}
}

func TestDescribeProcessingQueueTask(t *testing.T) {
	info := quickScanTaskInfo(t, "q1", asynq.TaskStateActive, QuickScanTaskPayload{Token: &structs.Claims{UserId: 5}, GroupId: 2})

	task, err := describeProcessingQueueTask(info, resolveByPayload)
	if err != nil || task.GroupId != 2 || task.TaskType != models.QUICK_SCAN || task.UploaderUserId == nil || *task.UploaderUserId != 5 {
		t.Fatalf("unexpected: %v %+v", err, task)
	}

	completed := quickScanTaskInfo(t, "q2", asynq.TaskStateCompleted, QuickScanTaskPayload{GroupId: 2})
	if _, err := describeProcessingQueueTask(completed, resolveByPayload); !errors.Is(err, ErrProcessingQueueTaskNotFound) {
		t.Errorf("a completed task is not part of the queue, got %v", err)
	}

	email := emailTaskInfo(t, "e1", asynq.TaskStateRetry, EmailProcessTaskPayload{GroupSettingsId: 99})
	if _, err := describeProcessingQueueTask(email, resolveByPayload); !errors.Is(err, ErrProcessingQueueTaskNotFound) {
		t.Errorf("an unresolvable group must read as not found, got %v", err)
	}
}

func TestIsProcessingQueueName(t *testing.T) {
	for _, name := range []string{"quick_scan", "email_receipt_processing"} {
		if !IsProcessingQueueName(name) {
			t.Errorf("%s should be an upload queue", name)
		}
	}
	for _, name := range []string{"", "email_polling", "system_clean_up", "../quick_scan"} {
		if IsProcessingQueueName(name) {
			t.Errorf("%q should not be an upload queue", name)
		}
	}
}

// The refusals happen before any Redis call, so they run without Redis. Only a
// retrying task or a failed one with its file may be run; pending, active and
// completed must never reach Inspector.RunTask, which would push them back.
func TestRunProcessingQueueTaskNow_RefusesWhatItMust(t *testing.T) {
	directory := containedTempDir(t)
	missing := directory + "/gone.jpg"

	for _, state := range []asynq.TaskState{asynq.TaskStatePending, asynq.TaskStateActive, asynq.TaskStateScheduled, asynq.TaskStateCompleted} {
		task := ProcessingQueueTask{
			Info:     quickScanTaskInfo(t, "q", state, QuickScanTaskPayload{GroupId: 1}),
			TaskType: models.QUICK_SCAN,
		}
		if err := RunProcessingQueueTaskNow(task); !errors.Is(err, ErrProcessingQueueNotRetryable) {
			t.Errorf("state %v: expected not retryable, got %v", state, err)
		}
	}

	archived := ProcessingQueueTask{
		Info:     quickScanTaskInfo(t, "q", asynq.TaskStateArchived, QuickScanTaskPayload{GroupId: 1, TempPath: missing}),
		TaskType: models.QUICK_SCAN,
	}
	if err := RunProcessingQueueTaskNow(archived); !errors.Is(err, ErrProcessingQueueFileMissing) {
		t.Errorf("archived without its file: expected file missing, got %v", err)
	}
}

func TestListProcessingQueueItems_BelowTheCapIsNotTruncated(t *testing.T) {
	var few []*asynq.TaskInfo
	for i := 0; i < processingQueuePerStateCap-1; i++ {
		few = append(few, quickScanTaskInfo(t, fmt.Sprintf("t%d", i), asynq.TaskStatePending, QuickScanTaskPayload{GroupId: 1}))
	}
	lister := &fakeQueueLister{tasks: map[string]map[asynq.TaskState][]*asynq.TaskInfo{
		string(models.QuickScanQueue): {asynq.TaskStatePending: few},
	}}

	items, truncated, err := listProcessingQueueItems(lister, resolveByPayload)

	if err != nil || truncated || len(items) != processingQueuePerStateCap-1 {
		t.Fatalf("unexpected: %v truncated=%v items=%d", err, truncated, len(items))
	}
}
