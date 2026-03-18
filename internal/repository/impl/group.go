package memory

import (
	"context"
	"mins_EduCenter/internal/models"
	"mins_EduCenter/internal/repository"
	"mins_EduCenter/pkg/errors"
	"sync"
	"time"
)

type groupRepository struct {
	mu    sync.RWMutex
	store map[string]*models.Group
}

func NewGroupRepository() repository.GroupRepository {
	return &groupRepository{
		store: make(map[string]*models.Group),
	}
}

func (r *groupRepository) Create(ctx context.Context, group *models.Group) error {
	const op = "GroupRepository.Create"
	r.mu.Lock()
	defer r.mu.Unlock()

	if group.ID == "" {
		group.ID = generateID("grp")
	}
	if _, exists := r.store[group.ID]; exists {
		return errors.NewDuplicateError(op, "Group", group.ID)
	}
	now := time.Now()
	group.CreatedAt = now
	group.UpdatedAt = now
	r.store[group.ID] = group
	return nil
}

func (r *groupRepository) GetByID(ctx context.Context, id string) (*models.Group, error) {
	const op = "GroupRepository.GetByID"
	r.mu.RLock()
	defer r.mu.RUnlock()

	group, exists := r.store[id]
	if !exists {
		return nil, errors.NewNotFoundError(op, "Group")
	}
	return group, nil
}

func (r *groupRepository) GetAll(ctx context.Context) ([]*models.Group, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	groups := make([]*models.Group, 0, len(r.store))
	for _, g := range r.store {
		groups = append(groups, g)
	}
	return groups, nil
}

func (r *groupRepository) Update(ctx context.Context, group *models.Group) error {
	const op = "GroupRepository.Update"
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[group.ID]; !exists {
		return errors.NewNotFoundError(op, "Group")
	}
	group.UpdatedAt = time.Now()
	r.store[group.ID] = group
	return nil
}

func (r *groupRepository) Delete(ctx context.Context, id string) error {
	const op = "GroupRepository.Delete"
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[id]; !exists {
		return errors.NewNotFoundError(op, "Group")
	}
	delete(r.store, id)
	return nil
}

func (r *groupRepository) AddStudent(ctx context.Context, groupID, studentID string) error {
	const op = "GroupRepository.AddStudent"
	r.mu.Lock()
	defer r.mu.Unlock()

	group, exists := r.store[groupID]
	if !exists {
		return errors.NewNotFoundError(op, "Group")
	}
	for _, id := range group.StudentIDs {
		if id == studentID {
			return errors.NewDuplicateError(op, "Student", "already in group")
		}
	}
	group.StudentIDs = append(group.StudentIDs, studentID)
	group.UpdatedAt = time.Now()
	return nil
}

func (r *groupRepository) RemoveStudent(ctx context.Context, groupID, studentID string) error {
	const op = "GroupRepository.RemoveStudent"
	r.mu.Lock()
	defer r.mu.Unlock()

	group, exists := r.store[groupID]
	if !exists {
		return errors.NewNotFoundError(op, "Group")
	}
	for i, id := range group.StudentIDs {
		if id == studentID {
			group.StudentIDs = append(group.StudentIDs[:i], group.StudentIDs[i+1:]...)
			group.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.NewNotFoundError(op, "Student in group")
}

func (r *groupRepository) GetByCourse(ctx context.Context, courseID string) ([]*models.Group, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Group
	for _, g := range r.store {
		if g.CourseID == courseID {
			result = append(result, g)
		}
	}
	return result, nil
}
