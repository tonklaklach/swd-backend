package repository

import "github.com/tonklaklach/swd-backend/4_rest_api_go/service"

var (
	_ service.SchoolRepository    = (*Repository)(nil)
	_ service.ClassroomRepository = (*Repository)(nil)
	_ service.TeacherRepository   = (*Repository)(nil)
	_ service.StudentRepository   = (*Repository)(nil)
)
