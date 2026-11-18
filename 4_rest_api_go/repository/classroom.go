package repository

import (
	"context"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
	"gorm.io/gorm"
)

func (s *Repository) ListClassroom(ctx context.Context, input domain.Filter) ([]domain.Classroom, error) {
	items := []domain.Classroom{}
	err := s.db.WithContext(ctx).Scopes(buildFilter("classrooms", input)).Find(&items).Error
	return items, mapError(err)
}
func (s *Repository) GetClassroom(ctx context.Context, id int64) (domain.Classroom, error) {
	var item domain.Classroom
	err := s.db.WithContext(ctx).First(&item, id).Error
	return item, mapError(err)
}
func (s *Repository) SaveClassroom(ctx context.Context, id int64, item domain.Classroom) (domain.Classroom, error) {
	item.ID = id
	var saved domain.Classroom
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if id == 0 {
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		} else {
			result := tx.Model(&domain.Classroom{}).Where("id = ?", id).Select("school_id", "year", "section").Updates(item)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return domain.ErrNotFound
			}
		}
		return tx.First(&saved, item.ID).Error
	})
	if err != nil {
		return domain.Classroom{}, mapError(err)
	}
	return saved, nil
}
func (s *Repository) DeleteClassroom(ctx context.Context, id int64) error {
	result := s.db.WithContext(ctx).Delete(&domain.Classroom{}, id)
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
