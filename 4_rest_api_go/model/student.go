package model

type StudentRequest struct {
	ID          int64  `json:"id"`
	Firstname   string `json:"firstname" validate:"required,notblank"`
	Lastname    string `json:"lastname" validate:"required,notblank"`
	Gender      string `json:"gender" validate:"required,oneof=male female other"`
	ClassroomID int64  `json:"classroom_id" validate:"gt=0"`
}
