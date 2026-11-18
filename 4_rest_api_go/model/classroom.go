package model

type ClassroomRequest struct {
	ID       int64 `json:"id"`
	SchoolID int64 `json:"school_id" validate:"gt=0"`
	Year     int   `json:"year" validate:"gt=0"`
	Section  int   `json:"section" validate:"gt=0"`
}
