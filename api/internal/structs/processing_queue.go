package structs

import "time"

// Processing queue item states. Plain strings on the wire, never an OpenAPI enum
// and never a SystemTaskStatus value: a closed enum on a response model fails the
// whole payload on every released Dart client that does not know the new value.
const (
	// ProcessingQueueStateQueued is asynq pending or scheduled: waiting for a worker.
	ProcessingQueueStateQueued = "QUEUED"
	// ProcessingQueueStateProcessing is asynq active: a worker is running it now.
	ProcessingQueueStateProcessing = "PROCESSING"
	// ProcessingQueueStateRetrying is asynq retry: an attempt failed and another is scheduled.
	ProcessingQueueStateRetrying = "RETRYING"
	// ProcessingQueueStateFailed is asynq archived: retries exhausted; waits for a rerun.
	ProcessingQueueStateFailed = "FAILED"
)

// Processing queue task types, matching the system task types they produce.
const (
	ProcessingQueueTaskTypeQuickScan   = "QUICK_SCAN"
	ProcessingQueueTaskTypeEmailUpload = "EMAIL_UPLOAD"
)

// ProcessingQueueItem is one quick scan or email upload that has not finished
// successfully. It is read from asynq, not the database.
type ProcessingQueueItem struct {
	AsynqTaskId string `json:"asynqTaskId"`
	Queue       string `json:"queue"`
	State       string `json:"state"`
	TaskType    string `json:"taskType"`
	GroupId     uint   `json:"groupId"`
	// UploaderUserId is the user who ran the quick scan. Email uploads are polled,
	// not run by a user, so it is absent for them.
	UploaderUserId *uint  `json:"uploaderUserId,omitempty"`
	FileName       string `json:"fileName"`
	// EmailSubject and EmailSender are set for email uploads only.
	EmailSubject string `json:"emailSubject,omitempty"`
	EmailSender  string `json:"emailSender,omitempty"`
	// UploadedAt is the modification time of the upload in temp/, asynq not
	// recording when a task was enqueued. Absent once the file is gone.
	UploadedAt *time.Time `json:"uploadedAt,omitempty"`
	// RetryCount is the number of failed attempts so far; MaxRetries the most
	// asynq will schedule after the first attempt.
	RetryCount    int        `json:"retryCount"`
	MaxRetries    int        `json:"maxRetries"`
	NextRetryAt   *time.Time `json:"nextRetryAt,omitempty"`
	LastError     string     `json:"lastError"`
	LastFailedAt  *time.Time `json:"lastFailedAt,omitempty"`
	HasSourceFile bool       `json:"hasSourceFile"`
	// CanBeRestarted is true for a failed (archived) task whose files are still on
	// disk and whose caller may rerun it. CanRetryNow is true for a RETRYING task.
	CanBeRestarted bool `json:"canBeRestarted"`
	CanRetryNow    bool `json:"canRetryNow"`
	// SystemTaskId addresses the existing /systemTask/{id}/sourceFile,
	// /sourceFile/download and /rerunActivity/{id} endpoints. Absent until a
	// processing step has finished and written a row.
	SystemTaskId *uint `json:"systemTaskId,omitempty"`
}

// ProcessingQueueList is the response of the processing queue list endpoint.
type ProcessingQueueList struct {
	Items []ProcessingQueueItem `json:"items"`
	// Truncated is true when asynq held more tasks in some state than the per-state
	// cap read, so the list may omit some.
	Truncated bool `json:"truncated"`
}
