package main

import (
	"context"
	"log"
	"mins_EduCenter/internal/delivery"
	"mins_EduCenter/internal/observer"
	"mins_EduCenter/internal/repository/impl"
	"mins_EduCenter/internal/strategy"
	"mins_EduCenter/internal/usecase"
)

func main() {
	studentRepo := impl.NewStudentRepository()
	groupRepo := impl.NewGroupRepository()
	gradeRepo := impl.NewGradeRepository()
	lessonRepo := impl.NewLessonRepository()
	attendanceRepo := impl.NewAttendanceRepository()

	notifier := observer.NewNotifier()
	// logger := &observer.LoggerObserver{}
	console := &observer.ConsoleObserver{}

	notifier.Subscribe(observer.EventStudentRegistered, console)
	// notifier.Subscribe(observer.EventStudentRegistered, logger)
	// notifier.Subscribe(observer.EventStudentEnrolled, logger)
	// notifier.Subscribe(observer.EventGradeAdded, logger)

	initialStrategy := &strategy.ArithmeticMean{}

	gradingUsecase := usecase.NewGradingUsecase(gradeRepo, studentRepo, lessonRepo, groupRepo, initialStrategy)
	studentUsecase := usecase.NewStudentUsecase(studentRepo, groupRepo, gradeRepo, gradingUsecase, notifier)
	lessonUsecase := usecase.NewLessonUsecase(lessonRepo, attendanceRepo, groupRepo, studentRepo)
	groupUsecase := usecase.NewGroupUsecase(groupRepo, studentRepo)
	estimatorUsecase := usecase.NewEstimatorUsecase()

	handler := delivery.NewHandler(studentUsecase, lessonUsecase, gradingUsecase, groupUsecase, estimatorUsecase)
	ctx := context.Background()
	log.Println(" Запуск учебного центра...")
	handler.Run(ctx)
}
