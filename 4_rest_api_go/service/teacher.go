package service

import (
	"context"
	"errors"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/model"
	responsestruct "github.com/tonklaklach/swd-backend/4_rest_api_go/response/struct"
)

type TeacherService interface {
	List(context.Context, domain.Filter) ([]domain.Teacher, error)
	Get(context.Context, int64) (responsestruct.TeacherDetail, error)
	Create(context.Context, model.TeacherRequest) (domain.Teacher, error)
	Update(context.Context, int64, model.TeacherRequest) (domain.Teacher, error)
	Delete(context.Context, int64) error
}

type teacherService struct {
	repo       TeacherRepository
	classrooms ClassroomReader
}

func NewTeacher(repo TeacherRepository, classrooms ClassroomReader) TeacherService {
	return &teacherService{repo: repo, classrooms: classrooms}
}
func (s *teacherService) List(ctx context.Context, filter domain.Filter) ([]domain.Teacher, error) {
	return s.repo.ListTeacher(ctx, filter)
}
func (s *teacherService) Create(ctx context.Context, req model.TeacherRequest) (domain.Teacher, error) {
	return s.save(ctx, 0, req)
}
func (s *teacherService) Update(ctx context.Context, id int64, req model.TeacherRequest) (domain.Teacher, error) {
	return s.save(ctx, id, req)
}
func (s *teacherService) save(ctx context.Context, id int64, req model.TeacherRequest) (domain.Teacher, error) {
	item := domain.Teacher{Firstname: req.Firstname, Lastname: req.Lastname, Gender: req.Gender, ClassroomIDs: req.ClassroomIDs}
	return s.repo.SaveTeacher(ctx, id, item)
}
func (s *teacherService) Delete(ctx context.Context, id int64) error {
	err := s.repo.DeleteTeacher(ctx, id)
	if errors.Is(err, domain.ErrInvalidReference) {
		return domain.ErrConflict
	}
	return err
}
func (s *teacherService) Get(ctx context.Context, id int64) (responsestruct.TeacherDetail, error) {
	item, err := s.repo.GetTeacher(ctx, id)
	if err != nil {
		return responsestruct.TeacherDetail{}, err
	}
	classrooms, err := s.classrooms.ListClassroom(ctx, domain.Filter{TeacherID: id, Limit: -1})
	return responsestruct.TeacherDetail{Teacher: item, Classrooms: classrooms}, err
}
