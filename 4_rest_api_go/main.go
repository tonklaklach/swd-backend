package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/database"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/handler"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/repository"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/response"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/service"

	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
func run() error {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "school.db"
	}
	db, err := database.Open(dsn)
	if err != nil {
		return err
	}
	defer database.Close(db)
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	server := &http.Server{Addr: addr, Handler: newApp(db), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	result := make(chan error, 1)
	go func() { slog.Info("server listening", "address", addr); result <- server.ListenAndServe() }()
	select {
	case err := <-result:
		return err
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdown)
	}
}

func newApp(db *gorm.DB) *echo.Echo {
	app := echo.New()
	app.HTTPErrorHandler = response.ErrorHandler
	repo := repository.New(db)
	school := handler.NewSchool(service.NewSchool(repo))
	classroom := handler.NewClassroom(service.NewClassroom(repo, repo, repo))
	teacher := handler.NewTeacher(service.NewTeacher(repo, repo))
	student := handler.NewStudent(service.NewStudent(repo, repo))
	app.GET("/healthz", func(c *echo.Context) error {
		if err := db.WithContext(c.Request().Context()).Exec("SELECT 1").Error; err != nil {
			return echo.NewHTTPError(503, "database unavailable")
		}
		return response.Success(c, 200, map[string]string{"status": "ok"})
	})
	app.GET("/api/v1/schools", school.List)
	app.POST("/api/v1/schools", school.Create)
	app.GET("/api/v1/schools/:id", school.Get)
	app.PUT("/api/v1/schools/:id", school.Update)
	app.DELETE("/api/v1/schools/:id", school.Delete)
	app.GET("/api/v1/classrooms", classroom.List)
	app.POST("/api/v1/classrooms", classroom.Create)
	app.GET("/api/v1/classrooms/:id", classroom.Get)
	app.PUT("/api/v1/classrooms/:id", classroom.Update)
	app.DELETE("/api/v1/classrooms/:id", classroom.Delete)
	app.GET("/api/v1/teachers", teacher.List)
	app.POST("/api/v1/teachers", teacher.Create)
	app.GET("/api/v1/teachers/:id", teacher.Get)
	app.PUT("/api/v1/teachers/:id", teacher.Update)
	app.DELETE("/api/v1/teachers/:id", teacher.Delete)
	app.GET("/api/v1/students", student.List)
	app.POST("/api/v1/students", student.Create)
	app.GET("/api/v1/students/:id", student.Get)
	app.PUT("/api/v1/students/:id", student.Update)
	app.DELETE("/api/v1/students/:id", student.Delete)
	return app
}
