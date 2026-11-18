package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/database"
)

func setup(t *testing.T) http.Handler {
	t.Helper()
	s, err := database.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close(s) })
	return newApp(s)
}
func request(t *testing.T, h http.Handler, method, path, body string, status int) []byte {
	t.Helper()
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != status {
		t.Fatalf("%s %s: want %d got %d %s", method, path, status, w.Code, w.Body.String())
	}
	if status == http.StatusCreated && w.Header().Get("Location") != "" {
		t.Fatal("create response must not set Location")
	}
	if status == 204 {
		if w.Body.Len() != 0 {
			t.Fatal("204 must have an empty body")
		}
		return nil
	}
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
		t.Fatal(err)
	}
	if status >= 400 {
		var code, message string
		if err := json.Unmarshal(envelope["code"], &code); err != nil || code == "" {
			t.Fatal("error code is required")
		}
		if err := json.Unmarshal(envelope["error"], &message); err != nil || message == "" {
			t.Fatal("error message is required")
		}
		return w.Body.Bytes()
	}
	data, ok := envelope["data"]
	if !ok {
		t.Fatal("success data is required")
	}
	return data
}
func TestSchoolLifecycle(t *testing.T) {
	h := setup(t)
	request(t, h, "POST", "/api/v1/schools", `{"name":"Alpha","abbreviation":"A","address":"Bangkok"}`, 201)
	request(t, h, "POST", "/api/v1/classrooms", `{"school_id":1,"year":1,"section":1}`, 201)
	request(t, h, "POST", "/api/v1/classrooms", `{"school_id":1,"year":1,"section":2}`, 201)
	request(t, h, "POST", "/api/v1/teachers", `{"firstname":"Ann","lastname":"Smith","gender":"female","classroom_ids":[1,2]}`, 201)
	request(t, h, "POST", "/api/v1/students", `{"firstname":"Bob","lastname":"Lee","gender":"male","classroom_id":1}`, 201)
	var detail struct {
		Counts map[string]int `json:"counts"`
	}
	json.Unmarshal(request(t, h, "GET", "/api/v1/schools/1", "", 200), &detail)
	if detail.Counts["teacher_count"] != 1 || detail.Counts["classroom_count"] != 2 || detail.Counts["student_count"] != 1 {
		t.Fatal(detail)
	}
	for _, path := range []string{"classrooms/1", "teachers/1", "students/1"} {
		request(t, h, "GET", "/api/v1/"+path, "", 200)
	}
	var teachers []struct {
		ClassroomIDs []int64 `json:"classroom_ids"`
	}
	json.Unmarshal(request(t, h, "GET", "/api/v1/teachers?school=1", "", 200), &teachers)
	if len(teachers) != 1 || len(teachers[0].ClassroomIDs) != 2 {
		t.Fatal(teachers)
	}
	request(t, h, "PUT", "/api/v1/teachers/1", `{"firstname":"Changed","lastname":"Smith","gender":"female","classroom_ids":[1,999]}`, 422)
	var teacher struct {
		Firstname    string
		ClassroomIDs []int64 `json:"classroom_ids"`
	}
	json.Unmarshal(request(t, h, "GET", "/api/v1/teachers/1", "", 200), &teacher)
	if teacher.Firstname != "Ann" || len(teacher.ClassroomIDs) != 2 {
		t.Fatal("transaction did not roll back", teacher)
	}
	var classroom struct {
		Teachers []json.RawMessage
		Students []json.RawMessage
	}
	json.Unmarshal(request(t, h, "GET", "/api/v1/classrooms/1", "", 200), &classroom)
	if len(classroom.Teachers) != 1 || len(classroom.Students) != 1 {
		t.Fatal("missing classroom members")
	}
	var teacherDetail struct{ Classrooms []json.RawMessage }
	json.Unmarshal(request(t, h, "GET", "/api/v1/teachers/1", "", 200), &teacherDetail)
	if len(teacherDetail.Classrooms) != 2 {
		t.Fatal("missing teacher classrooms")
	}
	var studentDetail struct{ Classroom struct{ ID int64 } }
	json.Unmarshal(request(t, h, "GET", "/api/v1/students/1", "", 200), &studentDetail)
	if studentDetail.Classroom.ID != 1 {
		t.Fatal("missing student classroom")
	}
	for range 2 {
		request(t, h, "PUT", "/api/v1/schools/1", `{"name":"Alpha","abbreviation":"A","address":"Bangkok"}`, 200)
	}
	var schools []json.RawMessage
	json.Unmarshal(request(t, h, "GET", "/api/v1/schools", "", 200), &schools)
	if len(schools) != 1 {
		t.Fatal("PUT created duplicates")
	}
	request(t, h, "DELETE", "/api/v1/schools/1", "", 409)
	request(t, h, "DELETE", "/api/v1/classrooms/1", "", 409)
	request(t, h, "PUT", "/api/v1/schools/1", `{"name":"Beta","abbreviation":"B","address":"Bangkok"}`, 200)
	request(t, h, "PUT", "/api/v1/classrooms/2", `{"school_id":1,"year":2,"section":1}`, 200)
	request(t, h, "PUT", "/api/v1/students/1", `{"firstname":"Bob","lastname":"Lee","gender":"male","classroom_id":2}`, 200)
	request(t, h, "PUT", "/api/v1/teachers/1", `{"firstname":"Ann","lastname":"Smith","gender":"female","classroom_ids":[2]}`, 200)
	request(t, h, "DELETE", "/api/v1/classrooms/1", "", 204)
	request(t, h, "DELETE", "/api/v1/students/1", "", 204)
	request(t, h, "DELETE", "/api/v1/teachers/1", "", 204)
	request(t, h, "DELETE", "/api/v1/classrooms/2", "", 204)
	request(t, h, "DELETE", "/api/v1/schools/1", "", 204)
	request(t, h, "GET", "/api/v1/schools/1", "", 404)
}
func TestInvalidRequests(t *testing.T) {
	h := setup(t)
	for _, tc := range []struct {
		method, path, body string
		status             int
	}{
		{"POST", "schools", `{}`, 422},
		{"POST", "schools", `{"name":"A","abbreviation":"A","address":"B","extra":true}`, 400},
		{"POST", "schools", `{} {}`, 400},
		{"POST", "students", `{"firstname":"A","lastname":"B","gender":"other","classroom_id":999}`, 422},
		{"POST", "classrooms", `{"school_id":999,"year":1,"section":1}`, 422},
		{"POST", "teachers", `{"firstname":"A","lastname":"B","gender":"other","classroom_ids":[1,1]}`, 422},
		{"GET", "missing", "", 404},
		{"PATCH", "schools/1", "{}", 405},
		{"GET", "schools/invalid", "", 400}, {"GET", "schools?limit=101", "", 400},
		{"GET", "schools?unknown=1", "", 400}, {"GET", "teachers?school=no", "", 400},
		{"GET", "students?gender=invalid", "", 400}, {"DELETE", "teachers/999", "", 404},
	} {
		request(t, h, tc.method, "/api/v1/"+tc.path, tc.body, tc.status)
	}
	r := httptest.NewRequest("POST", "/api/v1/schools", strings.NewReader("{}"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 415 {
		t.Fatal(w.Code)
	}
	request(t, h, "GET", "/healthz", "", 200)
}
func TestUniqueClassroomAndPagination(t *testing.T) {
	h := setup(t)
	request(t, h, "POST", "/api/v1/schools", `{"name":"A","abbreviation":"A","address":"B"}`, 201)
	body := `{"school_id":1,"year":1,"section":1}`
	request(t, h, "POST", "/api/v1/classrooms", body, 201)
	request(t, h, "POST", "/api/v1/classrooms", body, 409)
	var empty []any
	json.Unmarshal(request(t, h, "GET", "/api/v1/classrooms?limit=1&offset=1", "", 200), &empty)
	if empty == nil || len(empty) != 0 {
		t.Fatal("expected empty array")
	}
}

func TestPersistence(t *testing.T) {
	dsn := filepath.Join(t.TempDir(), "persistent.db")
	s, err := database.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	request(t, newApp(s), "POST", "/api/v1/schools", `{"name":"Persist","abbreviation":"P","address":"Bangkok"}`, 201)
	if err = database.Close(s); err != nil {
		t.Fatal(err)
	}
	reopened, err := database.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close(reopened)
	body := request(t, newApp(reopened), "GET", "/api/v1/schools/1", "", 200)
	if !strings.Contains(string(body), "Persist") {
		t.Fatal("data was not persisted")
	}
}

func TestUpdateDoesNotCreateMissingResources(t *testing.T) {
	h := setup(t)
	for _, tc := range []struct{ path, body string }{
		{"schools", `{"name":"Swift","abbreviation":"SWD","address":"Bangkok"}`},
		{"classrooms", `{"school_id":1,"year":1,"section":1}`},
		{"teachers", `{"firstname":"Ann","lastname":"Lee","gender":"female","classroom_ids":[]}`},
		{"students", `{"firstname":"Bob","lastname":"Lee","gender":"male","classroom_id":1}`},
	} {
		request(t, h, "PUT", "/api/v1/"+tc.path+"/999", tc.body, 404)
		var items []json.RawMessage
		if err := json.Unmarshal(request(t, h, "GET", "/api/v1/"+tc.path, "", 200), &items); err != nil {
			t.Fatal(err)
		}
		if len(items) != 0 {
			t.Fatalf("PUT created an unexpected %s resource", tc.path)
		}
	}
}

func TestTeacherCanClearClassroomAssignments(t *testing.T) {
	h := setup(t)
	request(t, h, "POST", "/api/v1/schools", `{"name":"Swift","abbreviation":"SWD","address":"Bangkok"}`, 201)
	request(t, h, "POST", "/api/v1/classrooms", `{"school_id":1,"year":1,"section":1}`, 201)
	request(t, h, "POST", "/api/v1/teachers", `{"firstname":"Ann","lastname":"Lee","gender":"female","classroom_ids":[1]}`, 201)
	request(t, h, "PUT", "/api/v1/teachers/1", `{"firstname":"Ann","lastname":"Lee","gender":"female","classroom_ids":[]}`, 200)
	var teacher struct {
		ClassroomIDs []int64 `json:"classroom_ids"`
	}
	if err := json.Unmarshal(request(t, h, "GET", "/api/v1/teachers/1", "", 200), &teacher); err != nil {
		t.Fatal(err)
	}
	if teacher.ClassroomIDs == nil || len(teacher.ClassroomIDs) != 0 {
		t.Fatalf("unexpected classroom IDs: %v", teacher.ClassroomIDs)
	}
	request(t, h, "DELETE", "/api/v1/classrooms/1", "", 204)
}
