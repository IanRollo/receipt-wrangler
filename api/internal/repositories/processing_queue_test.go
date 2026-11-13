package repositories

import (
	"receipt-wrangler/api/internal/models"
	"testing"
)

func TestGetSourceSystemTaskIdsByAsynqTaskIds(t *testing.T) {
	defer TruncateTestDb()
	db := GetDB()

	group := models.Group{Name: "pq-group"}
	db.Create(&group)
	groupId := group.ID

	create := func(taskType models.SystemTaskType, asynqId string, withGroup bool) models.SystemTask {
		task := models.SystemTask{
			Type:                 taskType,
			Status:               models.SYSTEM_TASK_FAILED,
			AssociatedEntityType: models.NOOP_ENTITY_TYPE,
			AsynqTaskId:          asynqId,
		}
		if withGroup {
			task.GroupId = &groupId
		}
		db.Create(&task)
		return task
	}

	primary := create(models.QUICK_SCAN, "task-a", true)
	create(models.QUICK_SCAN, "task-a", true) // the fallback row sharing the asynq id
	create(models.RECEIPT_UPLOADED, "task-b", true)
	create(models.QUICK_SCAN, "task-c", false)
	email := create(models.EMAIL_UPLOAD, "task-d", true)

	got, err := NewSystemTaskRepository(nil).GetSourceSystemTaskIdsByAsynqTaskIds([]string{"task-a", "task-b", "task-c", "task-d", "task-none"})
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 || got["task-a"] != primary.ID || got["task-d"] != email.ID {
		t.Fatalf("expected the lowest id for task-a and task-d only, got %v", got)
	}

	empty, err := NewSystemTaskRepository(nil).GetSourceSystemTaskIdsByAsynqTaskIds(nil)
	if err != nil || empty == nil || len(empty) != 0 {
		t.Fatalf("expected an empty map, got %v %v", empty, err)
	}
}
