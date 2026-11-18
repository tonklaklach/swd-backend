package domain

type Filter struct {
	Name        string
	Firstname   string
	Lastname    string
	Gender      string
	SchoolID    int64
	ClassroomID int64
	TeacherID   int64
	Limit       int
	Offset      int
}
