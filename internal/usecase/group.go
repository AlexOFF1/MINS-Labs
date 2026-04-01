package usecase

import (
	"context"
	"mins_EduCenter/internal/models"
	"mins_EduCenter/internal/repository"
	"mins_EduCenter/pkg/errors"
	"time"
)

type GroupUsecase struct {
	groupRepo   repository.GroupRepository
	studentRepo repository.StudentRepository
}

func NewGroupUsecase(
	gr repository.GroupRepository,
	sr repository.StudentRepository,
) *GroupUsecase {
	return &GroupUsecase{
		groupRepo:   gr,
		studentRepo: sr,
	}
}

type CreateGroupDTO struct {
	Name        string
	CourseID    string
	StartDate   time.Time
	EndDate     time.Time
	MaxStudents int
}

func (u *GroupUsecase) CreateGroup(ctx context.Context, dto CreateGroupDTO) (*models.Group, error) {
	const op = "GroupUsecase.CreateGroup"

	if dto.Name == "" {
		return nil, errors.NewValidationError(op, "Name", "required")
	}
	if dto.MaxStudents <= 0 {
		dto.MaxStudents = 20
	}
	if dto.StartDate.IsZero() {
		dto.StartDate = time.Now()
	}
	if dto.EndDate.IsZero() {
		dto.EndDate = dto.StartDate.AddDate(0, 3, 0)
	}

	group := &models.Group{
		Name:        dto.Name,
		CourseID:    dto.CourseID,
		StartDate:   dto.StartDate,
		EndDate:     dto.EndDate,
		Status:      "active",
		MaxStudents: dto.MaxStudents,
		StudentIDs:  []string{},
	}

	if err := u.groupRepo.Create(ctx, group); err != nil {
		return nil, errors.NewInternalError(op, err)
	}
	return group, nil
}

func (u *GroupUsecase) GetGroup(ctx context.Context, groupID string) (*models.Group, error) {
	const op = "GroupUsecase.GetGroup"
	return u.groupRepo.GetByID(ctx, groupID)
}

func (u *GroupUsecase) GetAllGroups(ctx context.Context) ([]*models.Group, error) {
	const op = "GroupUsecase.GetAllGroups"
	return u.groupRepo.GetAll(ctx)
}

func (u *GroupUsecase) GetGroupStudents(ctx context.Context, groupID string) ([]*models.Student, error) {
	const op = "GroupUsecase.GetGroupStudents"

	if _, err := u.groupRepo.GetByID(ctx, groupID); err != nil {
		return nil, errors.NewValidationError(op, "GroupID", "group not found")
	}

	return u.studentRepo.GetByGroup(ctx, groupID)
}

func (u *GroupUsecase) UpdateGroupStatus(ctx context.Context, groupID, status string) error {
	const op = "GroupUsecase.UpdateGroupStatus"

	group, err := u.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return errors.NewValidationError(op, "GroupID", "group not found")
	}

	group.Status = status
	return u.groupRepo.Update(ctx, group)
}
