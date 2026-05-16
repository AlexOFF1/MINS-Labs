package usecase

import (
	"context"
	"mins_EduCenter/pkg/errors"
	"mins_EduCenter/pkg/validation"
	"mins_EduCenter/serviceA/internal/models"
	"mins_EduCenter/serviceA/internal/observer"
	"mins_EduCenter/serviceA/internal/repository"
	"time"
)

type StudentUsecase struct {
	studentRepo repository.StudentRepository
	groupRepo   repository.GroupRepository
	gradeRepo   repository.GradeRepository
	avgCalc     repository.AverageCalculator
	notifier    *observer.Notifier
}

func NewStudentUsecase(
	sr repository.StudentRepository,
	gr repository.GroupRepository,
	gdr repository.GradeRepository,
	avgCalc repository.AverageCalculator,
	n *observer.Notifier,
) *StudentUsecase {
	return &StudentUsecase{
		studentRepo: sr,
		groupRepo:   gr,
		gradeRepo:   gdr,
		avgCalc:     avgCalc,
		notifier:    n,
	}
}

type RegisterDTO struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

func (u *StudentUsecase) Register(ctx context.Context, dto RegisterDTO) (*models.Student, error) {
	const op = "StudentUsecase.Register"

	studentData := &validation.StudentRegisterData{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Phone:     dto.Phone,
	}

	if err := validation.Validate(studentData); err != nil {
		if valErr, ok := err.(*validation.ValidationError); ok {
			return nil, errors.NewValidationError(op, valErr.Field, valErr.Message)
		}
		return nil, errors.NewValidationError(op, "unknown", err.Error())
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

	if err := u.studentRepo.Create(ctx, student); err != nil {
		return nil, errors.NewInternalError(op, err)
	}

	u.notifier.Notify(observer.EventStudentRegistered, map[string]interface{}{
		"student_id": student.ID,
		"email":      student.Email,
	})

	return student, nil
}

func (u *StudentUsecase) EnrollToGroup(ctx context.Context, studentID, groupID string) error {
	const op = "StudentUsecase.EnrollToGroup"

	if studentID == "" {
		return errors.NewValidationError(op, "studentID", "required")
	}
	if groupID == "" {
		return errors.NewValidationError(op, "groupID", "required")
	}

	student, err := u.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return errors.NewValidationError(op, "studentID", "student not found")
	}

	if student.GroupID != "" {
		return errors.NewValidationError(op, "studentID", "student already enrolled in a group")
	}

	group, err := u.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return errors.NewValidationError(op, "groupID", "group not found")
	}

	if len(group.StudentIDs) >= group.MaxStudents {
		return errors.NewValidationError(op, "groupID", "group is full")
	}

	if err := u.groupRepo.AddStudent(ctx, groupID, studentID); err != nil {
		return errors.NewInternalError(op, err)
	}

	student.GroupID = groupID
	if err := u.studentRepo.Update(ctx, student); err != nil {
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

func (u *StudentUsecase) GetProgress(ctx context.Context, studentID string) (*ProgressReport, error) {
	const op = "StudentUsecase.GetProgress"

	student, err := u.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return nil, errors.NewValidationError(op, "studentID", "student not found")
	}

	grades, err := u.gradeRepo.GetByStudent(ctx, studentID)
	if err != nil {
		return nil, errors.NewInternalError(op, err)
	}

	avg, _ := u.avgCalc.GetAverageForStudent(ctx, studentID)

	return &ProgressReport{
		Student:      student,
		Grades:       grades,
		AverageGrade: avg,
		TotalGrades:  len(grades),
	}, nil
}

func generateStudentCard() string {
	return "STU-" + time.Now().Format("2006-0001")
}
