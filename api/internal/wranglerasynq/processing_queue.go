package wranglerasynq

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"time"

	"github.com/hibiken/asynq"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

// processingQueuePerStateCap is the most tasks read from any one asynq state of
// one queue. Each list call is a Redis round trip that also ships every payload
// back, and an email payload carries its body twice (see completedTaskRetention),
// so an unbounded read lets a backlog make this endpoint slow. 100 is far above
// what this install ever holds, and a list that hits it reports Truncated.
// Pending lists oldest first, so in a backlog past the cap it is the newest
// uploads that are cut; archived lists the newest failure first.
const processingQueuePerStateCap = 100

// processingQueueQueues are the two queues that hold uploads, with the task type
// each carries.
var processingQueueQueues = []struct {
	queue    models.QueueName
	taskType string
}{
	{models.QuickScanQueue, structs.ProcessingQueueTaskTypeQuickScan},
	{models.EmailReceiptProcessingQueue, structs.ProcessingQueueTaskTypeEmailUpload},
}

// IsProcessingQueueName reports whether a queue name from a URL is one of the
// upload queues.
func IsProcessingQueueName(queue string) bool {
	for _, entry := range processingQueueQueues {
		if string(entry.queue) == queue {
			return true
		}
	}

	return false
}

// processingQueueLister is the part of *asynq.Inspector the listing uses. Named so
// the mapping can be tested without Redis; in production it is always the Inspector.
type processingQueueLister interface {
	ListPendingTasks(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error)
	ListScheduledTasks(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error)
	ListActiveTasks(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error)
	ListRetryTasks(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error)
	ListArchivedTasks(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error)
}

// processingQueueState maps an asynq state to the queue's own plain-string state.
// Completed and aggregating tasks are not part of the queue (ok == false).
func processingQueueState(state asynq.TaskState) (string, bool) {
	switch state {
	case asynq.TaskStatePending, asynq.TaskStateScheduled:
		return structs.ProcessingQueueStateQueued, true
	case asynq.TaskStateActive:
		return structs.ProcessingQueueStateProcessing, true
	case asynq.TaskStateRetry:
		return structs.ProcessingQueueStateRetrying, true
	case asynq.TaskStateArchived:
		return structs.ProcessingQueueStateFailed, true
	}

	return "", false
}

// groupIdResolver finds the group a decoded payload belongs to.
type groupIdResolver func(payload RerunTaskPayload) (uint, error)

// GetProcessingQueueItems reads every quick scan and email upload task that has
// not succeeded, newest first. Visibility is not applied here.
func GetProcessingQueueItems() ([]structs.ProcessingQueueItem, bool, error) {
	inspector, err := GetAsynqInspector()
	if err != nil {
		return nil, false, err
	}
	defer inspector.Close()

	return listProcessingQueueItems(inspector, memoizedGroupIdResolver())
}

func listProcessingQueueItems(
	lister processingQueueLister,
	resolveGroup groupIdResolver,
) ([]structs.ProcessingQueueItem, bool, error) {
	items := make([]structs.ProcessingQueueItem, 0)
	truncated := false

	for _, entry := range processingQueueQueues {
		queue := string(entry.queue)
		listers := []func(string, ...asynq.ListOption) ([]*asynq.TaskInfo, error){
			lister.ListPendingTasks,
			lister.ListScheduledTasks,
			lister.ListActiveTasks,
			lister.ListRetryTasks,
			lister.ListArchivedTasks,
		}

		for _, list := range listers {
			taskInfos, err := list(queue, asynq.PageSize(processingQueuePerStateCap), asynq.Page(1))
			if err != nil {
				if errors.Is(err, asynq.ErrQueueNotFound) {
					// Nothing has ever been enqueued on this queue.
					continue
				}

				return nil, false, err
			}

			if len(taskInfos) >= processingQueuePerStateCap {
				truncated = true
			}

			for _, taskInfo := range taskInfos {
				item, ok := buildProcessingQueueItem(taskInfo, entry.taskType, resolveGroup)
				if ok {
					items = append(items, item)
				}
			}
		}
	}

	sortProcessingQueueItems(items)

	return items, truncated, nil
}

// sortProcessingQueueItems orders newest first by the best time available: when
// the upload landed, else when it last failed. Equal keys keep their order.
func sortProcessingQueueItems(items []structs.ProcessingQueueItem) {
	sortKey := func(item structs.ProcessingQueueItem) time.Time {
		if item.UploadedAt != nil {
			return *item.UploadedAt
		}
		if item.LastFailedAt != nil {
			return *item.LastFailedAt
		}

		return time.Time{}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return sortKey(items[i]).After(sortKey(items[j]))
	})
}

// buildProcessingQueueItem decodes one task. A task that is not an upload, whose
// payload cannot be decoded, or whose group cannot be resolved is logged and
// skipped (ok == false): without a group there is nothing to authorize it against.
func buildProcessingQueueItem(
	taskInfo *asynq.TaskInfo,
	queueTaskType string,
	resolveGroup groupIdResolver,
) (structs.ProcessingQueueItem, bool) {
	state, ok := processingQueueState(taskInfo.State)
	if !ok {
		return structs.ProcessingQueueItem{}, false
	}

	var payload RerunTaskPayload
	var systemTaskType models.SystemTaskType
	var path string

	switch taskInfo.Type {
	case QuickScan:
		systemTaskType = models.QUICK_SCAN
		if err := json.Unmarshal(taskInfo.Payload, &payload.QuickScanTaskPayload); err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, "Processing queue: cannot decode task ", taskInfo.ID, ": ", err.Error())
			return structs.ProcessingQueueItem{}, false
		}
		path = payload.TempPath
	case EmailProcess:
		systemTaskType = models.EMAIL_UPLOAD
		if err := json.Unmarshal(taskInfo.Payload, &payload.EmailProcessTaskPayload); err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, "Processing queue: cannot decode task ", taskInfo.ID, ": ", err.Error())
			return structs.ProcessingQueueItem{}, false
		}
		path = payload.TempFilePath
	default:
		return structs.ProcessingQueueItem{}, false
	}

	groupId, err := resolveGroup(payload)
	if err != nil || groupId == 0 {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Processing queue: cannot resolve the group of task ", taskInfo.ID)
		return structs.ProcessingQueueItem{}, false
	}

	files, err := taskSourceFilesFromPayload(systemTaskType, taskInfo.Payload)
	if err != nil {
		return structs.ProcessingQueueItem{}, false
	}
	flags := flagsForState(taskInfo.State, files)

	item := structs.ProcessingQueueItem{
		AsynqTaskId:    taskInfo.ID,
		Queue:          taskInfo.Queue,
		State:          state,
		TaskType:       queueTaskType,
		GroupId:        groupId,
		RetryCount:     taskInfo.Retried,
		MaxRetries:     taskInfo.MaxRetry,
		LastError:      taskInfo.LastErr,
		HasSourceFile:  flags.hasSourceFile,
		CanBeRestarted: flags.canBeRestarted,
		CanRetryNow:    taskInfo.State == asynq.TaskStateRetry,
	}

	if systemTaskType == models.QUICK_SCAN {
		item.FileName = payload.OriginalFileName
		if payload.Token != nil {
			uploader := payload.Token.UserId
			item.UploaderUserId = &uploader
		}
	} else {
		item.FileName = payload.Attachment.Filename
		item.EmailSubject = payload.Metadata.Subject
		item.EmailSender = payload.Metadata.FromEmail
	}

	if flags.hasSourceFile {
		if stat, err := os.Stat(path); err == nil {
			modified := stat.ModTime()
			item.UploadedAt = &modified
		}
	}
	if !taskInfo.LastFailedAt.IsZero() {
		failedAt := taskInfo.LastFailedAt
		item.LastFailedAt = &failedAt
	}
	if taskInfo.State == asynq.TaskStateRetry && !taskInfo.NextProcessAt.IsZero() {
		next := taskInfo.NextProcessAt
		item.NextRetryAt = &next
	}

	return item, true
}

// memoizedGroupIdResolver resolves a payload's group: quick scan carries it,
// email carries the group settings it was polled for, looked up once per id.
func memoizedGroupIdResolver() groupIdResolver {
	cache := make(map[uint]uint)
	repository := repositories.NewGroupSettingsRepository(nil)

	return func(payload RerunTaskPayload) (uint, error) {
		if payload.GroupSettingsId == 0 {
			return payload.GroupId, nil
		}

		if groupId, ok := cache[payload.GroupSettingsId]; ok {
			return groupId, nil
		}

		groupSettings, err := repository.GetGroupSettingsById(utils.UintToString(payload.GroupSettingsId))
		if err != nil {
			return 0, err
		}

		cache[payload.GroupSettingsId] = groupSettings.GroupId
		return groupSettings.GroupId, nil
	}
}

// ProcessingQueueTask is one upload task located by queue and id, with the group
// and uploader the permission gates need.
type ProcessingQueueTask struct {
	Info           *asynq.TaskInfo
	TaskType       models.SystemTaskType
	GroupId        uint
	UploaderUserId *uint
}

// ErrProcessingQueueTaskNotFound covers a missing task, one that is not an upload
// and one whose group cannot be resolved; callers answer all three the same way.
var ErrProcessingQueueTaskNotFound = errors.New("processing queue task not found")

// LoadProcessingQueueTask finds one upload task. The queue must be one of the
// upload queues.
func LoadProcessingQueueTask(queue string, asynqTaskId string) (ProcessingQueueTask, error) {
	if !IsProcessingQueueName(queue) || len(asynqTaskId) == 0 {
		return ProcessingQueueTask{}, ErrProcessingQueueTaskNotFound
	}

	inspector, err := GetAsynqInspector()
	if err != nil {
		return ProcessingQueueTask{}, err
	}
	defer inspector.Close()

	info, err := inspector.GetTaskInfo(queue, asynqTaskId)
	if err != nil {
		if errors.Is(err, asynq.ErrTaskNotFound) || errors.Is(err, asynq.ErrQueueNotFound) {
			return ProcessingQueueTask{}, ErrProcessingQueueTaskNotFound
		}

		return ProcessingQueueTask{}, err
	}

	return describeProcessingQueueTask(info, memoizedGroupIdResolver())
}

func describeProcessingQueueTask(info *asynq.TaskInfo, resolveGroup groupIdResolver) (ProcessingQueueTask, error) {
	if _, ok := processingQueueState(info.State); !ok {
		return ProcessingQueueTask{}, ErrProcessingQueueTaskNotFound
	}

	var payload RerunTaskPayload
	task := ProcessingQueueTask{Info: info}

	switch info.Type {
	case QuickScan:
		task.TaskType = models.QUICK_SCAN
		if err := json.Unmarshal(info.Payload, &payload.QuickScanTaskPayload); err != nil {
			return ProcessingQueueTask{}, ErrProcessingQueueTaskNotFound
		}
		if payload.Token != nil {
			uploader := payload.Token.UserId
			task.UploaderUserId = &uploader
		}
	case EmailProcess:
		task.TaskType = models.EMAIL_UPLOAD
		if err := json.Unmarshal(info.Payload, &payload.EmailProcessTaskPayload); err != nil {
			return ProcessingQueueTask{}, ErrProcessingQueueTaskNotFound
		}
	default:
		return ProcessingQueueTask{}, ErrProcessingQueueTaskNotFound
	}

	groupId, err := resolveGroup(payload)
	if err != nil || groupId == 0 {
		return ProcessingQueueTask{}, ErrProcessingQueueTaskNotFound
	}
	task.GroupId = groupId

	return task, nil
}

// ErrProcessingQueueNotRetryable is returned for a task that is not in a state
// that may be run again.
var ErrProcessingQueueNotRetryable = errors.New("only a retrying or failed upload can be retried")

// ErrProcessingQueueFileMissing is returned when a rerun's upload is gone.
var ErrProcessingQueueFileMissing = errors.New("the file this upload needs is no longer available")

// RunProcessingQueueTaskNow runs a task straight away. Inspector.RunTask does not
// refuse by state, so the rule is here (see rerunnableState): a RETRYING task is
// moved to pending, which asynq would do on its own at the next retry time, and a
// FAILED (archived) one is rerun exactly as the rerun endpoint does. Active and
// pending tasks are refused.
func RunProcessingQueueTaskNow(task ProcessingQueueTask) error {
	switch task.Info.State {
	case asynq.TaskStateRetry:
	case asynq.TaskStateArchived:
		if !CanRerunTask(task.Info) || !RerunSourceFilesPresent(task.TaskType, task.Info.Payload) {
			return ErrProcessingQueueFileMissing
		}
	default:
		return ErrProcessingQueueNotRetryable
	}

	inspector, err := GetAsynqInspector()
	if err != nil {
		return err
	}
	defer inspector.Close()

	return inspector.RunTask(task.Info.Queue, task.Info.ID)
}
