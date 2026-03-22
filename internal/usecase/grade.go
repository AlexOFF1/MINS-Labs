package usecase

import (
	"context"
	"fmt"
	"mins_EduCenter/internal/models"
	"mins_EduCenter/internal/repository"
	"mins_EduCenter/pkg/errors"
	"time"
)

type GradingUsecase struct {
	gradeRepo   repository.GradeRepository
	studentRepo repository.StudentRepository
	lessonRepo  repository.LessonRepository
	groupRepo   repository.GroupRepository
}

func NewGradingUsecase(
	gr repository.GradeRepository,
	sr repository.StudentRepository,
	lr repository.LessonRepository,
	gpr repository.GroupRepository,
) *GradingUsecase {
	return &GradingUsecase{
		gradeRepo:   gr,
		studentRepo: sr,
		lessonRepo:  lr,
		groupRepo:   gpr,
	}
}

type SetGradeDTO struct {
	StudentID string
	LessonID  string
	Value     int
	Comment   string
	GradedBy  string
	Type      string
}

func (u *GradingUsecase) SetGrade(ctx context.Context, dto SetGradeDTO) error {
	const op = "GradingUsecase.SetGrade"

	if dto.Value < 1 || dto.Value > 5 {
		return errors.NewValidationError(op, "Value", "must be between 1 and 5")
	}

	student, err := u.studentRepo.GetByID(ctx, dto.StudentID)
	if err != nil {
		return errors.NewValidationError(op, "StudentID", "student not found")
	}

	lesson, err := u.lessonRepo.GetByID(ctx, dto.LessonID)
	if err != nil {
		return errors.NewValidationError(op, "LessonID", "lesson not found")
	}

	if student.GroupID != lesson.GroupID {
		return errors.NewValidationError(op, "StudentID", "student not in this lesson's group")
	}

	grade := &models.Grade{
		StudentID: dto.StudentID,
		LessonID:  dto.LessonID,
		Value:     dto.Value,
		Comment:   dto.Comment,
		GradedAt:  time.Now(),
		GradedBy:  dto.GradedBy,
		Type:      dto.Type,
	}

	if err := u.gradeRepo.Set(ctx, grade); err != nil {
		return errors.NewInternalError(op, err)
	}

	return nil
}

func (u *GradingUsecase) GetStudentGrades(ctx context.Context, studentID string) ([]*models.Grade, error) {
	const op = "GradingUsecase.GetStudentGrades"

	grades, err := u.gradeRepo.GetByStudent(ctx, studentID)
	if err != nil {
		return nil, errors.NewInternalError(op, err)
	}
	return grades, nil
}

func (u *GradingUsecase) GetGradeBook(ctx context.Context, groupID string) (*models.GradeBook, error) {
	const op = "GradingUsecase.GetGradeBook"

	gradeBook, err := u.gradeRepo.GetGradeBook(ctx, groupID)
	if err != nil {
		return nil, errors.NewInternalError(op, err)
	}
	return gradeBook, nil
}

func (u *GradingUsecase) GetGroupAverage(ctx context.Context, groupID string) (float64, error) {
	const op = "GradingUsecase.GetGroupAverage"

	students, err := u.studentRepo.GetByGroup(ctx, groupID)
	if err != nil {
		return 0, errors.NewInternalError(op, err)
	}

	if len(students) == 0 {
		return 0, nil
	}

	var total float64
	var count int
	for _, student := range students {
		avg, err := u.gradeRepo.GetAverageForStudent(ctx, student.ID)
		if err != nil {
			continue
		}
		if avg > 0 {
			total += avg
			count++
		}
	}

	if count == 0 {
		return 0, nil
	}
	return total / float64(count), nil
}

func (u *GradingUsecase) GenerateReportCard(ctx context.Context, studentID string) (string, error) {
	const op = "GradingUsecase.GenerateReportCard"

	student, err := u.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return "", errors.NewValidationError(op, "StudentID", "student not found")
	}

	grades, err := u.gradeRepo.GetByStudent(ctx, studentID)
	if err != nil {
		return "", errors.NewInternalError(op, err)
	}

	avg, _ := u.gradeRepo.GetAverageForStudent(ctx, studentID)

	report := fmt.Sprintf("\n=== ТАБЕЛЬ УСПЕВАЕМОСТИ ===\n")
	report += fmt.Sprintf("Студент: %s %s\n", student.FirstName, student.LastName)
	report += fmt.Sprintf("Группа: %s\n", student.GroupID)
	report += fmt.Sprintf("Студенческий билет: %s\n", student.StudentCard)
	report += "----------------------------\n"

	if len(grades) == 0 {
		report += "Оценок пока нет\n"
	} else {
		report += "Оценки:\n"
		for i, g := range grades {
			report += fmt.Sprintf("  %d. %d (%s) - %s\n", i+1, g.Value, g.Type, g.Comment)
		}
		report += fmt.Sprintf("\nСредний балл: %.2f\n", avg)
	}
	report += "============================\n"

	return report, nil
}
