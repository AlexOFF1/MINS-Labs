package report

import (
	"encoding/json"
	"fmt"
	"mins_EduCenter/serviceA/internal/models"
	"time"
)

type ReportGenerator interface {
	Generate(student *models.Student, grades []*models.Grade) string
	Name() string
}

// Консольный отчёт (уже есть, просто выделяем)
type ConsoleReport struct{}

func (r *ConsoleReport) Name() string { return "console" }

func (r *ConsoleReport) Generate(student *models.Student, grades []*models.Grade) string {
	avg := 0.0
	if len(grades) > 0 {
		sum := 0
		for _, g := range grades {
			sum += g.Value
		}
		avg = float64(sum) / float64(len(grades))
	}

	report := "\n=== ТАБЕЛЬ УСПЕВАЕМОСТИ ===\n"
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
	return report
}

// HTML отчёт
type HTMLReport struct{}

func (r *HTMLReport) Name() string { return "html" }

func (r *HTMLReport) Generate(student *models.Student, grades []*models.Grade) string {
	// Строим HTML-таблицу
	rows := ""
	for _, g := range grades {
		rows += fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%s</td></tr>", g.Type, g.Value, g.Comment)
	}
	return fmt.Sprintf(`<html><body>
<h1>Табель студента %s %s</h1>
<p>Группа: %s</p>
<p>Студ. билет: %s</p>
<table border="1">
<tr><th>Тип</th><th>Оценка</th><th>Комментарий</th></tr>
%s
</table>
</body></html>`, student.FirstName, student.LastName, student.GroupID, student.StudentCard, rows)
}

// JSON отчёт
type JSONReport struct{}

func (r *JSONReport) Name() string { return "json" }

func (r *JSONReport) Generate(student *models.Student, grades []*models.Grade) string {
	data := struct {
		Student   *models.Student `json:"student"`
		Grades    []*models.Grade `json:"grades"`
		Generated string          `json:"generated"`
	}{
		Student:   student,
		Grades:    grades,
		Generated: time.Now().Format(time.RFC3339),
	}
	b, _ := json.MarshalIndent(data, "", "  ")
	return string(b)
}
