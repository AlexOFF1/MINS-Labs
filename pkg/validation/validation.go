package validation

import (
	"fmt"
	"regexp"
	"time"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("field '%s': %s", e.Field, e.Message)
}

type Validator interface {
	Validate() []ValidationError
}

type StudentRegisterData struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

func (s *StudentRegisterData) Validate() []ValidationError {
	var errors []ValidationError

	if s.FirstName == "" {
		errors = append(errors, ValidationError{Field: "FirstName", Message: "required"})
	}
	if s.LastName == "" {
		errors = append(errors, ValidationError{Field: "LastName", Message: "required"})
	}
	if s.Email == "" {
		errors = append(errors, ValidationError{Field: "Email", Message: "required"})
	} else {
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
		if !emailRegex.MatchString(s.Email) {
			errors = append(errors, ValidationError{Field: "Email", Message: "invalid format"})
		}
	}

	return errors
}

type GroupCreateData struct {
	Name        string
	MaxStudents int
	StartDate   time.Time
	EndDate     time.Time
}

func (g *GroupCreateData) Validate() []ValidationError {
	var errors []ValidationError

	if g.Name == "" {
		errors = append(errors, ValidationError{Field: "Name", Message: "required"})
	}
	if g.MaxStudents <= 0 {
		errors = append(errors, ValidationError{Field: "MaxStudents", Message: "must be greater than 0"})
	}
	if g.MaxStudents > 100 {
		errors = append(errors, ValidationError{Field: "MaxStudents", Message: "cannot exceed 100"})
	}
	if !g.StartDate.IsZero() && !g.EndDate.IsZero() && g.StartDate.After(g.EndDate) {
		errors = append(errors, ValidationError{Field: "StartDate", Message: "must be before EndDate"})
	}

	return errors
}

type LessonCreateData struct {
	GroupID   string
	Topic     string
	StartTime time.Time
	EndTime   time.Time
	Room      string
}

func (l *LessonCreateData) Validate() []ValidationError {
	var errors []ValidationError

	if l.GroupID == "" {
		errors = append(errors, ValidationError{Field: "GroupID", Message: "required"})
	}
	if l.Topic == "" {
		errors = append(errors, ValidationError{Field: "Topic", Message: "required"})
	}
	if l.Room == "" {
		errors = append(errors, ValidationError{Field: "Room", Message: "required"})
	}
	if l.StartTime.IsZero() {
		errors = append(errors, ValidationError{Field: "StartTime", Message: "required"})
	}
	if l.EndTime.IsZero() {
		errors = append(errors, ValidationError{Field: "EndTime", Message: "required"})
	}
	if !l.StartTime.IsZero() && !l.EndTime.IsZero() {
		if l.StartTime.After(l.EndTime) {
			errors = append(errors, ValidationError{Field: "EndTime", Message: "must be after StartTime"})
		}
		if l.StartTime.Before(time.Now()) {
			errors = append(errors, ValidationError{Field: "StartTime", Message: "cannot be in past"})
		}
	}

	return errors
}

type AttendanceMarkData struct {
	LessonID  string
	StudentID string
	Present   bool
}

func (a *AttendanceMarkData) Validate() []ValidationError {
	var errors []ValidationError

	if a.LessonID == "" {
		errors = append(errors, ValidationError{Field: "LessonID", Message: "required"})
	}
	if a.StudentID == "" {
		errors = append(errors, ValidationError{Field: "StudentID", Message: "required"})
	}

	return errors
}

type GradeSetData struct {
	StudentID string
	LessonID  string
	Value     int
	Type      string
}

func (g *GradeSetData) Validate() []ValidationError {
	var errors []ValidationError

	if g.StudentID == "" {
		errors = append(errors, ValidationError{Field: "StudentID", Message: "required"})
	}
	if g.LessonID == "" {
		errors = append(errors, ValidationError{Field: "LessonID", Message: "required"})
	}
	if g.Value < 1 || g.Value > 5 {
		errors = append(errors, ValidationError{Field: "Value", Message: "must be between 1 and 5"})
	}
	if g.Type == "" {
		errors = append(errors, ValidationError{Field: "Type", Message: "required (exam, test, homework)"})
	}
	if g.Type != "" && g.Type != "exam" && g.Type != "test" && g.Type != "homework" {
		errors = append(errors, ValidationError{Field: "Type", Message: "must be exam, test, or homework"})
	}

	return errors
}

func Validate(v Validator) error {
	errors := v.Validate()
	if len(errors) == 0 {
		return nil
	}

	first := errors[0]
	return &ValidationError{
		Field:   first.Field,
		Message: first.Message,
	}
}
