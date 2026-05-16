package usecase

import (
	"context"
	"fmt"
	"mins_EduCenter/pkg/errors"
	"mins_EduCenter/pkg/validation"
	"mins_EduCenter/serviceA/internal/models"
	"mins_EduCenter/serviceA/internal/report"
	"mins_EduCenter/serviceA/internal/repository"
	"mins_EduCenter/serviceA/internal/strategy"
	"time"
)

type GradingUsecase struct {
	gradeRepo   repository.GradeRepository
	studentRepo repository.StudentRepository
	lessonRepo  repository.LessonRepository
	groupRepo   repository.GroupRepository
	strategy    strategy.AverageStrategy
}

func NewGradingUsecase(
	gr repository.GradeRepository,
	sr repository.StudentRepository,
	lr repository.LessonRepository,
	gpr repository.GroupRepository,
	strategy strategy.AverageStrategy,
) *GradingUsecase {
	return &GradingUsecase{
		gradeRepo:   gr,
		studentRepo: sr,
		lessonRepo:  lr,
		groupRepo:   gpr,
		strategy:    strategy,
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

	gradeData := &validation.GradeSetData{
		StudentID: dto.StudentID,
		LessonID:  dto.LessonID,
		Value:     dto.Value,
		Type:      dto.Type,
	}

	if err := validation.Validate(gradeData); err != nil {
		if valErr, ok := err.(*validation.ValidationError); ok {
			return errors.NewValidationError(op, valErr.Field, valErr.Message)
		}
		return errors.NewValidationError(op, "unknown", err.Error())
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
	_, err := u.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return nil, errors.NewValidationError(op, "GroupID", "group not found")
	}
	students, err := u.studentRepo.GetByGroup(ctx, groupID)
	if err != nil {
		return nil, errors.NewInternalError(op, err)
	}
	lessons, err := u.lessonRepo.GetByGroup(ctx, groupID)
	if err != nil {
		return nil, errors.NewInternalError(op, err)
	}
	gradesMap := make(map[string][]*models.Grade)
	for _, student := range students {
		grades, err := u.gradeRepo.GetByStudent(ctx, student.ID)
		if err != nil {
			continue
		}
		gradesMap[student.ID] = grades
	}
	gradeBook := &models.GradeBook{
		GroupID:     groupID,
		CourseName:  "Course Name",
		Lessons:     lessons,
		Students:    students,
		Grades:      gradesMap,
		GeneratedAt: time.Now(),
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
		return "", errors.NewValidationError(op, "studentID", "student not found")
	}
	grades, err := u.gradeRepo.GetByStudent(ctx, studentID)
	if err != nil {
		return "", errors.NewInternalError(op, err)
	}
	avg, _ := u.GetAverageForStudent(ctx, studentID)

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

func (u *GradingUsecase) GenerateReport(ctx context.Context, studentID string, reportType string) (string, error) {
	const op = "GradingUsecase.GenerateReport"

	student, err := u.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return "", errors.NewValidationError(op, "studentID", "student not found")
	}
	grades, err := u.gradeRepo.GetByStudent(ctx, studentID)
	if err != nil {
		return "", errors.NewInternalError(op, err)
	}

	factory := &report.ReportFactory{}
	generator, err := factory.CreateReport(reportType)
	if err != nil {
		return "", errors.NewValidationError(op, "reportType", err.Error())
	}

	return generator.Generate(student, grades), nil
}

func (u *GradingUsecase) SetStrategy(strategy strategy.AverageStrategy) {
	u.strategy = strategy
}

func (u *GradingUsecase) GetCurrentStrategyName() string {
	return u.strategy.Name()
}

func (u *GradingUsecase) GetAverageForStudent(ctx context.Context, studentID string) (float64, error) {
	grades, err := u.gradeRepo.GetByStudent(ctx, studentID)
	if err != nil {
		return 0, err
	}
	values := make([]int, len(grades))
	for i, g := range grades {
		values[i] = g.Value
	}
	fmt.Printf("[DEBUG] Использую стратегию: %s\n", u.strategy.Name())
	result := u.strategy.Calculate(values)
	fmt.Printf("[DEBUG] Оценки: %v, результат: %.2f\n", values, result)
	return result, nil
}
