package models

import "time"

type Base struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type Student struct {
	Base
	Person
	GroupID     string    `json:"group_id,omitempty"`
	EnrolledAt  time.Time `json:"enrolled_at"`
	IsActive    bool      `json:"is_active"`
	StudentCard string    `json:"student_card"`
}

type Group struct {
	Base
	Name        string    `json:"name"`
	CourseID    string    `json:"course_id"`
	StudentIDs  []string  `json:"student_ids"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"`
	MaxStudents int       `json:"max_students"`
}

type Lesson struct {
	Base
	GroupID     string    `json:"group_id"`
	Topic       string    `json:"topic"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Room        string    `json:"room"`
	TeacherID   string    `json:"teacher_id"`
	Status      string    `json:"status"`
}

type Attendance struct {
	LessonID  string    `json:"lesson_id"`
	StudentID string    `json:"student_id"`
	Present   bool      `json:"present"`
	MarkedAt  time.Time `json:"marked_at"`
	MarkedBy  string    `json:"marked_by"`
}

type Grade struct {
	StudentID string    `json:"student_id"`
	LessonID  string    `json:"lesson_id"`
	Value     int       `json:"value"`
	Comment   string    `json:"comment"`
	GradedAt  time.Time `json:"graded_at"`
	GradedBy  string    `json:"graded_by"`
	Type      string    `json:"type"` // exam, test, homework
}

type GradeBook struct {
	GroupID     string              `json:"group_id"`
	CourseName  string              `json:"course_name"`
	Lessons     []*Lesson           `json:"lessons"`
	Students    []*Student          `json:"students"`
	Grades      map[string][]*Grade `json:"grades"`
	GeneratedAt time.Time           `json:"generated_at"`
}
