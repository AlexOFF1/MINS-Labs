package memory

import (
	"context"
	"fmt"
	"mins_EduCenter/internal/models"
	"mins_EduCenter/internal/repository"
	"mins_EduCenter/pkg/errors"
	"sync"
	"time"
)

type gradeRepository struct {
	mu        sync.RWMutex
	store     map[string][]*models.Grade // studentID -> grades
	lessonMap map[string][]string        // lessonID -> []studentID (для быстрого поиска)
}

func NewGradeRepository() repository.GradeRepository {
	return &gradeRepository{
		store:     make(map[string][]*models.Grade),
		lessonMap: make(map[string][]string),
	}
}

// Set - выставление оценки
func (r *gradeRepository) Set(ctx context.Context, grade *models.Grade) error {
	const op = "GradeRepository.Set"
	r.mu.Lock()
	defer r.mu.Unlock()

	if grade.StudentID == "" || grade.LessonID == "" {
		return errors.NewValidationError(op, "grade", "studentID and lessonID are required")
	}

	grade.GradedAt = time.Now()

	// Добавляем оценку
	r.store[grade.StudentID] = append(r.store[grade.StudentID], grade)

	// Сохраняем в индекс по уроку
	r.lessonMap[grade.LessonID] = append(r.lessonMap[grade.LessonID], grade.StudentID)

	return nil
}

// GetByStudent - получить все оценки студента
func (r *gradeRepository) GetByStudent(ctx context.Context, studentID string) ([]*models.Grade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	grades, exists := r.store[studentID]
	if !exists {
		return []*models.Grade{}, nil // пустой слайс, а не nil и не ошибка
	}

	// Возвращаем копию, чтобы не меняли оригинал
	result := make([]*models.Grade, len(grades))
	copy(result, grades)
	return result, nil
}

// GetByLesson - получить все оценки за урок
func (r *gradeRepository) GetByLesson(ctx context.Context, lessonID string) ([]*models.Grade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	studentIDs, exists := r.lessonMap[lessonID]
	if !exists {
		return []*models.Grade{}, nil
	}

	var result []*models.Grade
	seen := make(map[string]bool) // чтобы не дублировать, если студент несколько оценок за урок

	for _, studentID := range studentIDs {
		if seen[studentID] {
			continue
		}
		seen[studentID] = true

		grades := r.store[studentID]
		for _, g := range grades {
			if g.LessonID == lessonID {
				result = append(result, g)
			}
		}
	}

	return result, nil
}

// GetAverageForStudent - средний балл студента
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

// GetGradeBook - получить ведомость группы
func (r *gradeRepository) GetGradeBook(ctx context.Context, groupID string) (*models.GradeBook, error) {
	const op = "GradeRepository.GetGradeBook"

	// Это сложный метод, требует доступа к другим репозиториям
	// В реальности он должен быть в usecase, но для полноты вернем заглушку
	return nil, errors.NewInternalError(op, fmt.Errorf("use GradingUsecase.GetGradeBook instead"))
}

// UpdateGrade - обновить оценку
func (r *gradeRepository) UpdateGrade(ctx context.Context, studentID, lessonID string, value int) error {
	const op = "GradeRepository.UpdateGrade"
	r.mu.Lock()
	defer r.mu.Unlock()

	grades, exists := r.store[studentID]
	if !exists {
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
