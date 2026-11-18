package model

type TeacherRequest struct {
	ID           int64   `json:"id"`
	Firstname    string  `json:"firstname" validate:"required,notblank"`
	Lastname     string  `json:"lastname" validate:"required,notblank"`
	Gender       string  `json:"gender" validate:"required,oneof=male female other"`
	ClassroomIDs []int64 `json:"classroom_ids" validate:"omitempty,unique,dive,gt=0"`
}
