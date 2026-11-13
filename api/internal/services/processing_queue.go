package services

import (
	"slices"

	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"

	"gorm.io/gorm"
)

// ProcessingQueueService decides which processing queue items a caller may see.
// Reading the queue itself lives in wranglerasynq, which imports this package.
type ProcessingQueueService struct {
	BaseService
}

func NewProcessingQueueService(tx *gorm.DB) ProcessingQueueService {
	return ProcessingQueueService{BaseService: BaseService{DB: repositories.GetDB(), TX: tx}}
}

// processingQueueAccess answers the three per-group questions the selection asks.
type processingQueueAccess struct {
	// canRead is the gate on seeing the group's activity (group.activities.read).
	canRead func(groupId uint) bool
	// canRerun is group.activities.rerun.
	canRerun func(groupId uint) bool
	// uploaderVisible is member isolation: whether the caller may see this user
	// within the group. Only asked for items that have an uploader.
	uploaderVisible func(groupId uint, uploaderUserId uint) bool
}

// SelectProcessingQueueItems returns the items the caller may see. scopeGroupId is
// the group being browsed, or nil for the synthetic All group, which spans every
// group the caller can read. Each item is judged in ITS OWN group, never the one
// browsed, so the All group's membership grants nothing.
//
// The rule matches the activity feed: group.activities.read in the item's group,
// and an item run by a member the caller cannot see there is hidden. An item with
// no uploader (every email upload) stays visible, like ran_by_user_id IS NULL.
func (service ProcessingQueueService) SelectProcessingQueueItems(
	userId uint,
	scopeGroupId *uint,
	items []structs.ProcessingQueueItem,
	systemTaskIds func(asynqTaskIds []string) (map[string]uint, error),
) ([]structs.ProcessingQueueItem, error) {
	permissionService := NewPermissionService(service.TX)
	groupIds := make([]uint, 0)
	seen := make(map[uint]bool)
	for _, item := range items {
		if !seen[item.GroupId] {
			seen[item.GroupId] = true
			groupIds = append(groupIds, item.GroupId)
		}
	}

	visibility, err := permissionService.ActivityVisibilityResolver(userId, groupIds)
	if err != nil {
		return nil, err
	}

	permissionCache := func(permission string) func(groupId uint) bool {
		cache := make(map[uint]bool)
		return func(groupId uint) bool {
			if allowed, ok := cache[groupId]; ok {
				return allowed
			}
			allowed, err := permissionService.HasGroupPermissions(userId, groupId, permission)
			if err != nil {
				logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
				allowed = false
			}
			cache[groupId] = allowed
			return allowed
		}
	}

	access := processingQueueAccess{
		canRead:  permissionCache(permissions.GroupActivitiesRead),
		canRerun: permissionCache(permissions.GroupActivitiesRerun),
		uploaderVisible: func(groupId uint, uploaderUserId uint) bool {
			visibleIds, unrestricted, err := visibility(groupId)
			if err != nil {
				logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
				return false
			}
			return unrestricted || slices.Contains(visibleIds, uploaderUserId)
		},
	}

	selected := filterProcessingQueueItems(items, scopeGroupId, access)

	ids := make([]string, 0, len(selected))
	for _, item := range selected {
		ids = append(ids, item.AsynqTaskId)
	}

	systemTaskIdByAsynqId, err := systemTaskIds(ids)
	if err != nil {
		return nil, err
	}
	for i := range selected {
		if id, ok := systemTaskIdByAsynqId[selected[i].AsynqTaskId]; ok {
			idCopy := id
			selected[i].SystemTaskId = &idCopy
		}
	}

	return selected, nil
}

// filterProcessingQueueItems is the pure selection, with the permission lookups
// injected so the rules can be tested without a database. It also drops the
// rerun affordances from items the caller may see but not rerun.
func filterProcessingQueueItems(
	items []structs.ProcessingQueueItem,
	scopeGroupId *uint,
	access processingQueueAccess,
) []structs.ProcessingQueueItem {
	selected := make([]structs.ProcessingQueueItem, 0, len(items))

	for _, item := range items {
		if scopeGroupId != nil && item.GroupId != *scopeGroupId {
			continue
		}
		if !access.canRead(item.GroupId) {
			continue
		}
		if item.UploaderUserId != nil && !access.uploaderVisible(item.GroupId, *item.UploaderUserId) {
			continue
		}

		if !access.canRerun(item.GroupId) {
			item.CanBeRestarted = false
			item.CanRetryNow = false
		}

		selected = append(selected, item)
	}

	return selected
}
