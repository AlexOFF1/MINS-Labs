package main

import (
	"context"
	"log"
	"mins_EduCenter/internal/delivery"
	"mins_EduCenter/internal/repository/impl"
	"mins_EduCenter/internal/usecase"
)

func main() {
	studentRepo := impl.NewStudentRepository()
	groupRepo := impl.NewGroupRepository()
	gradeRepo := impl.NewGradeRepository()
	lessonRepo := impl.NewLessonRepository()
	attendanceRepo := impl.NewAttendanceRepository()

	studentUsecase := usecase.NewStudentUsecase(studentRepo, groupRepo, gradeRepo)
	lessonUsecase := usecase.NewLessonUsecase(lessonRepo, attendanceRepo, groupRepo, studentRepo)
	gradingUsecase := usecase.NewGradingUsecase(gradeRepo, studentRepo, lessonRepo, groupRepo)
	groupUsecase := usecase.NewGroupUsecase(groupRepo, studentRepo)

	handler := delivery.NewHandler(studentUsecase, lessonUsecase, gradingUsecase, groupUsecase)
	ctx := context.Background()
	log.Println(" Запуск учебного центра...")
	handler.Run(ctx)
}
