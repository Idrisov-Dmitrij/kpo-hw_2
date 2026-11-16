package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"pos-app/internal/application"
)

type CSVExporter struct{}

func NewCSVExporter() *CSVExporter {
	return &CSVExporter{}
}

func (e *CSVExporter) GetFormat() string { return "csv" }

func (e *CSVExporter) Export(report application.GradeReport, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("не удалось создать файл %s: %w", path, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(e.buildRecords(report)); err != nil {
		return fmt.Errorf("не удалось записать CSV: %w", err)
	}
	return nil
}

func (e *CSVExporter) buildRecords(report application.GradeReport) [][]string {
	header := []string{"student_id", "full_name", "group"}
	for _, criterion := range report.Criteria {
		header = append(header, fmt.Sprintf("%s (max %s, вес %s)",
			criterion.Name, formatNumber(criterion.MaxScore), formatNumber(criterion.Weight)))
	}
	header = append(header, "graded", "final_grade")

	records := [][]string{header}
	for _, row := range report.Rows {
		record := []string{row.StudentID, row.FullName, row.Group}
		for _, cell := range row.Scores {
			if cell.Graded {
				record = append(record, formatNumber(cell.Value))
			} else {
				record = append(record, "")
			}
		}
		record = append(record,
			fmt.Sprintf("%d/%d", row.GradedCount, row.TotalCriteria),
			strconv.FormatFloat(row.FinalGrade, 'f', 2, 64),
		)
		records = append(records, record)
	}

	return records
}

func formatNumber(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
