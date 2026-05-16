package delivery

import (
	"bufio"
	"context"
	"fmt"
	"mins_EduCenter/pkg/errors"
	"mins_EduCenter/serviceA/internal/models"
	"mins_EduCenter/serviceA/internal/strategy"
	"mins_EduCenter/serviceA/internal/usecase"
	"os"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	studentUsecase *usecase.StudentUsecase
	lessonUsecase  *usecase.LessonUsecase
	gradingUsecase *usecase.GradingUsecase
	groupUsecase   *usecase.GroupUsecase
	reader         *bufio.Reader
	strategies     map[string]strategy.AverageStrategy
	estimatorUC    *usecase.EstimatorUsecase
}

func NewHandler(
	su *usecase.StudentUsecase,
	lu *usecase.LessonUsecase,
	gu *usecase.GradingUsecase,
	grpU *usecase.GroupUsecase,
	estU *usecase.EstimatorUsecase,
) *Handler {
	strategies := map[string]strategy.AverageStrategy{
		"arithmetic": &strategy.ArithmeticMean{},
		"median":     &strategy.MedianStrategy{},
		"drop_worst": &strategy.DropWorstStrategy{},
	}

	return &Handler{
		studentUsecase: su,
		lessonUsecase:  lu,
		gradingUsecase: gu,
		groupUsecase:   grpU,
		reader:         bufio.NewReader(os.Stdin),
		strategies:     strategies,
		estimatorUC:    estU,
	}
}

func (h *Handler) handleRegister(ctx context.Context, args []string) {
	if len(args) < 3 {
		fmt.Println("❌ Использование: register <имя> <фамилия> <email> [телефон]")
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
	fmt.Printf("✅ Студент успешно зарегистрирован!\n")
	fmt.Printf("   ID: %s\n", student.ID)
	fmt.Printf("   Имя: %s %s\n", student.FirstName, student.LastName)
	fmt.Printf("   Студенческий билет: %s\n", student.StudentCard)
}

func (h *Handler) handleEnroll(ctx context.Context, args []string) {
	if len(args) < 2 {
		fmt.Println("❌ Использование: enroll <student_id> <group_id>")
		return
	}
	err := h.studentUsecase.EnrollToGroup(ctx, args[0], args[1])
	if err != nil {
		h.handleError(err)
		return
	}
	fmt.Printf("✅ Студент %s зачислен в группу %s\n", args[0], args[1])
}

func (h *Handler) handleProgress(ctx context.Context, args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Использование: progress <student_id>")
		return
	}
	progress, err := h.studentUsecase.GetProgress(ctx, args[0])
	if err != nil {
		h.handleError(err)
		return
	}
	fmt.Println("\n📊 Отчет об успеваемости")
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

func (h *Handler) handleCreateLesson(ctx context.Context, args []string) {
	if len(args) < 6 {
		fmt.Println("❌ Использование: lesson <group_id> <topic> <дата> <начало> <конец> <аудитория>")
		fmt.Println("   Пример: lesson grp123 'Go Basics' 2026-03-20 10:00 11:30 B-201")
		return
	}
	dateStr := args[2]
	startStr := args[3]
	endStr := args[4]
	startTime, err := time.Parse("2006-01-02 15:04", dateStr+" "+startStr)
	if err != nil {
		fmt.Println("❌ Неверный формат времени начала. Используйте ЧЧ:ММ")
		return
	}
	endTime, err := time.Parse("2006-01-02 15:04", dateStr+" "+endStr)
	if err != nil {
		fmt.Println("❌ Неверный формат времени конца. Используйте ЧЧ:ММ")
		return
	}
	dto := usecase.CreateLessonDTO{
		GroupID:     args[0],
		Topic:       args[1],
		Description: "",
		StartTime:   startTime,
		EndTime:     endTime,
		Room:        args[5],
		TeacherID:   "teacher-1",
	}
	lesson, err := h.lessonUsecase.CreateLesson(ctx, dto)
	if err != nil {
		h.handleError(err)
		return
	}
	fmt.Printf("✅ Занятие создано! ID: %s\n", lesson.ID)
}

func (h *Handler) handleSchedule(ctx context.Context, args []string) {
	if len(args) < 2 {
		fmt.Println("❌ Использование: schedule student <student_id>  или  schedule group <group_id>")
		return
	}
	var lessons []*models.Lesson
	var err error
	switch args[0] {
	case "student":
		lessons, err = h.lessonUsecase.GetStudentSchedule(ctx, args[1])
	case "group":
		lessons, err = h.lessonUsecase.GetGroupSchedule(ctx, args[1])
	default:
		fmt.Println("❌ Укажите student или group")
		return
	}
	if err != nil {
		h.handleError(err)
		return
	}
	fmt.Println("\n📅 РАСПИСАНИЕ:")
	fmt.Println("────────────────────")
	if len(lessons) == 0 {
		fmt.Println("Нет занятий")
		return
	}
	for _, l := range lessons {
		fmt.Printf("📌 %s\n", l.Topic)
		fmt.Printf("   Время: %s - %s\n", l.StartTime.Format("02.01.2006 15:04"), l.EndTime.Format("15:04"))
		fmt.Printf("   Аудитория: %s\n", l.Room)
		fmt.Printf("   ID: %s\n\n", l.ID)
	}
}

func (h *Handler) handleMarkAttendance(ctx context.Context, args []string) {
	if len(args) < 3 {
		fmt.Println("❌ Использование: mark <lesson_id> <student_id> <+|->")
		fmt.Println("   + - присутствовал, - - отсутствовал")
		return
	}
	present := args[2] == "+"
	dto := usecase.MarkAttendanceDTO{
		LessonID:  args[0],
		StudentID: args[1],
		Present:   present,
		MarkedBy:  "teacher",
	}
	err := h.lessonUsecase.MarkAttendance(ctx, dto)
	if err != nil {
		h.handleError(err)
		return
	}
	fmt.Printf("✅ Посещаемость отмечена\n")
}

func (h *Handler) handleAttendance(ctx context.Context, args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Использование: attendance <lesson_id>")
		return
	}
	attendance, err := h.lessonUsecase.GetLessonAttendance(ctx, args[0])
	if err != nil {
		h.handleError(err)
		return
	}
	fmt.Println("\n📋 ПОСЕЩАЕМОСТЬ:")
	fmt.Println("────────────────")
	if len(attendance) == 0 {
		fmt.Println("Нет данных о посещаемости")
		return
	}
	for _, a := range attendance {
		status := "❌"
		if a.Present {
			status = "✅"
		}
		fmt.Printf("%s %s\n", status, a.StudentID)
	}
}

func (h *Handler) handleSetGrade(ctx context.Context, args []string) {
	if len(args) < 4 {
		fmt.Println("❌ Использование: grade <student_id> <lesson_id> <оценка(1-5)> <тип> [комментарий]")
		fmt.Println("   Тип: exam, test, homework")
		return
	}
	value, err := strconv.Atoi(args[2])
	if err != nil || value < 1 || value > 5 {
		fmt.Println("❌ Оценка должна быть числом от 1 до 5")
		return
	}
	dto := usecase.SetGradeDTO{
		StudentID: args[0],
		LessonID:  args[1],
		Value:     value,
		Type:      args[3],
		GradedBy:  "teacher",
	}
	if len(args) > 4 {
		dto.Comment = strings.Join(args[4:], " ")
	}
	err = h.gradingUsecase.SetGrade(ctx, dto)
	if err != nil {
		h.handleError(err)
		return
	}
	fmt.Printf("✅ Оценка выставлена\n")
}

func (h *Handler) handleGetGrades(ctx context.Context, args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Использование: grades <student_id>")
		return
	}
	grades, err := h.gradingUsecase.GetStudentGrades(ctx, args[0])
	if err != nil {
		h.handleError(err)
		return
	}
	fmt.Println("\n📊 ОЦЕНКИ СТУДЕНТА:")
	fmt.Println("───────────────────")
	if len(grades) == 0 {
		fmt.Println("Оценок нет")
		return
	}
	for i, g := range grades {
		fmt.Printf("%d. %d (%s) - %s\n", i+1, g.Value, g.Type, g.Comment)
	}
}

func (h *Handler) handleGradeBook(ctx context.Context, args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Использование: gradebook <group_id>")
		return
	}
	gradeBook, err := h.gradingUsecase.GetGradeBook(ctx, args[0])
	if err != nil {
		h.handleError(err)
		return
	}
	fmt.Printf("\n📚 ВЕДОМОСТЬ ГРУППЫ %s\n", gradeBook.GroupID)
	fmt.Printf("Курс: %s\n", gradeBook.CourseName)
	fmt.Printf("Сгенерировано: %s\n", gradeBook.GeneratedAt.Format("02.01.2006 15:04"))
	fmt.Println("=================================")
	for _, student := range gradeBook.Students {
		fmt.Printf("\n👤 %s %s\n", student.FirstName, student.LastName)
		studentGrades := gradeBook.Grades[student.ID]
		if len(studentGrades) > 0 {
			for _, g := range studentGrades {
				fmt.Printf("   %d ", g.Value)
			}
			fmt.Println()
		} else {
			fmt.Println("   Нет оценок")
		}
	}
}

func (h *Handler) handleReportCard(ctx context.Context, args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Использование: reportcard <student_id>")
		return
	}
	report, err := h.gradingUsecase.GenerateReportCard(ctx, args[0])
	if err != nil {
		h.handleError(err)
		return
	}
	fmt.Print(report)
}

func (h *Handler) handleError(err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		switch appErr.Code {
		case "NOT_FOUND":
			fmt.Printf("❌ Не найдено: %v\n", appErr)
		case "VALIDATION_ERROR":
			fmt.Printf("❌ Ошибка валидации: %v\n", appErr)
		case "DUPLICATE_ENTRY":
			fmt.Printf("❌ Дубликат: %v\n", appErr)
		default:
			fmt.Printf("❌ Ошибка: %v\n", appErr)
		}
	} else {
		fmt.Printf("❌ Неизвестная ошибка: %v\n", err)
	}
}

func (h *Handler) handleCreateGroup(ctx context.Context, args []string) {
	if len(args) < 2 {
		fmt.Println("❌ Использование: create-group <название> <макс_студентов> [курс_id]")
		fmt.Println("   Пример: create-group 'Go Developers' 20 course-1")
		fmt.Println("   Пример: create-group 'Python Basics' 15")
		return
	}

	maxStudents, err := strconv.Atoi(args[1])
	if err != nil || maxStudents <= 0 {
		fmt.Println("❌ Максимальное количество студентов должно быть положительным числом")
		return
	}

	courseID := ""
	if len(args) >= 3 {
		courseID = args[2]
	}

	dto := usecase.CreateGroupDTO{
		Name:        args[0],
		CourseID:    courseID,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 3, 0),
		MaxStudents: maxStudents,
	}

	group, err := h.groupUsecase.CreateGroup(ctx, dto)
	if err != nil {
		h.handleError(err)
		return
	}

	fmt.Printf("✅ Группа успешно создана!\n")
	fmt.Printf("   ID: %s\n", group.ID)
	fmt.Printf("   Название: %s\n", group.Name)
	fmt.Printf("   Макс. студентов: %d\n", group.MaxStudents)
	fmt.Printf("   Статус: %s\n", group.Status)
}

func (h *Handler) handleListGroups(ctx context.Context, args []string) {
	groups, err := h.groupUsecase.GetAllGroups(ctx)
	if err != nil {
		h.handleError(err)
		return
	}

	fmt.Println("\n📚 СПИСОК ГРУПП:")
	fmt.Println("─────────────────────────────────")
	if len(groups) == 0 {
		fmt.Println("Нет созданных групп")
		return
	}

	for _, g := range groups {
		fmt.Printf("📌 %s\n", g.Name)
		fmt.Printf("   ID: %s\n", g.ID)
		fmt.Printf("   Студентов: %d / %d\n", len(g.StudentIDs), g.MaxStudents)
		fmt.Printf("   Статус: %s\n", g.Status)
		fmt.Printf("   Период: %s - %s\n",
			g.StartDate.Format("02.01.2006"),
			g.EndDate.Format("02.01.2006"))
		fmt.Println("─────────────────────────────────")
	}
}

func (h *Handler) handleGroupStudents(ctx context.Context, args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Использование: group-students <group_id>")
		return
	}

	students, err := h.groupUsecase.GetGroupStudents(ctx, args[0])
	if err != nil {
		h.handleError(err)
		return
	}

	group, _ := h.groupUsecase.GetGroup(ctx, args[0])

	fmt.Printf("\n👥 СТУДЕНТЫ ГРУППЫ %s\n", group.Name)
	fmt.Println("─────────────────────────────────")
	if len(students) == 0 {
		fmt.Println("В группе пока нет студентов")
		return
	}

	for i, s := range students {
		fmt.Printf("%d. %s %s\n", i+1, s.FirstName, s.LastName)
		fmt.Printf("   ID: %s\n", s.ID)
		fmt.Printf("   Email: %s\n", s.Email)
		fmt.Printf("   Студ. билет: %s\n", s.StudentCard)
		fmt.Println("─────────────────────────────────")
	}
}

func (h *Handler) handleGroupStatus(ctx context.Context, args []string) {
	if len(args) < 2 {
		fmt.Println("❌ Использование: group-status <group_id> <active|finished|cancelled>")
		return
	}

	status := args[1]
	if status != "active" && status != "finished" && status != "cancelled" {
		fmt.Println("❌ Статус должен быть: active, finished или cancelled")
		return
	}

	err := h.groupUsecase.UpdateGroupStatus(ctx, args[0], status)
	if err != nil {
		h.handleError(err)
		return
	}

	fmt.Printf("✅ Статус группы изменен на: %s\n", status)
}

func (h *Handler) handleReport(ctx context.Context, args []string) {
	if len(args) < 2 {
		fmt.Println("❌ Использование: report <student_id> <console|html|json>")
		return
	}
	content, err := h.gradingUsecase.GenerateReport(ctx, args[0], args[1])
	if err != nil {
		h.handleError(err)
		return
	}
	fmt.Println(content)
}

func (h *Handler) handleListStrategies(ctx context.Context, args []string) {
	current := h.gradingUsecase.GetCurrentStrategyName()
	fmt.Println("\n📊 Доступные стратегии расчёта среднего балла:")
	for name := range h.strategies {
		marker := " "
		if name == current {
			marker = "✓"
		}
		fmt.Printf("  %s %s\n", marker, name)
	}
	fmt.Println("\nИспользуйте: strategy <название> для выбора стратегии")
}

func (h *Handler) handleStrategy(ctx context.Context, args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Использование: strategy <arithmetic|median|drop_worst>")
		h.handleListStrategies(ctx, nil)
		return
	}
	strategyName := args[0]
	newStrategy, exists := h.strategies[strategyName]
	if !exists {
		fmt.Printf("❌ Неизвестная стратегия: %s\n", strategyName)
		h.handleListStrategies(ctx, nil)
		return
	}
	h.gradingUsecase.SetStrategy(newStrategy)
	fmt.Printf("✅ Стратегия расчёта среднего балла изменена на: %s\n", strategyName)
}

func (h *Handler) handleEstimate(ctx context.Context, args []string) {
	if len(args) < 2 {
		fmt.Println("❌ Использование: estimate <название курса> <количество часов>")
		fmt.Println("   Пример: estimate \"Go Advanced\" 50")
		return
	}
	hours, err := strconv.Atoi(args[1])
	if err != nil || hours <= 0 {
		fmt.Println("❌ Часы должны быть положительным числом")
		return
	}

	price, err := h.estimatorUC.EstimateCoursePrice(ctx, args[0], hours)
	if err != nil {
		h.handleError(err)
		return
	}

	fmt.Printf("\n💰 ПРИМЕРНАЯ СТОИМОСТЬ КУРСА\n")
	fmt.Printf("   Курс: %s\n", args[0])
	fmt.Printf("   Часы: %d\n", hours)
	fmt.Printf("   Стоимость (с налогом): %.2f руб.\n", price)
}

func (h *Handler) Run(ctx context.Context) {
	fmt.Println("===================================")
	fmt.Println("📚 Учебный центр - Система управления (Service A)")
	fmt.Println("===================================")
	h.printHelp()

	for {
		fmt.Print("\n➤ ")
		input, _ := h.reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Генерируем traceID на каждую команду
		traceID := fmt.Sprintf("trace-%d", time.Now().UnixNano())
		ctxWithTrace := context.WithValue(ctx, "traceID", traceID)

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}
		cmd := parts[0]
		args := parts[1:]

		// В каждый обработчик передаём ctxWithTrace
		switch cmd {
		case "create-group", "cg":
			h.handleCreateGroup(ctxWithTrace, args)
		case "list-groups", "lg":
			h.handleListGroups(ctxWithTrace, args)
		case "group-students", "gs":
			h.handleGroupStudents(ctxWithTrace, args)
		case "group-status", "gst":
			h.handleGroupStatus(ctxWithTrace, args)

		case "register", "reg":
			h.handleRegister(ctxWithTrace, args)
		case "enroll":
			h.handleEnroll(ctxWithTrace, args)
		case "progress", "prog":
			h.handleProgress(ctxWithTrace, args)

		case "strategy", "strat":
			h.handleStrategy(ctxWithTrace, args)
		case "list-strategies", "lstr":
			h.handleListStrategies(ctxWithTrace, args)

		case "lesson", "ls":
			h.handleCreateLesson(ctxWithTrace, args)
		case "schedule", "sched":
			h.handleSchedule(ctxWithTrace, args)

		case "mark":
			h.handleMarkAttendance(ctxWithTrace, args)
		case "attendance", "attend":
			h.handleAttendance(ctxWithTrace, args)

		case "grade", "g":
			h.handleSetGrade(ctxWithTrace, args)
		case "grades":
			h.handleGetGrades(ctxWithTrace, args)
		case "gradebook", "gb":
			h.handleGradeBook(ctxWithTrace, args)
		case "reportcard", "rc":
			h.handleReportCard(ctxWithTrace, args)
		case "report", "r":
			h.handleReport(ctxWithTrace, args)

		case "estimate", "est":
			h.handleEstimate(ctxWithTrace, args)

		case "help", "h":
			h.printHelp()
		case "exit", "quit":
			fmt.Println("👋 До свидания!")
			return
		default:
			fmt.Printf("❌ Неизвестная команда: %s\n", cmd)
			h.printHelp()
		}
	}
}

func (h *Handler) printHelp() {
	fmt.Println("\n📖 Доступные команды:")
	fmt.Println("")
	fmt.Println("  ГРУППЫ:")
	fmt.Println("    create-group, cg <название> <макс_студентов> [курс_id] - создать группу")
	fmt.Println("    list-groups, lg - список всех групп")
	fmt.Println("    group-students, gs <group_id> - список студентов группы")
	fmt.Println("    group-status, gst <group_id> <active|finished|cancelled> - изменить статус группы")
	fmt.Println("")
	fmt.Println("  СТУДЕНТЫ:")
	fmt.Println("    register, reg <имя> <фамилия> <email> [телефон] - регистрация студента")
	fmt.Println("    enroll <student_id> <group_id> - зачисление в группу")
	fmt.Println("    progress, prog <student_id> - успеваемость студента")
	fmt.Println("")
	fmt.Println("  РАСПИСАНИЕ:")
	fmt.Println("    lesson <group_id> <topic> <дата> <начало> <конец> <аудитория> - создать занятие")
	fmt.Println("    schedule student <student_id> - расписание студента")
	fmt.Println("    schedule group <group_id> - расписание группы")
	fmt.Println("")
	fmt.Println("  ПОСЕЩАЕМОСТЬ:")
	fmt.Println("    mark <lesson_id> <student_id> <+|-> - отметить одного")
	fmt.Println("    attendance <lesson_id> - показать посещаемость урока")
	fmt.Println("")
	fmt.Println("  ОЦЕНКИ:")
	fmt.Println("    grade <student_id> <lesson_id> <оценка> <тип> [коммент] - выставить оценку")
	fmt.Println("    grades <student_id> - оценки студента")
	fmt.Println("    gradebook <group_id> - ведомость группы")
	fmt.Println("    reportcard <student_id> - табель успеваемости")
	fmt.Println("    report <student_id> <console|html|json> - табель успеваемости")
	fmt.Println("")
	fmt.Println("  ОБЩЕЕ:")
	fmt.Println("    estimate <название курса> <количество часов> - расчет примерной стоимости курса")
	fmt.Println("    help, h - показать справку")
	fmt.Println("    exit, quit - выход")
}
