package service

import (
	"context"
	"errors"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/model"
	responsestruct "github.com/tonklaklach/swd-backend/4_rest_api_go/response/struct"
)

type SchoolService interface {
	List(context.Context, domain.Filter) ([]domain.School, error)
	Get(context.Context, int64) (responsestruct.SchoolDetail, error)
	Create(context.Context, model.SchoolRequest) (domain.School, error)
	Update(context.Context, int64, model.SchoolRequest) (domain.School, error)
	Delete(context.Context, int64) error
}

type schoolService struct{ repo SchoolRepository }

func NewSchool(repo SchoolRepository) SchoolService { return &schoolService{repo: repo} }
func (s *schoolService) List(ctx context.Context, filter domain.Filter) ([]domain.School, error) {
	return s.repo.ListSchool(ctx, filter)
}
func (s *schoolService) Create(ctx context.Context, req model.SchoolRequest) (domain.School, error) {
	return s.save(ctx, 0, req)
}
func (s *schoolService) Update(ctx context.Context, id int64, req model.SchoolRequest) (domain.School, error) {
	return s.save(ctx, id, req)
}
func (s *schoolService) save(ctx context.Context, id int64, req model.SchoolRequest) (domain.School, error) {
	item := domain.School{Name: req.Name, Abbreviation: req.Abbreviation, Address: req.Address}
	return s.repo.SaveSchool(ctx, id, item)
}
func (s *schoolService) Delete(ctx context.Context, id int64) error {
	err := s.repo.DeleteSchool(ctx, id)
	if errors.Is(err, domain.ErrInvalidReference) {
		return domain.ErrConflict
	}
	return err
}
func (s *schoolService) Get(ctx context.Context, id int64) (responsestruct.SchoolDetail, error) {
	item, err := s.repo.GetSchool(ctx, id)
	if err != nil {
		return responsestruct.SchoolDetail{}, err
	}
	counts, err := s.repo.SchoolCounts(ctx, id)
	return responsestruct.SchoolDetail{School: item, Counts: counts}, err
}
