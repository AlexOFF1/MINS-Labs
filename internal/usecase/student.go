package usecase

import (
	"context"
	"mins_EduCenter/internal/models"
	"mins_EduCenter/internal/repository"
	"mins_EduCenter/pkg/errors"
	"regexp"
	"time"
)

type StudentUsecase struct {
	studentRepo repository.StudentRepository
	groupRepo   repository.GroupRepository
	gradeRepo   repository.GradeRepository
}

func NewStudentUsecase(
	sr repository.StudentRepository,
	gr repository.GroupRepository,
	gdr repository.GradeRepository,
) *StudentUsecase {
	return &StudentUsecase{
		studentRepo: sr,
		groupRepo:   gr,
		gradeRepo:   gdr,
	}
}

type RegisterDTO struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

func (s *StudentUsecase) Register(ctx context.Context, dto RegisterDTO) (*models.Student, error) {
	const op = "StudentUsecase.Register"

	if err := s.validateRegisterData(dto); err != nil {
		return nil, err
	}

	student := &models.Student{
		Person: models.Person{
			FirstName: dto.FirstName,
			LastName:  dto.LastName,
			Email:     dto.Email,
			Phone:     dto.Phone,
		},
		EnrolledAt:  time.Now(),
		IsActive:    true,
		StudentCard: generateStudentCard(),
	}

	if err := s.studentRepo.Create(ctx, student); err != nil {
		return nil, errors.NewInternalError(op, err)
	}

	return student, nil
}

func (s *StudentUsecase) EnrollToGroup(ctx context.Context, studentID, groupID string) error {
	const op = "StudentUsecase.EnrollToGroup"

	// Check student exists
	student, err := s.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return errors.NewValidationError(op, "studentID", "student not found")
	}

	// Check group exists
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return errors.NewValidationError(op, "groupID", "group not found")
	}

	// Check group capacity
	if len(group.StudentIDs) >= group.MaxStudents {
		return errors.NewValidationError(op, "groupID", "group is full")
	}

	// Check if already enrolled
	for _, id := range group.StudentIDs {
		if id == studentID {
			return errors.NewDuplicateError(op, "Student", "already in group")
		}
	}

	// Add to group
	if err := s.groupRepo.AddStudent(ctx, groupID, studentID); err != nil {
		return errors.NewInternalError(op, err)
	}

	// Update student
	student.GroupID = groupID
	if err := s.studentRepo.Update(ctx, student); err != nil {
		return errors.NewInternalError(op, err)
	}

	return nil
}

type ProgressReport struct {
	Student      *models.Student
	Grades       []*models.Grade
	AverageGrade float64
	TotalGrades  int
}

func (s *StudentUsecase) GetProgress(ctx context.Context, studentID string) (*ProgressReport, error) {
	const op = "StudentUsecase.GetProgress"

	student, err := s.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return nil, errors.NewValidationError(op, "studentID", "student not found")
	}

	grades, err := s.gradeRepo.GetByStudent(ctx, studentID)
	if err != nil {
		return nil, errors.NewInternalError(op, err)
	}

	avg, err := s.gradeRepo.GetAverageForStudent(ctx, studentID)
	if err != nil {
		return nil, errors.NewInternalError(op, err)
	}

	return &ProgressReport{
		Student:      student,
		Grades:       grades,
		AverageGrade: avg,
		TotalGrades:  len(grades),
	}, nil
}

func (s *StudentUsecase) validateRegisterData(dto RegisterDTO) error {
	const op = "StudentUsecase.validateRegisterData"

	if dto.FirstName == "" {
		return errors.NewValidationError(op, "FirstName", "required")
	}
	if dto.LastName == "" {
		return errors.NewValidationError(op, "LastName", "required")
	}
	if dto.Email == "" {
		return errors.NewValidationError(op, "Email", "required")
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(dto.Email) {
		return errors.NewValidationError(op, "Email", "invalid format")
	}
	return nil
}

func generateStudentCard() string {
	return "STU-" + time.Now().Format("2006-0001")
}
