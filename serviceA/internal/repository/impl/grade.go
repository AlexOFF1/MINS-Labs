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

type gradeRepository struct {
	mu        sync.RWMutex
	store     map[string][]*models.Grade
	lessonMap map[string][]string
}

func NewGradeRepository() repository.GradeRepository {
	return &gradeRepository{
		store:     make(map[string][]*models.Grade),
		lessonMap: make(map[string][]string),
	}
}

func (r *gradeRepository) Set(ctx context.Context, grade *models.Grade) error {
	const op = "gradeRepository.Set"
	r.mu.Lock()
	defer r.mu.Unlock()

	if grade.StudentID == "" || grade.LessonID == "" {
		return errors.NewValidationError(op, "grade", "studentID and lessonID required")
	}
	grade.GradedAt = time.Now()

	r.store[grade.StudentID] = append(r.store[grade.StudentID], grade)
	r.lessonMap[grade.LessonID] = append(r.lessonMap[grade.LessonID], grade.StudentID)
	return nil
}

func (r *gradeRepository) GetByStudent(ctx context.Context, studentID string) ([]*models.Grade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	grades, ok := r.store[studentID]
	if !ok {
		return []*models.Grade{}, nil
	}
	result := make([]*models.Grade, len(grades))
	copy(result, grades)
	return result, nil
}

func (r *gradeRepository) GetByLesson(ctx context.Context, lessonID string) ([]*models.Grade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	studentIDs, ok := r.lessonMap[lessonID]
	if !ok {
		return []*models.Grade{}, nil
	}
	var result []*models.Grade
	seen := make(map[string]bool)
	for _, sid := range studentIDs {
		if seen[sid] {
			continue
		}
		seen[sid] = true
		grades := r.store[sid]
		for _, g := range grades {
			if g.LessonID == lessonID {
				result = append(result, g)
			}
		}
	}
	return result, nil
}

func (r *gradeRepository) GetAverageForStudent(ctx context.Context, studentID string) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	grades, ok := r.store[studentID]
	if !ok || len(grades) == 0 {
		return 0, nil
	}
	sum := 0
	for _, g := range grades {
		sum += g.Value
	}
	return float64(sum) / float64(len(grades)), nil
}

func (r *gradeRepository) GetGradeBook(ctx context.Context, groupID string) (*models.GradeBook, error) {
	return nil, errors.NewInternalError("gradeRepository.GetGradeBook",
		liberrors.New("use GradingUsecase.GetGradeBook instead"))
}

func (r *gradeRepository) UpdateGrade(ctx context.Context, studentID, lessonID string, value int) error {
	const op = "gradeRepository.UpdateGrade"
	r.mu.Lock()
	defer r.mu.Unlock()

	grades, ok := r.store[studentID]
	if !ok {
		return errors.NewNotFoundError(op, "grades for student")
	}
	for i, g := range grades {
		if g.LessonID == lessonID {
			grades[i].Value = value
			grades[i].GradedAt = time.Now()
			return nil
		}
	}
	return errors.NewNotFoundError(op, "grade for this lesson")
}
