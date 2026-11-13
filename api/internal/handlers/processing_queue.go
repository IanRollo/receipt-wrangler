package handlers

import (
	"errors"
	"net/http"

	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/stacktrace"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"receipt-wrangler/api/internal/wranglerasynq"

	"github.com/go-chi/chi/v5"
)

// The processing queue (swagger tag ProcessingQueue): quick scans and email
// uploads that are waiting, running, retrying or failed, read from asynq. It
// reuses the activity feed permissions rather than adding keys to the
// Permission enum (a new value fails already-released Dart clients): seeing an
// item and its file is group.activities.read, running one again is
// group.activities.rerun, both in the item's OWN group.

// GetProcessingQueueForGroup lists the queue for one group, or for the synthetic
// All group across every group the caller can read.
func GetProcessingQueueForGroup(w http.ResponseWriter, r *http.Request) {
	errorMsg := "Error getting processing queue"
	groupId := chi.URLParam(r, "groupId")

	uintGroupId, err := utils.StringToUint(groupId)
	if err != nil {
		utils.WriteCustomErrorResponse(w, errorMsg, http.StatusBadRequest)
		return
	}

	isAllGroup, err := repositories.NewGroupRepository(nil).IsAllGroup(uintGroupId)
	if err != nil {
		// An unknown group answers like one the caller may not reach.
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(w, unauthorizedEntityMessage, http.StatusForbidden)
		return
	}

	handler := structs.Handler{
		ErrorMessage: errorMsg,
		Writer:       w,
		Request:      r,
		ResponseType: constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			items, truncated, err := wranglerasynq.GetProcessingQueueItems()
			if err != nil {
				return http.StatusInternalServerError, stacktrace.WithStack(err)
			}

			var scopeGroupId *uint
			if !isAllGroup {
				scopeGroupId = &uintGroupId
			}

			selected, err := services.NewProcessingQueueService(nil).SelectProcessingQueueItems(
				structs.GetClaims(r).UserId,
				scopeGroupId,
				items,
				repositories.NewSystemTaskRepository(nil).GetSourceSystemTaskIdsByAsynqTaskIds,
			)
			if err != nil {
				return http.StatusInternalServerError, stacktrace.WithStack(err)
			}

			return writeClaimResponse(w, structs.ProcessingQueueList{Items: selected, Truncated: truncated})
		},
	}

	// A real group is gated up front. The All group has no activity of its own and
	// its membership grants nothing, so it is gated per item instead.
	if !isAllGroup {
		handler.GroupId = groupId
		handler.GroupPermissions = []string{permissions.GroupActivitiesRead}
	}

	HandleRequest(handler)
}

// GetProcessingQueueTaskSourceFile previews the upload behind a queue item. It is
// keyed by queue and asynq task id because a task that has not yet finished a
// processing step has no system task row for the activity routes to address.
func GetProcessingQueueTaskSourceFile(w http.ResponseWriter, r *http.Request) {
	errorMsg := "Error getting processing queue source file"

	task, groupId, handled := loadProcessingQueueTask(w, r, errorMsg)
	if handled {
		return
	}

	HandleRequest(structs.Handler{
		ErrorMessage:     errorMsg,
		Writer:           w,
		Request:          r,
		GroupId:          groupId,
		GroupPermissions: []string{permissions.GroupActivitiesRead},
		ResponseType:     constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			sourceFile, ok := authorizeAndResolveSourceFile(w, r, processingQueueSystemTask(task), errorMsg)
			if !ok {
				return 0, nil
			}

			return writeSourceFilePreview(w, sourceFile)
		},
	})
}

// DownloadProcessingQueueTaskSourceFile serves the upload behind a queue item verbatim.
func DownloadProcessingQueueTaskSourceFile(w http.ResponseWriter, r *http.Request) {
	errorMsg := "Error downloading processing queue source file"

	task, groupId, handled := loadProcessingQueueTask(w, r, errorMsg)
	if handled {
		return
	}

	HandleRequest(structs.Handler{
		ErrorMessage:     errorMsg,
		Writer:           w,
		Request:          r,
		GroupId:          groupId,
		GroupPermissions: []string{permissions.GroupActivitiesRead},
		// Deliberately unset so http.ServeFile derives the Content-Type itself.
		ResponseType: "",
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			sourceFile, ok := authorizeAndResolveSourceFile(w, r, processingQueueSystemTask(task), errorMsg)
			if !ok {
				return 0, nil
			}

			return writeSourceFileDownload(w, r, sourceFile)
		},
	})
}

// RetryProcessingQueueTask runs a RETRYING task now, or reruns a FAILED one.
// Anything else is a 400: Inspector.RunTask would push back a pending or
// completed task too, so the state rule lives in
// wranglerasynq.RunProcessingQueueTaskNow, as it does for the rerun endpoint.
func RetryProcessingQueueTask(w http.ResponseWriter, r *http.Request) {
	errorMsg := "Error retrying upload"

	task, groupId, handled := loadProcessingQueueTask(w, r, errorMsg)
	if handled {
		return
	}

	HandleRequest(structs.Handler{
		ErrorMessage:     errorMsg,
		Writer:           w,
		Request:          r,
		GroupId:          groupId,
		GroupPermissions: []string{permissions.GroupActivitiesRerun},
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			if denied := enforceActivityActorVisible(w, r, processingQueueSystemTask(task), errorMsg); denied {
				return 0, nil
			}

			err := wranglerasynq.RunProcessingQueueTaskNow(task)
			if errors.Is(err, wranglerasynq.ErrProcessingQueueNotRetryable) {
				utils.WriteCustomErrorResponse(w, "Only a retrying or failed upload can be retried.", http.StatusBadRequest)
				return 0, nil
			}
			if errors.Is(err, wranglerasynq.ErrProcessingQueueFileMissing) {
				utils.WriteCustomErrorResponse(w, "The file this upload needs is no longer available.", http.StatusBadRequest)
				return 0, nil
			}
			if err != nil {
				return http.StatusInternalServerError, stacktrace.WithStack(err)
			}

			return 0, nil
		},
	})
}

// loadProcessingQueueTask finds the task named in the URL and returns the group
// the request must be gated on. It reports true when it has already written a
// response. A missing task, a foreign queue and an undecodable payload all answer
// with the access-denied body a group the caller cannot reach would, so the status
// is no oracle for task ids.
func loadProcessingQueueTask(
	w http.ResponseWriter,
	r *http.Request,
	errorMsg string,
) (wranglerasynq.ProcessingQueueTask, string, bool) {
	task, err := wranglerasynq.LoadProcessingQueueTask(chi.URLParam(r, "queue"), chi.URLParam(r, "taskId"))
	if err != nil {
		if errors.Is(err, wranglerasynq.ErrProcessingQueueTaskNotFound) {
			utils.WriteCustomErrorResponse(w, activityAccessDeniedMessage, http.StatusForbidden)
			return wranglerasynq.ProcessingQueueTask{}, "", true
		}

		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(w, errorMsg, http.StatusInternalServerError)
		return wranglerasynq.ProcessingQueueTask{}, "", true
	}

	return task, utils.UintToString(task.GroupId), false
}

// processingQueueSystemTask is the part of a system task the source-file helpers
// read, built from the queue item so they can be shared.
func processingQueueSystemTask(task wranglerasynq.ProcessingQueueTask) models.SystemTask {
	groupId := task.GroupId

	return models.SystemTask{
		Type:        task.TaskType,
		AsynqTaskId: task.Info.ID,
		GroupId:     &groupId,
		RanByUserId: task.UploaderUserId,
	}
}
