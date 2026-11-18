package service

import (
	"context"
	"errors"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/model"
	responsestruct "github.com/tonklaklach/swd-backend/4_rest_api_go/response/struct"
)

type ClassroomService interface {
	List(context.Context, domain.Filter) ([]domain.Classroom, error)
	Get(context.Context, int64) (responsestruct.ClassroomDetail, error)
	Create(context.Context, model.ClassroomRequest) (domain.Classroom, error)
	Update(context.Context, int64, model.ClassroomRequest) (domain.Classroom, error)
	Delete(context.Context, int64) error
}

type classroomService struct {
	repo     ClassroomRepository
	teachers TeacherReader
	students StudentReader
}

func NewClassroom(repo ClassroomRepository, teachers TeacherReader, students StudentReader) ClassroomService {
	return &classroomService{repo: repo, teachers: teachers, students: students}
}
func (s *classroomService) List(ctx context.Context, filter domain.Filter) ([]domain.Classroom, error) {
	return s.repo.ListClassroom(ctx, filter)
}
func (s *classroomService) Create(ctx context.Context, req model.ClassroomRequest) (domain.Classroom, error) {
	return s.save(ctx, 0, req)
}
func (s *classroomService) Update(ctx context.Context, id int64, req model.ClassroomRequest) (domain.Classroom, error) {
	return s.save(ctx, id, req)
}
func (s *classroomService) save(ctx context.Context, id int64, req model.ClassroomRequest) (domain.Classroom, error) {
	item := domain.Classroom{SchoolID: req.SchoolID, Year: req.Year, Section: req.Section}
	return s.repo.SaveClassroom(ctx, id, item)
}
func (s *classroomService) Delete(ctx context.Context, id int64) error {
	err := s.repo.DeleteClassroom(ctx, id)
	if errors.Is(err, domain.ErrInvalidReference) {
		return domain.ErrConflict
	}
	return err
}
func (s *classroomService) Get(ctx context.Context, id int64) (responsestruct.ClassroomDetail, error) {
	item, err := s.repo.GetClassroom(ctx, id)
	if err != nil {
		return responsestruct.ClassroomDetail{}, err
	}
	teachers, err := s.teachers.ListTeacher(ctx, domain.Filter{ClassroomID: id, Limit: -1})
	if err != nil {
		return responsestruct.ClassroomDetail{}, err
	}
	students, err := s.students.ListStudent(ctx, domain.Filter{ClassroomID: id, Limit: -1})
	return responsestruct.ClassroomDetail{Classroom: item, Teachers: teachers, Students: students}, err
}
