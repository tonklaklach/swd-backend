package domain

type School struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Address      string `json:"address"`
}
type Classroom struct {
	ID       int64 `json:"id"`
	SchoolID int64 `json:"school_id"`
	Year     int   `json:"year"`
	Section  int   `json:"section"`
}
type Teacher struct {
	ID           int64   `json:"id"`
	Firstname    string  `json:"firstname"`
	Lastname     string  `json:"lastname"`
	Gender       string  `json:"gender"`
	ClassroomIDs []int64 `json:"classroom_ids"`
}
type Student struct {
	ID          int64  `json:"id"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Gender      string `json:"gender"`
	ClassroomID int64  `json:"classroom_id"`
}
