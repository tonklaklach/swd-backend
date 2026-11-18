package service

import (
	"context"
	"errors"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/model"
	responsestruct "github.com/tonklaklach/swd-backend/4_rest_api_go/response/struct"
)

type StudentService interface {
	List(context.Context, domain.Filter) ([]domain.Student, error)
	Get(context.Context, int64) (responsestruct.StudentDetail, error)
	Create(context.Context, model.StudentRequest) (domain.Student, error)
	Update(context.Context, int64, model.StudentRequest) (domain.Student, error)
	Delete(context.Context, int64) error
}

type studentService struct {
	repo       StudentRepository
	classrooms ClassroomReader
}

func NewStudent(repo StudentRepository, classrooms ClassroomReader) StudentService {
	return &studentService{repo: repo, classrooms: classrooms}
}
func (s *studentService) List(ctx context.Context, filter domain.Filter) ([]domain.Student, error) {
	return s.repo.ListStudent(ctx, filter)
}
func (s *studentService) Create(ctx context.Context, req model.StudentRequest) (domain.Student, error) {
	return s.save(ctx, 0, req)
}
func (s *studentService) Update(ctx context.Context, id int64, req model.StudentRequest) (domain.Student, error) {
	return s.save(ctx, id, req)
}
func (s *studentService) save(ctx context.Context, id int64, req model.StudentRequest) (domain.Student, error) {
	item := domain.Student{Firstname: req.Firstname, Lastname: req.Lastname, Gender: req.Gender, ClassroomID: req.ClassroomID}
	return s.repo.SaveStudent(ctx, id, item)
}
func (s *studentService) Delete(ctx context.Context, id int64) error {
	err := s.repo.DeleteStudent(ctx, id)
	if errors.Is(err, domain.ErrInvalidReference) {
		return domain.ErrConflict
	}
	return err
}
func (s *studentService) Get(ctx context.Context, id int64) (responsestruct.StudentDetail, error) {
	item, err := s.repo.GetStudent(ctx, id)
	if err != nil {
		return responsestruct.StudentDetail{}, err
	}
	classroom, err := s.classrooms.GetClassroom(ctx, item.ClassroomID)
	return responsestruct.StudentDetail{Student: item, Classroom: classroom}, err
}
