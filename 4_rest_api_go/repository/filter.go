package repository

import (
	"strings"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
	"gorm.io/gorm"
)

func buildFilter(resource string, f domain.Filter) func(*gorm.DB) *gorm.DB {
	var args []any
	clauses := []string{}
	add := func(query string, arg any) { clauses = append(clauses, query); args = append(args, arg) }
	if f.Name != "" {
		add("instr(lower(name),lower(?))>0", f.Name)
	}
	if f.Firstname != "" {
		add("instr(lower(firstname),lower(?))>0", f.Firstname)
	}
	if f.Lastname != "" {
		add("instr(lower(lastname),lower(?))>0", f.Lastname)
	}
	if f.Gender != "" {
		add("gender=?", f.Gender)
	}
	if f.SchoolID > 0 {
		switch resource {
		case "teachers":
			add("EXISTS (SELECT 1 FROM teacher_classrooms tc JOIN classrooms c ON c.id=tc.classroom_id WHERE tc.teacher_id=teachers.id AND c.school_id=?)", f.SchoolID)
		case "students":
			add("classroom_id IN (SELECT id FROM classrooms WHERE school_id=?)", f.SchoolID)
		default:
			add("school_id=?", f.SchoolID)
		}
	}
	if f.ClassroomID > 0 {
		if resource == "teachers" {
			add("EXISTS (SELECT 1 FROM teacher_classrooms tc WHERE tc.teacher_id=teachers.id AND tc.classroom_id=?)", f.ClassroomID)
		} else {
			add("classroom_id=?", f.ClassroomID)
		}
	}
	if f.TeacherID > 0 {
		add("id IN (SELECT classroom_id FROM teacher_classrooms WHERE teacher_id=?)", f.TeacherID)
	}
	return func(db *gorm.DB) *gorm.DB {
		if len(clauses) > 0 {
			db = db.Where(strings.Join(clauses, " AND "), args...)
		}
		return db.Order("id").Limit(f.Limit).Offset(f.Offset)
	}
}
