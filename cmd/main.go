package main

import (
	"context"
	"log"
	cli "mins_EduCenter/internal/delivery"
	memory "mins_EduCenter/internal/repository/impl"
	"mins_EduCenter/internal/usecase"
)

func main() {
	studentRepo := memory.NewStudentRepository()
	groupRepo := memory.NewGroupRepository()
	gradeRepo := memory.NewGradeRepository()

	studentUsecase := usecase.NewStudentUsecase(
		studentRepo,
		groupRepo,
		gradeRepo,
	)

	handler := cli.NewHandler(studentUsecase)
	ctx := context.Background()
	log.Println("Запуск учебного центра...")
	handler.Run(ctx)
}
