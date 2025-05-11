package data

import (
	"TaskLogger/internal/validator"
	"slices"
	"time"
)

func ValidateCategory(vld *validator.Validator, category *Categories) {
	vld.CheckError(category.Name != "", "name", "must not be empty")
	vld.CheckError(len(category.Name) > 0 &&
		len(category.Name) <= 50, "name", "cannot be longer than 50 chars")

	vld.CheckError(category.UserID > 0, "user_id", "cannot be zero or negative")
}

func (ctg *Categories) ApplyPartialUpdatesToCtg(name, color *string, userID *int64) {
	if name != nil {
		ctg.Name = *name
	}
	if color != nil {
		ctg.Color = *color
	}
	if userID != nil {
		ctg.UserID = *userID
	}
}

func ValidateTask(vld *validator.Validator, task *Tasks) {
	vld.CheckError(task.Name != "", "name", "must not be empty")
	vld.CheckError(len(task.Name) <= 50, "name", "cannot be longer than 50 characters")

	vld.CheckError(len(task.Description) <= 250, "description", "cannot exceed 250 characters")

	validStatuses := []StatusType{StatusNotStarted, StatusInProgress, StatusPaused, StatusCompleted}
	vld.CheckError(slices.Contains(validStatuses, task.Status),
		"status", "cannot be other than (Not Started, In Progress, Paused, Completed)")

	vld.CheckError(task.Priority >= 1 && task.Priority <= 5, "priority", "must be between 1 and 5")

	if task.Deadline != nil {
		vld.CheckError(task.Deadline.After(time.Now()), "deadline", "must be in the future")
	}
	vld.CheckError(task.UserID > 0, "user_id", "cannot be zero or negative")

	if task.CategoryID != nil {
		vld.CheckError(*task.CategoryID > 0, "category_id", "cannot be zero or negative")
	}
}

func (task *Tasks) ApplyPartialUpdatesToTask(name, description, image *string, status *StatusType,
	priority *int, deadline *time.Time, userId, categoryID *int64,
) {
	if name != nil {
		task.Name = *name
	}
	if description != nil {
		task.Description = *description
	}
	if image != nil {
		task.Image = *image
	}
	if status != nil {
		task.Status = *status
	}
	if priority != nil {
		task.Priority = *priority
	}
	if deadline != nil {
		task.Deadline = deadline
	}
	if userId != nil {
		task.UserID = *userId
	}
	if categoryID != nil {
		task.CategoryID = categoryID
	}
}
