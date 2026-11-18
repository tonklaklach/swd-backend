package repository

import (
	"context"
	"database/sql"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
	"gorm.io/gorm"
)

func (s *Repository) ListSchool(ctx context.Context, input domain.Filter) ([]domain.School, error) {
	items := []domain.School{}
	err := s.db.WithContext(ctx).Scopes(buildFilter("schools", input)).Find(&items).Error
	return items, mapError(err)
}
func (s *Repository) GetSchool(ctx context.Context, id int64) (domain.School, error) {
	var item domain.School
	err := s.db.WithContext(ctx).First(&item, id).Error
	return item, mapError(err)
}
func (s *Repository) SaveSchool(ctx context.Context, id int64, item domain.School) (domain.School, error) {
	item.ID = id
	var saved domain.School
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if id == 0 {
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		} else {
			result := tx.Model(&domain.School{}).Where("id = ?", id).Select("name", "abbreviation", "address").Updates(item)
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
		return domain.School{}, mapError(err)
	}
	return saved, nil
}
func (s *Repository) DeleteSchool(ctx context.Context, id int64) error {
	result := s.db.WithContext(ctx).Delete(&domain.School{}, id)
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Repository) SchoolCounts(ctx context.Context, id int64) (map[string]int, error) {
	counts := map[string]int{}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var classrooms, teachers, students int64
		if err := tx.Model(&domain.Classroom{}).Where("school_id = ?", id).Count(&classrooms).Error; err != nil {
			return err
		}
		if err := tx.Table("teacher_classrooms AS tc").Joins("JOIN classrooms c ON c.id = tc.classroom_id").Where("c.school_id = ?", id).Distinct("tc.teacher_id").Count(&teachers).Error; err != nil {
			return err
		}
		if err := tx.Model(&domain.Student{}).Joins("JOIN classrooms c ON c.id = students.classroom_id").Where("c.school_id = ?", id).Count(&students).Error; err != nil {
			return err
		}
		counts["classroom_count"] = int(classrooms)
		counts["teacher_count"] = int(teachers)
		counts["student_count"] = int(students)
		return nil
	}, &sql.TxOptions{ReadOnly: true})
	return counts, mapError(err)
}
