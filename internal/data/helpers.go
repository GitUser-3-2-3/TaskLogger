package data

import "TaskLogger/internal/validator"

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

// ValidateTask todo -> incomplete validation method
func ValidateTask(vld *validator.Validator, task *Tasks) {
	vld.CheckError(task.Name != "", "name", "must not be empty")
	vld.CheckError(len(task.Name) > 0 &&
		len(task.Name) <= 50, "name", "cannot be longer than 50 chars")

	vld.CheckError(task.UserID > 0, "user_id", "cannot be zero or negative")
}
