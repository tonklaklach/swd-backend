package repository

import (
	"context"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
	"gorm.io/gorm"
)

func (s *Repository) ListStudent(ctx context.Context, input domain.Filter) ([]domain.Student, error) {
	items := []domain.Student{}
	err := s.db.WithContext(ctx).Scopes(buildFilter("students", input)).Find(&items).Error
	return items, mapError(err)
}
func (s *Repository) GetStudent(ctx context.Context, id int64) (domain.Student, error) {
	var item domain.Student
	err := s.db.WithContext(ctx).First(&item, id).Error
	return item, mapError(err)
}
func (s *Repository) SaveStudent(ctx context.Context, id int64, item domain.Student) (domain.Student, error) {
	item.ID = id
	var saved domain.Student
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if id == 0 {
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		} else {
			result := tx.Model(&domain.Student{}).Where("id = ?", id).Select("firstname", "lastname", "gender", "classroom_id").Updates(item)
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
		return domain.Student{}, mapError(err)
	}
	return saved, nil
}
func (s *Repository) DeleteStudent(ctx context.Context, id int64) error {
	result := s.db.WithContext(ctx).Delete(&domain.Student{}, id)
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
