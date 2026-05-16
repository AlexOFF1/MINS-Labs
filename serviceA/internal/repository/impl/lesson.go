package impl

import (
	"context"
	liberrors "errors"
	"mins_EduCenter/pkg/errors"
	"mins_EduCenter/serviceA/internal/models"
	"mins_EduCenter/serviceA/internal/repository"
	"sort"
	"sync"
	"time"
)

type lessonRepository struct {
	mu           sync.RWMutex
	store        map[string]*models.Lesson
	groupLessons map[string][]string
}

func NewLessonRepository() repository.LessonRepository {
	return &lessonRepository{
		store:        make(map[string]*models.Lesson),
		groupLessons: make(map[string][]string),
	}
}

func (r *lessonRepository) Create(ctx context.Context, lesson *models.Lesson) error {
	const op = "lessonRepository.Create"
	r.mu.Lock()
	defer r.mu.Unlock()

	if lesson.ID == "" {
		lesson.ID = generateID("lsn")
	}
	if _, exists := r.store[lesson.ID]; exists {
		return errors.NewDuplicateError(op, "Lesson", lesson.ID)
	}
	now := time.Now()
	lesson.CreatedAt = now
	lesson.UpdatedAt = now
	if lesson.Status == "" {
		lesson.Status = "scheduled"
	}
	r.store[lesson.ID] = lesson
	r.groupLessons[lesson.GroupID] = append(r.groupLessons[lesson.GroupID], lesson.ID)
	return nil
}

func (r *lessonRepository) GetByID(ctx context.Context, id string) (*models.Lesson, error) {
	const op = "lessonRepository.GetByID"
	r.mu.RLock()
	defer r.mu.RUnlock()

	lesson, exists := r.store[id]
	if !exists {
		return nil, errors.NewNotFoundError(op, "Lesson")
	}
	return lesson, nil
}

func (r *lessonRepository) GetByGroup(ctx context.Context, groupID string) ([]*models.Lesson, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lessonIDs, exists := r.groupLessons[groupID]
	if !exists {
		return []*models.Lesson{}, nil
	}
	result := make([]*models.Lesson, 0, len(lessonIDs))
	for _, id := range lessonIDs {
		if lesson, ok := r.store[id]; ok {
			result = append(result, lesson)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartTime.Before(result[j].StartTime)
	})
	return result, nil
}

func (r *lessonRepository) GetByDateRange(ctx context.Context, from, to time.Time) ([]*models.Lesson, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Lesson
	for _, lesson := range r.store {
		if (lesson.StartTime.After(from) || lesson.StartTime.Equal(from)) &&
			(lesson.EndTime.Before(to) || lesson.EndTime.Equal(to)) {
			result = append(result, lesson)
		}
	}
	return result, nil
}

func (r *lessonRepository) Update(ctx context.Context, lesson *models.Lesson) error {
	const op = "lessonRepository.Update"
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[lesson.ID]; !exists {
		return errors.NewNotFoundError(op, "Lesson")
	}
	lesson.UpdatedAt = time.Now()
	r.store[lesson.ID] = lesson
	return nil
}

func (r *lessonRepository) Delete(ctx context.Context, id string) error {
	const op = "lessonRepository.Delete"
	r.mu.Lock()
	defer r.mu.Unlock()

	lesson, exists := r.store[id]
	if !exists {
		return errors.NewNotFoundError(op, "Lesson")
	}
	if lessonIDs, ok := r.groupLessons[lesson.GroupID]; ok {
		for i, lid := range lessonIDs {
			if lid == id {
				r.groupLessons[lesson.GroupID] = append(lessonIDs[:i], lessonIDs[i+1:]...)
				break
			}
		}
	}
	delete(r.store, id)
	return nil
}

func (r *lessonRepository) GetScheduleForStudent(ctx context.Context, studentID string) ([]*models.Lesson, error) {
	return nil, errors.NewInternalError("lessonRepository.GetScheduleForStudent",
		liberrors.New("use LessonUsecase.GetStudentSchedule instead"))
}

func (r *lessonRepository) GetUpcomingLessons(ctx context.Context, groupID string, limit int) ([]*models.Lesson, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lessonIDs, exists := r.groupLessons[groupID]
	if !exists {
		return []*models.Lesson{}, nil
	}
	var upcoming []*models.Lesson
	now := time.Now()
	for _, id := range lessonIDs {
		lesson := r.store[id]
		if lesson.StartTime.After(now) {
			upcoming = append(upcoming, lesson)
		}
	}
	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].StartTime.Before(upcoming[j].StartTime)
	})
	if limit > 0 && len(upcoming) > limit {
		upcoming = upcoming[:limit]
	}
	return upcoming, nil
}
