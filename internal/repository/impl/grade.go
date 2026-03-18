package memory

import (
	"context"
	"mins_EduCenter/internal/models"
	"mins_EduCenter/internal/repository"
	"sync"
	"time"
)

type gradeRepository struct {
	mu    sync.RWMutex
	store map[string][]*models.Grade // key: studentID
}

func NewGradeRepository() repository.GradeRepository {
	return &gradeRepository{
		store: make(map[string][]*models.Grade),
	}
}

func (r *gradeRepository) Set(ctx context.Context, grade *models.Grade) error {
	const op = "GradeRepository.Set"
	r.mu.Lock()
	defer r.mu.Unlock()

	grade.GradedAt = time.Now()
	r.store[grade.StudentID] = append(r.store[grade.StudentID], grade)
	return nil
}

func (r *gradeRepository) GetByStudent(ctx context.Context, studentID string) ([]*models.Grade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	grades, exists := r.store[studentID]
	if !exists {
		return []*models.Grade{}, nil
	}
	return grades, nil
}

func (r *gradeRepository) GetByLesson(ctx context.Context, lessonID string) ([]*models.Grade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Grade
	for _, grades := range r.store {
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

	grades, exists := r.store[studentID]
	if !exists || len(grades) == 0 {
		return 0, nil
	}
	var sum int
	for _, g := range grades {
		sum += g.Value
	}
	return float64(sum) / float64(len(grades)), nil
}
