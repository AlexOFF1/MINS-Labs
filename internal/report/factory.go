package report

import "fmt"

type ReportFactory struct{}

func (f *ReportFactory) CreateReport(reportType string) (ReportGenerator, error) {
	switch reportType {
	case "console", "text":
		return &ConsoleReport{}, nil
	case "html":
		return &HTMLReport{}, nil
	case "json":
		return &JSONReport{}, nil
	default:
		return nil, fmt.Errorf("unknown report type: %s", reportType)
	}
}
