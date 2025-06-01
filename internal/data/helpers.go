package data

import (
	"TaskLogger/internal/validator"
	"slices"
	"time"

	"github.com/google/uuid"
)

func ValidateCategory(vld *validator.Validator, ctg *Categories) {
	vld.CheckError(ctg.Name != "", "name", "must not be empty")
	vld.CheckError(len(ctg.Name) > 0 &&
		len(ctg.Name) <= 50, "name", "cannot be longer than 50 chars")

	vld.CheckError(ctg.UserID != "", "user_id", "cannot be zero or negative")
	vld.CheckError(isValidUUID(ctg.UserID), "user_id", "user id is not valid")
}

func isValidUUID(id string) bool {
	if _, err := uuid.Parse(id); err != nil {
		return false
	}
	return true
}

func (ctg *Categories) ApplyPartialUpdatesToCtg(name, color *string, userID *string) {
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

	validStatuses := []StatusType{StatusPending,
		StatusInProgress, StatusPaused, StatusCompleted, StatusCancelled}

	vld.CheckError(slices.Contains(validStatuses, task.Status),
		"status", "cannot be other than (Not Started, In Progress, Paused, Completed)")

	validPriority := []PriorityType{PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent}

	vld.CheckError(slices.Contains(validPriority, task.Priority),
		"priority", "cannot be other than (low, medium, high, urgent)")

	if task.Deadline != nil {
		vld.CheckError(task.Deadline.After(time.Now()), "deadline", "must be in the future")
	}
	vld.CheckError(task.UserID != "", "user_id", "cannot be zero or negative")
	vld.CheckError(isValidUUID(task.UserID), "user_id", "user_id is not valid")

	if task.CategoryID != nil {
		vld.CheckError(*task.CategoryID > 0, "category_id", "cannot be zero or negative")
	}
}

func (task *Tasks) ApplyPartialUpdatesToTask(name, description, image *string, status *StatusType,
	priority *PriorityType, deadline *time.Time, userId *string, categoryID *int64,
) {
	if name != nil {
		task.Name = *name
	}
	if description != nil {
		task.Description = *description
	}
	if image != nil {
		task.ImageUrl = *image
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

func ValidateSession(vld *validator.Validator, session *Session) {
	vld.CheckError(session.TaskID != "", "task_id", "must be a positive integer")
	vld.CheckError(isValidUUID(session.TaskID), "task_id", "task_id is not valid")
	vld.CheckError(!session.StartedAt.IsZero(), "session_start", "must be provided")
	vld.CheckError(!session.EndedAt.IsZero(), "session_end", "must be provided")
	vld.CheckError(session.EndedAt.After(session.StartedAt), "session_end", "must be after session start")
	vld.CheckError(len(session.Note) <= 500, "note", "cannot exceed 500 characters")

	expectedDuration := int(session.EndedAt.Sub(session.StartedAt).Minutes())
	vld.CheckError(session.Duration == expectedDuration, "duration",
		"must match the time difference between start and end times")
}
