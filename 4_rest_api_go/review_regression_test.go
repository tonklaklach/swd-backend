package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"slices"
	"testing"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/database"
	"gorm.io/gorm"
)

func TestListFiltersReturnMatchingIDs(t *testing.T) {
	h := setup(t)
	fixtures := []struct{ resource, body string }{
		{"schools", `{"name":"Alpha","abbreviation":"A","address":"Bangkok"}`},
		{"schools", `{"name":"Beta","abbreviation":"B","address":"Bangkok"}`},
		{"classrooms", `{"school_id":1,"year":1,"section":1}`},
		{"classrooms", `{"school_id":1,"year":1,"section":2}`},
		{"classrooms", `{"school_id":2,"year":1,"section":1}`},
		{"teachers", `{"firstname":"Ann","lastname":"Smith","gender":"female","classroom_ids":[1,2]}`},
		{"teachers", `{"firstname":"Andy","lastname":"Lee","gender":"male","classroom_ids":[3]}`},
		{"teachers", `{"firstname":"Cara","lastname":"Smith","gender":"female","classroom_ids":[2]}`},
		{"students", `{"firstname":"Bob","lastname":"Lee","gender":"male","classroom_id":1}`},
		{"students", `{"firstname":"Bobby","lastname":"Smith","gender":"female","classroom_id":3}`},
		{"students", `{"firstname":"Cara","lastname":"Lee","gender":"male","classroom_id":2}`},
	}
	for _, fixture := range fixtures {
		request(t, h, http.MethodPost, "/api/v1/"+fixture.resource, fixture.body, http.StatusCreated)
	}
	cases := []struct {
		path string
		want []int64
	}{
		{"schools?name=ALP", []int64{1}},
		{"classrooms?school=1", []int64{1, 2}},
		{"teachers?school=1", []int64{1, 3}},
		{"teachers?classroom=2", []int64{1, 3}},
		{"teachers?firstname=AN", []int64{1, 2}},
		{"teachers?lastname=SMI", []int64{1, 3}},
		{"teachers?gender=male", []int64{2}},
		{"teachers?school=1&classroom=2&firstname=AN&lastname=SMI&gender=female", []int64{1}},
		{"teachers?school=1&limit=1&offset=1", []int64{3}},
		{"students?school=1", []int64{1, 3}},
		{"students?classroom=2", []int64{3}},
		{"students?firstname=BOB", []int64{1, 2}},
		{"students?lastname=LEE", []int64{1, 3}},
		{"students?gender=female", []int64{2}},
		{"students?school=1&classroom=2&firstname=CAR&lastname=LEE&gender=male", []int64{3}},
		{"students?school=2&gender=male", nil},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			payload := request(t, h, http.MethodGet, "/api/v1/"+tc.path, "", http.StatusOK)
			var items []struct{ ID int64 }
			if err := json.Unmarshal(payload, &items); err != nil {
				t.Fatal(err)
			}
			if items == nil {
				t.Fatalf("expected a JSON array, got %s", payload)
			}
			ids := make([]int64, 0, len(items))
			for _, item := range items {
				ids = append(ids, item.ID)
			}
			if !slices.Equal(ids, tc.want) {
				t.Fatalf("got IDs %v, want %v", ids, tc.want)
			}
		})
	}
}

func TestTeacherWriteRollsBackWhenResponseReadFails(t *testing.T) {
	db, err := database.Open(filepath.Join(t.TempDir(), "rollback.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = database.Close(db) })
	h := newApp(db)
	request(t, h, http.MethodPost, "/api/v1/schools", `{"name":"Alpha","abbreviation":"A","address":"Bangkok"}`, 201)
	request(t, h, http.MethodPost, "/api/v1/classrooms", `{"school_id":1,"year":1,"section":1}`, 201)
	request(t, h, http.MethodPost, "/api/v1/classrooms", `{"school_id":1,"year":1,"section":2}`, 201)
	request(t, h, http.MethodPost, "/api/v1/teachers", `{"firstname":"Ann","lastname":"Smith","gender":"female","classroom_ids":[1]}`, 201)

	const callback = "test:fail_teacher_assignment_read"
	err = db.Callback().Query().Before("gorm:query").Register(callback, func(tx *gorm.DB) {
		if tx.Statement.Table == "teacher_classrooms" {
			tx.AddError(errors.New("could not read saved classroom assignments"))
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	body := `{"firstname":"Changed","lastname":"Smith","gender":"female","classroom_ids":[2]}`
	request(t, h, http.MethodPost, "/api/v1/teachers", body, 500)
	request(t, h, http.MethodPut, "/api/v1/teachers/1", body, 500)
	if err = db.Callback().Query().Remove(callback); err != nil {
		t.Fatal(err)
	}

	var teachers []struct {
		ID           int64
		Firstname    string
		ClassroomIDs []int64 `json:"classroom_ids"`
	}
	if err = json.Unmarshal(request(t, h, http.MethodGet, "/api/v1/teachers", "", 200), &teachers); err != nil {
		t.Fatal(err)
	}
	if len(teachers) != 1 || teachers[0].ID != 1 || teachers[0].Firstname != "Ann" || !slices.Equal(teachers[0].ClassroomIDs, []int64{1}) {
		t.Fatalf("failed create or update changed stored teachers: %+v", teachers)
	}
}
