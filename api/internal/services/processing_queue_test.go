package services

import (
	"testing"

	"receipt-wrangler/api/internal/structs"
)

func uintPointer(value uint) *uint { return &value }

func queueItem(taskId string, groupId uint, uploader *uint) structs.ProcessingQueueItem {
	return structs.ProcessingQueueItem{
		AsynqTaskId:    taskId,
		GroupId:        groupId,
		UploaderUserId: uploader,
		CanBeRestarted: true,
		CanRetryNow:    true,
	}
}

func taskIds(items []structs.ProcessingQueueItem) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.AsynqTaskId)
	}

	return ids
}

func allowGroups(groupIds ...uint) func(uint) bool {
	return func(groupId uint) bool {
		for _, allowed := range groupIds {
			if allowed == groupId {
				return true
			}
		}

		return false
	}
}

func everyoneVisible(uint, uint) bool { return true }

func TestFilterProcessingQueueItems_RealGroupKeepsOnlyThatGroup(t *testing.T) {
	items := []structs.ProcessingQueueItem{queueItem("a", 1, nil), queueItem("b", 2, nil)}
	access := processingQueueAccess{canRead: allowGroups(1, 2), canRerun: allowGroups(1, 2), uploaderVisible: everyoneVisible}

	got := taskIds(filterProcessingQueueItems(items, uintPointer(1), access))

	if len(got) != 1 || got[0] != "a" {
		t.Fatalf("expected only a, got %v", got)
	}
}

func TestFilterProcessingQueueItems_AllGroupSpansReadableGroupsOnly(t *testing.T) {
	items := []structs.ProcessingQueueItem{queueItem("a", 1, nil), queueItem("b", 2, nil), queueItem("c", 3, nil)}
	// Group 2 is one the caller belongs to but may not read activity in: the All
	// group's membership must grant nothing there.
	access := processingQueueAccess{canRead: allowGroups(1, 3), canRerun: allowGroups(1, 3), uploaderVisible: everyoneVisible}

	got := taskIds(filterProcessingQueueItems(items, nil, access))

	if len(got) != 2 || got[0] != "a" || got[1] != "c" {
		t.Fatalf("expected a and c, got %v", got)
	}
}

func TestFilterProcessingQueueItems_IsolationHidesOtherMembersUploads(t *testing.T) {
	mine := uintPointer(10)
	theirs := uintPointer(20)
	items := []structs.ProcessingQueueItem{
		queueItem("mine", 1, mine),
		queueItem("theirs", 1, theirs),
		queueItem("email", 1, nil),
	}
	access := processingQueueAccess{
		canRead:  allowGroups(1),
		canRerun: allowGroups(1),
		uploaderVisible: func(groupId uint, uploaderUserId uint) bool {
			return uploaderUserId == 10
		},
	}

	got := taskIds(filterProcessingQueueItems(items, uintPointer(1), access))

	// An email upload has no uploader and stays visible, like ran_by_user_id IS NULL.
	if len(got) != 2 || got[0] != "mine" || got[1] != "email" {
		t.Fatalf("expected mine and email, got %v", got)
	}
}

func TestFilterProcessingQueueItems_IsolationIsJudgedInTheItemsOwnGroup(t *testing.T) {
	uploader := uintPointer(20)
	items := []structs.ProcessingQueueItem{queueItem("open", 1, uploader), queueItem("isolated", 2, uploader)}
	access := processingQueueAccess{
		canRead:  allowGroups(1, 2),
		canRerun: allowGroups(1, 2),
		uploaderVisible: func(groupId uint, uploaderUserId uint) bool {
			return groupId == 1
		},
	}

	got := taskIds(filterProcessingQueueItems(items, nil, access))

	if len(got) != 1 || got[0] != "open" {
		t.Fatalf("expected only open, got %v", got)
	}
}

func TestFilterProcessingQueueItems_WithoutRerunPermissionKeepsItemButDropsActions(t *testing.T) {
	items := []structs.ProcessingQueueItem{queueItem("a", 1, nil), queueItem("b", 2, nil)}
	access := processingQueueAccess{canRead: allowGroups(1, 2), canRerun: allowGroups(2), uploaderVisible: everyoneVisible}

	got := filterProcessingQueueItems(items, nil, access)

	if len(got) != 2 {
		t.Fatalf("expected both items, got %d", len(got))
	}
	if got[0].CanBeRestarted || got[0].CanRetryNow {
		t.Error("group 1 caller cannot rerun; actions must be off")
	}
	if !got[1].CanBeRestarted || !got[1].CanRetryNow {
		t.Error("group 2 caller can rerun; actions must stay on")
	}
}

func TestFilterProcessingQueueItems_NothingReadableIsAnEmptySliceNotNil(t *testing.T) {
	access := processingQueueAccess{canRead: allowGroups(), canRerun: allowGroups(), uploaderVisible: everyoneVisible}

	got := filterProcessingQueueItems([]structs.ProcessingQueueItem{queueItem("a", 1, nil)}, nil, access)

	if got == nil || len(got) != 0 {
		t.Fatalf("expected a non-nil empty slice, got %v", got)
	}
}
