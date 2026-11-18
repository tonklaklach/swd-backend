package service

import (
	"context"
	"errors"
	"testing"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/model"
)

type schoolRepositoryStub struct {
	item      domain.School
	saveCalls int
	savedID   int64
	failure   error
}

func (r *schoolRepositoryStub) ListSchool(context.Context, domain.Filter) ([]domain.School, error) {
	return []domain.School{r.item}, r.failure
}
func (r *schoolRepositoryStub) GetSchool(context.Context, int64) (domain.School, error) {
	return r.item, r.failure
}
func (r *schoolRepositoryStub) SaveSchool(_ context.Context, id int64, item domain.School) (domain.School, error) {
	r.saveCalls++
	r.savedID = id
	r.item = item
	if id == 0 {
		id = 1
	}
	r.item.ID = id
	return r.item, r.failure
}
func (r *schoolRepositoryStub) DeleteSchool(context.Context, int64) error { return r.failure }
func (r *schoolRepositoryStub) SchoolCounts(context.Context, int64) (map[string]int, error) {
	return map[string]int{"classroom_count": 2, "teacher_count": 1, "student_count": 3}, r.failure
}
func TestCreateMapsWritableFields(t *testing.T) {
	repo := &schoolRepositoryStub{}
	school, err := NewSchool(repo).Create(context.Background(), model.SchoolRequest{ID: 999, Name: "Swift", Abbreviation: "SWD", Address: "Bangkok"})
	if err != nil {
		t.Fatal(err)
	}
	if repo.saveCalls != 1 || repo.savedID != 0 || school.ID != 1 || school.Name != "Swift" || school.Abbreviation != "SWD" || school.Address != "Bangkok" {
		t.Fatalf("unexpected create result: %+v", school)
	}
}
func TestUpdateUsesExistingIDAndReturnsSavedSchool(t *testing.T) {
	repo := &schoolRepositoryStub{}
	school, err := NewSchool(repo).Update(context.Background(), 7, model.SchoolRequest{Name: "Swift", Abbreviation: "SWD", Address: "Bangkok"})
	if err != nil {
		t.Fatal(err)
	}
	if repo.saveCalls != 1 || repo.savedID != 7 || school.Name != "Swift" || school.ID != 7 {
		t.Fatalf("unexpected update result: %+v", school)
	}
}
func TestGetBuildsSchoolDetail(t *testing.T) {
	repo := &schoolRepositoryStub{item: domain.School{ID: 1, Name: "Swift"}}
	detail, err := NewSchool(repo).Get(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if detail.ID != 1 || detail.Counts["teacher_count"] != 1 {
		t.Fatalf("unexpected detail: %+v", detail)
	}
}
func TestRepositoryErrorsPreserveIdentity(t *testing.T) {
	repo := &schoolRepositoryStub{failure: domain.ErrNotFound}
	_, err := NewSchool(repo).Get(context.Background(), 1)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatal(err)
	}
	repo.failure = domain.ErrInvalidReference
	if err = NewSchool(repo).Delete(context.Background(), 1); !errors.Is(err, domain.ErrConflict) {
		t.Fatal(err)
	}
}
