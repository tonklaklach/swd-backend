package utils

import (
	"errors"
	"testing"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/model"
)

func TestRequestValidation(t *testing.T) {
	for _, tc := range []struct {
		name  string
		req   any
		valid bool
	}{
		{"school valid", model.SchoolRequest{Name: "Swift", Abbreviation: "SWD", Address: "Bangkok"}, true},
		{"school missing", model.SchoolRequest{}, false},
		{"school whitespace", model.SchoolRequest{Name: " \t\n", Abbreviation: "SWD", Address: "Bangkok"}, false},
		{"classroom valid", model.ClassroomRequest{SchoolID: 1, Year: 1, Section: 1}, true},
		{"classroom zero", model.ClassroomRequest{SchoolID: 1, Year: 1}, false},
		{"teacher unassigned", model.TeacherRequest{Firstname: "Ann", Lastname: "Lee", Gender: "female"}, true},
		{"teacher assigned", model.TeacherRequest{Firstname: "Ann", Lastname: "Lee", Gender: "female", ClassroomIDs: []int64{1, 2}}, true},
		{"teacher duplicate", model.TeacherRequest{Firstname: "Ann", Lastname: "Lee", Gender: "female", ClassroomIDs: []int64{1, 1}}, false},
		{"teacher invalid ID", model.TeacherRequest{Firstname: "Ann", Lastname: "Lee", Gender: "female", ClassroomIDs: []int64{0}}, false},
		{"teacher whitespace", model.TeacherRequest{Firstname: " ", Lastname: "Lee", Gender: "female"}, false},
		{"student valid", model.StudentRequest{Firstname: "Bob", Lastname: "Lee", Gender: "other", ClassroomID: 1}, true},
		{"student gender", model.StudentRequest{Firstname: "Bob", Lastname: "Lee", Gender: "invalid", ClassroomID: 1}, false},
		{"student classroom", model.StudentRequest{Firstname: "Bob", Lastname: "Lee", Gender: "male", ClassroomID: -1}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.req)
			if tc.valid {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			var invalid *domain.ValidationError
			if !errors.As(err, &invalid) {
				t.Fatalf("expected validation error, got %v", err)
			}
			if invalid.Message == "" {
				t.Fatal("missing validation message")
			}
		})
	}
}
