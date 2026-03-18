package memory

import (
	"context"
	"mins_EduCenter/internal/models"
	"mins_EduCenter/internal/repository"
	"mins_EduCenter/pkg/errors"
	"sync"
	"time"
)

type studentRepository struct {
	mu    sync.RWMutex
	store map[string]*models.Student
}

func NewStudentRepository() repository.StudentRepository {
	return &studentRepository{
		store: make(map[string]*models.Student),
	}
}

func (r *studentRepository) Create(ctx context.Context, student *models.Student) error {
	const op = "StudentRepository.Create"
	r.mu.Lock()
	defer r.mu.Unlock()

	if student.ID == "" {
		student.ID = generateID("stu")
	}
	if _, exists := r.store[student.ID]; exists {
		return errors.NewDuplicateError(op, "Student", student.ID)
	}
	now := time.Now()
	student.CreatedAt = now
	student.UpdatedAt = now
	r.store[student.ID] = student
	return nil
}

func (r *studentRepository) GetByID(ctx context.Context, id string) (*models.Student, error) {
	const op = "StudentRepository.GetByID"
	r.mu.RLock()
	defer r.mu.RUnlock()

	student, exists := r.store[id]
	if !exists {
		return nil, errors.NewNotFoundError(op, "Student")
	}
	return student, nil
}

func (r *studentRepository) GetAll(ctx context.Context) ([]*models.Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	students := make([]*models.Student, 0, len(r.store))
	for _, s := range r.store {
		students = append(students, s)
	}
	return students, nil
}

func (r *studentRepository) Update(ctx context.Context, student *models.Student) error {
	const op = "StudentRepository.Update"
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[student.ID]; !exists {
		return errors.NewNotFoundError(op, "Student")
	}
	student.UpdatedAt = time.Now()
	r.store[student.ID] = student
	return nil
}

func (r *studentRepository) Delete(ctx context.Context, id string) error {
	const op = "StudentRepository.Delete"
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[id]; !exists {
		return errors.NewNotFoundError(op, "Student")
	}
	delete(r.store, id)
	return nil
}

func (r *studentRepository) GetByGroup(ctx context.Context, groupID string) ([]*models.Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Student
	for _, s := range r.store {
		if s.GroupID == groupID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (r *studentRepository) GetActive(ctx context.Context) ([]*models.Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Student
	for _, s := range r.store {
		if s.IsActive {
			result = append(result, s)
		}
	}
	return result, nil
}

func (r *studentRepository) Search(ctx context.Context, query string) ([]*models.Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Student
	for _, s := range r.store {
		if contains(s.FirstName, query) || contains(s.LastName, query) || contains(s.Email, query) {
			result = append(result, s)
		}
	}
	return result, nil
}

func generateID(prefix string) string {
	return prefix + "-" + time.Now().Format("20060102150405")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
