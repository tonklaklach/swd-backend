package repository

import (
	"context"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
	"gorm.io/gorm"
)

// teacherRecord excludes classroom IDs, which are stored in the join table.
type teacherRecord struct {
	ID        int64
	Firstname string
	Lastname  string
	Gender    string
}
type teacherClassroom struct {
	TeacherID   int64
	ClassroomID int64
}

func (s *Repository) ListTeacher(ctx context.Context, input domain.Filter) ([]domain.Teacher, error) {
	records := []teacherRecord{}
	err := s.db.WithContext(ctx).Table("teachers").Scopes(buildFilter("teachers", input)).Find(&records).Error
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]domain.Teacher, 0, len(records))
	ids := make([]int64, 0, len(records))
	byID := map[int64]int{}
	for i, record := range records {
		byID[record.ID] = i
		ids = append(ids, record.ID)
		items = append(items, domain.Teacher{ID: record.ID, Firstname: record.Firstname, Lastname: record.Lastname, Gender: record.Gender, ClassroomIDs: []int64{}})
	}
	if len(ids) == 0 {
		return items, nil
	}
	links := []teacherClassroom{}
	err = s.db.WithContext(ctx).Table("teacher_classrooms").Where("teacher_id IN ?", ids).Order("teacher_id, classroom_id").Find(&links).Error
	if err != nil {
		return nil, mapError(err)
	}
	for _, link := range links {
		index := byID[link.TeacherID]
		items[index].ClassroomIDs = append(items[index].ClassroomIDs, link.ClassroomID)
	}
	return items, nil
}
func (s *Repository) GetTeacher(ctx context.Context, id int64) (domain.Teacher, error) {
	var record teacherRecord
	err := s.db.WithContext(ctx).Table("teachers").First(&record, id).Error
	if err != nil {
		return domain.Teacher{}, mapError(err)
	}
	ids := []int64{}
	err = s.db.WithContext(ctx).Table("teacher_classrooms").Where("teacher_id = ?", id).Order("classroom_id").Pluck("classroom_id", &ids).Error
	return domain.Teacher{ID: record.ID, Firstname: record.Firstname, Lastname: record.Lastname, Gender: record.Gender, ClassroomIDs: ids}, mapError(err)
}
func (s *Repository) SaveTeacher(ctx context.Context, id int64, item domain.Teacher) (domain.Teacher, error) {
	record := teacherRecord{ID: id, Firstname: item.Firstname, Lastname: item.Lastname, Gender: item.Gender}
	var saved domain.Teacher
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if id == 0 {
			if err := tx.Table("teachers").Create(&record).Error; err != nil {
				return err
			}
		} else {
			result := tx.Table("teachers").Where("id = ?", id).Select("firstname", "lastname", "gender").Updates(&record)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return domain.ErrNotFound
			}
		}
		if err := tx.Where("teacher_id = ?", record.ID).Delete(&teacherClassroom{}).Error; err != nil {
			return err
		}
		for _, classroomID := range item.ClassroomIDs {
			link := teacherClassroom{TeacherID: record.ID, ClassroomID: classroomID}
			if err := tx.Create(&link).Error; err != nil {
				return err
			}
		}
		var err error
		saved, err = New(tx).GetTeacher(ctx, record.ID)
		return err
	})
	if err != nil {
		return domain.Teacher{}, mapError(err)
	}
	return saved, nil
}
func (s *Repository) DeleteTeacher(ctx context.Context, id int64) error {
	result := s.db.WithContext(ctx).Table("teachers").Delete(&teacherRecord{}, id)
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
