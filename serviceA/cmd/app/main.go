package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mins_EduCenter/serviceA/clients"
	"mins_EduCenter/serviceA/internal/delivery"
	"mins_EduCenter/serviceA/internal/observer"
	"mins_EduCenter/serviceA/internal/repository/impl"
	"mins_EduCenter/serviceA/internal/strategy"
	"mins_EduCenter/serviceA/internal/usecase"
)

func main() {
	// Подключаемся к Service B
	refClient, err := clients.NewReferenceClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect to Reference Service: %v", err)
	}
	defer refClient.Close()

	studentRepo := impl.NewStudentRepository()
	groupRepo := impl.NewGroupRepository()
	gradeRepo := impl.NewGradeRepository()
	lessonRepo := impl.NewLessonRepository()
	attendanceRepo := impl.NewAttendanceRepository()

	initialStrategy := &strategy.ArithmeticMean{}

	notifier := observer.NewNotifier()
	// logger := &observer.LoggerObserver{}
	console := &observer.ConsoleObserver{}
	estimatorUsecase := usecase.NewEstimatorUsecase()
	notifier.Subscribe(observer.EventStudentRegistered, console)

	gradingUsecase := usecase.NewGradingUsecase(gradeRepo, studentRepo, lessonRepo, groupRepo, initialStrategy)
	studentUsecase := usecase.NewStudentUsecase(studentRepo, groupRepo, gradeRepo, gradingUsecase, notifier)
	lessonUsecase := usecase.NewLessonUsecase(lessonRepo, attendanceRepo, groupRepo, studentRepo)
	groupUsecase := usecase.NewGroupUsecase(groupRepo, studentRepo, refClient)

	// CLI
	handler := delivery.NewHandler(studentUsecase, lessonUsecase, gradingUsecase, groupUsecase, estimatorUsecase)

	ctx := context.Background()
	log.Println("Service A (Core) started")

	go func() {
		handler.Run(ctx)
		os.Exit(0)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Service A...")
}
