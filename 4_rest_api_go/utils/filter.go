package utils

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
)

func ParseFilter(resource string, q url.Values) (domain.Filter, error) {
	f := domain.Filter{Limit: 50}
	allowed := map[string]bool{"limit": true, "offset": true}
	switch resource {
	case "schools":
		allowed["name"] = true
	case "classrooms":
		allowed["school"] = true
	default:
		for _, key := range []string{"school", "classroom", "firstname", "lastname", "gender"} {
			allowed[key] = true
		}
	}
	for key, values := range q {
		if !allowed[key] {
			return f, fmt.Errorf("unknown filter %s", key)
		}
		if len(values) != 1 {
			return f, errors.New("query parameters must appear once")
		}
		value := values[0]
		switch key {
		case "limit", "offset":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 || (key == "limit" && (n < 1 || n > 100)) {
				return f, errors.New("limit must be 1..100 and offset must be nonnegative")
			}
			if key == "limit" {
				f.Limit = n
			} else {
				f.Offset = n
			}
		case "school", "classroom":
			n, err := strconv.ParseInt(value, 10, 64)
			if err != nil || n < 1 {
				return f, fmt.Errorf("%s must be positive", key)
			}
			if key == "school" {
				f.SchoolID = n
			} else {
				f.ClassroomID = n
			}
		case "gender":
			if value != "male" && value != "female" && value != "other" {
				return f, errors.New("invalid gender")
			}
			f.Gender = value
		case "name":
			f.Name = value
		case "firstname":
			f.Firstname = value
		case "lastname":
			f.Lastname = value
		}
	}
	return f, nil
}
