package usecase

import (
	"context"
	"mins_EduCenter/pkg/errors"
	"mins_EduCenter/pkg/validation"
	"mins_EduCenter/serviceA/internal/models"
	"mins_EduCenter/serviceA/internal/repository"
	"time"
)

type LessonUsecase struct {
	lessonRepo     repository.LessonRepository
	attendanceRepo repository.AttendanceRepository
	groupRepo      repository.GroupRepository
	studentRepo    repository.StudentRepository
}

func NewLessonUsecase(
	lr repository.LessonRepository,
	ar repository.AttendanceRepository,
	gr repository.GroupRepository,
	sr repository.StudentRepository,
) *LessonUsecase {
	return &LessonUsecase{
		lessonRepo:     lr,
		attendanceRepo: ar,
		groupRepo:      gr,
		studentRepo:    sr,
	}
}

type CreateLessonDTO struct {
	GroupID     string
	Topic       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	Room        string
	TeacherID   string
}

func (u *LessonUsecase) CreateLesson(ctx context.Context, dto CreateLessonDTO) (*models.Lesson, error) {
	const op = "LessonUsecase.CreateLesson"

	lessonData := &validation.LessonCreateData{
		GroupID:   dto.GroupID,
		Topic:     dto.Topic,
		StartTime: dto.StartTime,
		EndTime:   dto.EndTime,
		Room:      dto.Room,
	}

	if err := validation.Validate(lessonData); err != nil {
		if valErr, ok := err.(*validation.ValidationError); ok {
			return nil, errors.NewValidationError(op, valErr.Field, valErr.Message)
		}
		return nil, errors.NewValidationError(op, "unknown", err.Error())
	}

	if _, err := u.groupRepo.GetByID(ctx, dto.GroupID); err != nil {
		return nil, errors.NewValidationError(op, "GroupID", "group not found")
	}

	lesson := &models.Lesson{
		GroupID:     dto.GroupID,
		Topic:       dto.Topic,
		Description: dto.Description,
		StartTime:   dto.StartTime,
		EndTime:     dto.EndTime,
		Room:        dto.Room,
		TeacherID:   dto.TeacherID,
		Status:      "scheduled",
	}

	if err := u.lessonRepo.Create(ctx, lesson); err != nil {
		return nil, errors.NewInternalError(op, err)
	}
	return lesson, nil
}

func (u *LessonUsecase) MarkAttendance(ctx context.Context, dto MarkAttendanceDTO) error {
	const op = "LessonUsecase.MarkAttendance"

	attendanceData := &validation.AttendanceMarkData{
		LessonID:  dto.LessonID,
		StudentID: dto.StudentID,
		Present:   dto.Present,
	}

	if err := validation.Validate(attendanceData); err != nil {
		if valErr, ok := err.(*validation.ValidationError); ok {
			return errors.NewValidationError(op, valErr.Field, valErr.Message)
		}
		return errors.NewValidationError(op, "unknown", err.Error())
	}

	lesson, err := u.lessonRepo.GetByID(ctx, dto.LessonID)
	if err != nil {
		return errors.NewValidationError(op, "LessonID", "lesson not found")
	}

	if lesson.StartTime.After(time.Now()) {
		return errors.NewValidationError(op, "LessonID", "cannot mark attendance for future lesson")
	}

	attendance := &models.Attendance{
		LessonID:  dto.LessonID,
		StudentID: dto.StudentID,
		Present:   dto.Present,
		MarkedBy:  dto.MarkedBy,
	}

	if err := u.attendanceRepo.Mark(ctx, attendance); err != nil {
		return errors.NewInternalError(op, err)
	}
	return nil
}

func (u *LessonUsecase) GetGroupSchedule(ctx context.Context, groupID string) ([]*models.Lesson, error) {
	const op = "LessonUsecase.GetGroupSchedule"
	lessons, err := u.lessonRepo.GetByGroup(ctx, groupID)
	if err != nil {
		return nil, errors.NewInternalError(op, err)
	}
	return lessons, nil
}

func (u *LessonUsecase) GetStudentSchedule(ctx context.Context, studentID string) ([]*models.Lesson, error) {
	const op = "LessonUsecase.GetStudentSchedule"
	student, err := u.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return nil, errors.NewValidationError(op, "studentID", "student not found")
	}
	if student.GroupID == "" {
		return []*models.Lesson{}, nil
	}
	return u.GetGroupSchedule(ctx, student.GroupID)
}

type MarkAttendanceDTO struct {
	LessonID  string
	StudentID string
	Present   bool
	MarkedBy  string
}

func (u *LessonUsecase) MarkBatchAttendance(ctx context.Context, lessonID string, attendanceMap map[string]bool, markedBy string) error {
	const op = "LessonUsecase.MarkBatchAttendance"
	lesson, err := u.lessonRepo.GetByID(ctx, lessonID)
	if err != nil {
		return errors.NewValidationError(op, "LessonID", "lesson not found")
	}
	students, err := u.studentRepo.GetByGroup(ctx, lesson.GroupID)
	if err != nil {
		return errors.NewInternalError(op, err)
	}
	var list []*models.Attendance
	for _, student := range students {
		present := attendanceMap[student.ID]
		list = append(list, &models.Attendance{
			LessonID:  lessonID,
			StudentID: student.ID,
			Present:   present,
			MarkedBy:  markedBy,
		})
	}
	if err := u.attendanceRepo.MarkBatch(ctx, list); err != nil {
		return errors.NewInternalError(op, err)
	}
	return nil
}

func (u *LessonUsecase) GetLessonAttendance(ctx context.Context, lessonID string) ([]*models.Attendance, error) {
	const op = "LessonUsecase.GetLessonAttendance"
	attendance, err := u.attendanceRepo.GetByLesson(ctx, lessonID)
	if err != nil {
		return nil, errors.NewInternalError(op, err)
	}
	return attendance, nil
}
