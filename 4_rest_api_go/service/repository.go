package service

import (
	"context"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
)

type SchoolRepository interface {
	ListSchool(context.Context, domain.Filter) ([]domain.School, error)
	GetSchool(context.Context, int64) (domain.School, error)
	SaveSchool(context.Context, int64, domain.School) (domain.School, error)
	DeleteSchool(context.Context, int64) error
	SchoolCounts(context.Context, int64) (map[string]int, error)
}
type ClassroomRepository interface {
	ListClassroom(context.Context, domain.Filter) ([]domain.Classroom, error)
	GetClassroom(context.Context, int64) (domain.Classroom, error)
	SaveClassroom(context.Context, int64, domain.Classroom) (domain.Classroom, error)
	DeleteClassroom(context.Context, int64) error
}
type TeacherRepository interface {
	ListTeacher(context.Context, domain.Filter) ([]domain.Teacher, error)
	GetTeacher(context.Context, int64) (domain.Teacher, error)
	SaveTeacher(context.Context, int64, domain.Teacher) (domain.Teacher, error)
	DeleteTeacher(context.Context, int64) error
}
type StudentRepository interface {
	ListStudent(context.Context, domain.Filter) ([]domain.Student, error)
	GetStudent(context.Context, int64) (domain.Student, error)
	SaveStudent(context.Context, int64, domain.Student) (domain.Student, error)
	DeleteStudent(context.Context, int64) error
}
type ClassroomReader interface {
	ListClassroom(context.Context, domain.Filter) ([]domain.Classroom, error)
	GetClassroom(context.Context, int64) (domain.Classroom, error)
}
type TeacherReader interface {
	ListTeacher(context.Context, domain.Filter) ([]domain.Teacher, error)
}
type StudentReader interface {
	ListStudent(context.Context, domain.Filter) ([]domain.Student, error)
}
