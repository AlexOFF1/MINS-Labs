package cli

import (
	"bufio"
	"context"
	"fmt"
	"mins_EduCenter/internal/usecase"
	"mins_EduCenter/pkg/errors"
	"os"
	"strings"
)

type Handler struct {
	studentUsecase *usecase.StudentUsecase
	reader         *bufio.Reader
}

func NewHandler(su *usecase.StudentUsecase) *Handler {
	return &Handler{
		studentUsecase: su,
		reader:         bufio.NewReader(os.Stdin),
	}
}

func (h *Handler) Run(ctx context.Context) {
	fmt.Println("===================================")
	fmt.Println("Учебный центр - Система управления")
	fmt.Println("===================================")
	h.printHelp()

	for {
		fmt.Print("\n➤ ")
		input, _ := h.reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		cmd := parts[0]
		args := parts[1:]

		switch cmd {
		case "register", "reg":
			h.handleRegister(ctx, args)
		case "enroll":
			h.handleEnroll(ctx, args)
		case "progress", "prog":
			h.handleProgress(ctx, args)
		case "help", "h":
			h.printHelp()
		case "exit", "quit":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Printf("Неизвестная команда: %s\n", cmd)
			h.printHelp()
		}
	}
}

func (h *Handler) handleRegister(ctx context.Context, args []string) {
	if len(args) < 3 {
		fmt.Println("Использование: register <имя> <фамилия> <email> [телефон]")
		return
	}

	dto := usecase.RegisterDTO{
		FirstName: args[0],
		LastName:  args[1],
		Email:     args[2],
	}
	if len(args) >= 4 {
		dto.Phone = args[3]
	}

	student, err := h.studentUsecase.Register(ctx, dto)
	if err != nil {
		h.handleError(err)
		return
	}

	fmt.Printf("   Студент успешно зарегистрирован!\n")
	fmt.Printf("   ID: %s\n", student.ID)
	fmt.Printf("   Имя: %s %s\n", student.FirstName, student.LastName)
	fmt.Printf("   Студенческий билет: %s\n", student.StudentCard)
}

func (h *Handler) handleEnroll(ctx context.Context, args []string) {
	if len(args) < 2 {
		fmt.Println("  Использование: enroll <student_id> <group_id>")
		return
	}

	err := h.studentUsecase.EnrollToGroup(ctx, args[0], args[1])
	if err != nil {
		h.handleError(err)
		return
	}

	fmt.Printf(" Студент %s зачислен в группу %s\n", args[0], args[1])
}

func (h *Handler) handleProgress(ctx context.Context, args []string) {
	if len(args) < 1 {
		fmt.Println(" Использование: progress <student_id>")
		return
	}

	progress, err := h.studentUsecase.GetProgress(ctx, args[0])
	if err != nil {
		h.handleError(err)
		return
	}

	fmt.Println("\n Отчет об успеваемости")
	fmt.Println("────────────────────────")
	fmt.Printf("Студент: %s %s\n", progress.Student.FirstName, progress.Student.LastName)
	fmt.Printf("Email: %s\n", progress.Student.Email)
	fmt.Printf("Группа: %s\n", progress.Student.GroupID)
	fmt.Printf("Всего оценок: %d\n", progress.TotalGrades)
	fmt.Printf("Средний балл: %.2f\n", progress.AverageGrade)

	if len(progress.Grades) > 0 {
		fmt.Println("\nОценки:")
		for i, g := range progress.Grades {
			if i >= 5 {
				fmt.Printf("   ... и еще %d\n", len(progress.Grades)-5)
				break
			}
			fmt.Printf("   %d: %d\n", i+1, g.Value)
		}
	}
}

func (h *Handler) handleError(err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		switch appErr.Code {
		case "NOT_FOUND":
			fmt.Printf(" Не найдено: %v\n", appErr)
		case "VALIDATION_ERROR":
			fmt.Printf(" Ошибка валидации: %v\n", appErr)
		case "DUPLICATE_ENTRY":
			fmt.Printf(" Дубликат: %v\n", appErr)
		default:
			fmt.Printf(" Ошибка: %v\n", appErr)
		}
	} else {
		fmt.Printf(" Неизвестная ошибка: %v\n", err)
	}
}

func (h *Handler) printHelp() {
	fmt.Println("\n Доступные команды:")
	fmt.Println("  register, reg <имя> <фамилия> <email> [телефон] - регистрация студента")
	fmt.Println("  enroll <student_id> <group_id> - зачисление в группу")
	fmt.Println("  progress, prog <student_id> - успеваемость студента")
	fmt.Println("  help, h - показать справку")
	fmt.Println("  exit, quit - выход")
}
