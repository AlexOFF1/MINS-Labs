package repository

import (
	"context"
	"mins_EduCenter/internal/models"
	"time"
)

type StudentRepository interface {
	Create(ctx context.Context, student *models.Student) error
	GetByID(ctx context.Context, id string) (*models.Student, error)
	GetAll(ctx context.Context) ([]*models.Student, error)
	Update(ctx context.Context, student *models.Student) error
	Delete(ctx context.Context, id string) error
	GetByGroup(ctx context.Context, groupID string) ([]*models.Student, error)
	GetActive(ctx context.Context) ([]*models.Student, error)
	Search(ctx context.Context, query string) ([]*models.Student, error)
}

type GroupRepository interface {
	Create(ctx context.Context, group *models.Group) error
	GetByID(ctx context.Context, id string) (*models.Group, error)
	GetAll(ctx context.Context) ([]*models.Group, error)
	Update(ctx context.Context, group *models.Group) error
	Delete(ctx context.Context, id string) error
	AddStudent(ctx context.Context, groupID, studentID string) error
	RemoveStudent(ctx context.Context, groupID, studentID string) error
	GetByCourse(ctx context.Context, courseID string) ([]*models.Group, error)
}

type LessonRepository interface {
	Create(ctx context.Context, lesson *models.Lesson) error
	GetByID(ctx context.Context, id string) (*models.Lesson, error)
	GetByGroup(ctx context.Context, groupID string) ([]*models.Lesson, error)
	GetByDateRange(ctx context.Context, from, to time.Time) ([]*models.Lesson, error)
	Update(ctx context.Context, lesson *models.Lesson) error
	Delete(ctx context.Context, id string) error
	GetScheduleForStudent(ctx context.Context, studentID string) ([]*models.Lesson, error)
	GetUpcomingLessons(ctx context.Context, groupID string, limit int) ([]*models.Lesson, error)
}

type AttendanceRepository interface {
	Mark(ctx context.Context, attendance *models.Attendance) error
	GetByLesson(ctx context.Context, lessonID string) ([]*models.Attendance, error)
	GetByStudent(ctx context.Context, studentID string) ([]*models.Attendance, error)
	GetByDateRange(ctx context.Context, from, to time.Time) ([]*models.Attendance, error)
	GetAttendanceStats(ctx context.Context, groupID string) (map[string]float64, error)
	MarkBatch(ctx context.Context, attendanceList []*models.Attendance) error
}

type GradeRepository interface {
	Set(ctx context.Context, grade *models.Grade) error
	GetByStudent(ctx context.Context, studentID string) ([]*models.Grade, error)
	GetByLesson(ctx context.Context, lessonID string) ([]*models.Grade, error)
	GetAverageForStudent(ctx context.Context, studentID string) (float64, error)
	GetGradeBook(ctx context.Context, groupID string) (*models.GradeBook, error)
	UpdateGrade(ctx context.Context, studentID, lessonID string, value int) error
}
