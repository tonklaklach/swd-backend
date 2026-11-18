package responsestruct

import "github.com/tonklaklach/swd-backend/4_rest_api_go/domain"

type SchoolDetail struct {
	domain.School
	Counts map[string]int `json:"counts"`
}
type ClassroomDetail struct {
	domain.Classroom
	Teachers []domain.Teacher `json:"teachers"`
	Students []domain.Student `json:"students"`
}
type TeacherDetail struct {
	domain.Teacher
	Classrooms []domain.Classroom `json:"classrooms"`
}
type StudentDetail struct {
	domain.Student
	Classroom domain.Classroom `json:"classroom"`
}
