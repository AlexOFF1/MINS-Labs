package usecase

import (
	"context"
	"mins_EduCenter/pkg/errors"
	"mins_EduCenter/serviceA/clients"
	"mins_EduCenter/serviceA/internal/models"
	"mins_EduCenter/serviceA/internal/repository"
	"time"
)

type GroupUsecase struct {
	groupRepo   repository.GroupRepository
	studentRepo repository.StudentRepository
	refClient   *clients.ReferenceClient
}

func NewGroupUsecase(
	gr repository.GroupRepository,
	sr repository.StudentRepository,
	refCli *clients.ReferenceClient,
) *GroupUsecase {
	return &GroupUsecase{
		groupRepo:   gr,
		studentRepo: sr,
		refClient:   refCli,
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

	// Проверяем существование курса через Reference Service
	if dto.CourseID != "" {
		_, err := u.refClient.GetCourse(ctx, dto.CourseID)
		if err != nil {
			return nil, errors.NewValidationError(op, "CourseID", err.Error())
		}
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
