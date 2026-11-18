package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/model"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/response"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/service"
)

type schoolServiceSpy struct {
	service.SchoolService
	calls int
	ctx   context.Context
	req   model.SchoolRequest
}

func (s *schoolServiceSpy) Create(ctx context.Context, req model.SchoolRequest) (domain.School, error) {
	s.calls++
	s.ctx = ctx
	s.req = req
	return domain.School{ID: 1, Name: req.Name, Abbreviation: req.Abbreviation, Address: req.Address}, nil
}
func TestCreateBindsValidatesAndPassesRequestContext(t *testing.T) {
	for _, tc := range []struct {
		name, body    string
		status, calls int
	}{
		{"malformed", `{`, http.StatusBadRequest, 0},
		{"invalid", `{"name":" "}`, http.StatusUnprocessableEntity, 0},
		{"valid", `{"name":"Swift","abbreviation":"SWD","address":"Bangkok"}`, http.StatusCreated, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			spy := &schoolServiceSpy{}
			app := echo.New()
			app.HTTPErrorHandler = response.ErrorHandler
			app.POST("/schools", NewSchool(spy).Create)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			req := httptest.NewRequest(http.MethodPost, "/schools", strings.NewReader(tc.body)).WithContext(ctx)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			app.ServeHTTP(recorder, req)
			if recorder.Code != tc.status || spy.calls != tc.calls {
				t.Fatalf("status=%d calls=%d body=%s", recorder.Code, spy.calls, recorder.Body.String())
			}
			if tc.calls > 0 && (spy.ctx != ctx || spy.req.Name != "Swift") {
				t.Fatal("request context or bound DTO was not forwarded")
			}
		})
	}
}
