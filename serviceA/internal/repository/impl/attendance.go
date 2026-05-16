package impl

import (
	"context"
	liberrors "errors"
	"mins_EduCenter/pkg/errors"
	"mins_EduCenter/serviceA/internal/models"
	"mins_EduCenter/serviceA/internal/repository"
	"sync"
	"time"
)

type attendanceRepository struct {
	mu         sync.RWMutex
	store      map[string][]*models.Attendance
	studentIdx map[string][]string
}

func NewAttendanceRepository() repository.AttendanceRepository {
	return &attendanceRepository{
		store:      make(map[string][]*models.Attendance),
		studentIdx: make(map[string][]string),
	}
}

func (r *attendanceRepository) Mark(ctx context.Context, attendance *models.Attendance) error {
	const op = "attendanceRepository.Mark"
	r.mu.Lock()
	defer r.mu.Unlock()

	if attendance.LessonID == "" || attendance.StudentID == "" {
		return errors.NewValidationError(op, "attendance", "lessonID and studentID required")
	}
	attendance.MarkedAt = time.Now()

	existing, ok := r.store[attendance.LessonID]
	if ok {
		for i, a := range existing {
			if a.StudentID == attendance.StudentID {
				existing[i] = attendance
				r.updateStudentIndex(attendance.StudentID, attendance.LessonID)
				return nil
			}
		}
	}
	r.store[attendance.LessonID] = append(r.store[attendance.LessonID], attendance)
	r.updateStudentIndex(attendance.StudentID, attendance.LessonID)
	return nil
}

func (r *attendanceRepository) MarkBatch(ctx context.Context, attendanceList []*models.Attendance) error {
	for _, a := range attendanceList {
		if err := r.Mark(ctx, a); err != nil {
			return err
		}
	}
	return nil
}

func (r *attendanceRepository) GetByLesson(ctx context.Context, lessonID string) ([]*models.Attendance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	attendance, ok := r.store[lessonID]
	if !ok {
		return []*models.Attendance{}, nil
	}
	result := make([]*models.Attendance, len(attendance))
	copy(result, attendance)
	return result, nil
}

func (r *attendanceRepository) GetByStudent(ctx context.Context, studentID string) ([]*models.Attendance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lessonIDs, ok := r.studentIdx[studentID]
	if !ok {
		return []*models.Attendance{}, nil
	}
	var result []*models.Attendance
	for _, lid := range lessonIDs {
		if list, ok := r.store[lid]; ok {
			for _, a := range list {
				if a.StudentID == studentID {
					result = append(result, a)
				}
			}
		}
	}
	return result, nil
}

func (r *attendanceRepository) GetByDateRange(ctx context.Context, from, to time.Time) ([]*models.Attendance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Attendance
	for _, list := range r.store {
		for _, a := range list {
			if (a.MarkedAt.After(from) || a.MarkedAt.Equal(from)) &&
				(a.MarkedAt.Before(to) || a.MarkedAt.Equal(to)) {
				result = append(result, a)
			}
		}
	}
	return result, nil
}

func (r *attendanceRepository) GetAttendanceStats(ctx context.Context, groupID string) (map[string]float64, error) {
	return nil, errors.NewInternalError("attendanceRepository.GetAttendanceStats",
		liberrors.New("use LessonUsecase for attendance stats"))
}

func (r *attendanceRepository) updateStudentIndex(studentID, lessonID string) {
	for _, lid := range r.studentIdx[studentID] {
		if lid == lessonID {
			return
		}
	}
	r.studentIdx[studentID] = append(r.studentIdx[studentID], lessonID)
}
