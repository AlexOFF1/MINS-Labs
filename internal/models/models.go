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

type Teacher struct {
	Base
	Person
	Specialization string   `json:"specialization"`
	Courses        []string `json:"course_ids"`
}

type Course struct {
	Base
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	TeacherID     string  `json:"teacher_id"`
	DurationWeeks int     `json:"duration_weeks"`
	Price         float64 `json:"price"`
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
}

type Report struct {
	StudentID    string       `json:"student_id"`
	GroupID      string       `json:"group_id"`
	PeriodStart  time.Time    `json:"period_start"`
	PeriodEnd    time.Time    `json:"period_end"`
	Grades       []Grade      `json:"grades"`
	Attendance   []Attendance `json:"attendance"`
	AverageGrade float64      `json:"average_grade"`
	TotalHours   int          `json:"total_hours"`
	MissedHours  int          `json:"missed_hours"`
}
