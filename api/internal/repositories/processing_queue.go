package repositories

import "receipt-wrangler/api/internal/models"

// GetSourceSystemTaskIdsByAsynqTaskIds maps asynq task ids to the id of a
// QUICK_SCAN or EMAIL_UPLOAD system task that the source-file and rerun
// endpoints can address. A primary and a fallback processing setting write two
// rows sharing one asynq id; the lowest id wins. Tasks with no group are left
// out, since those endpoints fail closed on them.
func (repository SystemTaskRepository) GetSourceSystemTaskIdsByAsynqTaskIds(asynqTaskIds []string) (map[string]uint, error) {
	result := make(map[string]uint)
	if len(asynqTaskIds) == 0 {
		return result, nil
	}

	var rows []models.SystemTask
	err := repository.GetDB().
		Model(&models.SystemTask{}).
		Select("id", "asynq_task_id").
		Where("asynq_task_id IN ?", asynqTaskIds).
		Where("type IN ?", []models.SystemTaskType{models.QUICK_SCAN, models.EMAIL_UPLOAD}).
		Where("group_id IS NOT NULL").
		Order("id ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		if _, exists := result[row.AsynqTaskId]; !exists {
			result[row.AsynqTaskId] = row.ID
		}
	}

	return result, nil
}
