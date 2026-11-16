package export

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"pos-app/internal/application"
)

type jsonCriterion struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	MaxScore float64 `json:"max_score"`
	Weight   float64 `json:"weight"`
}

type jsonStudent struct {
	ID         string              `json:"id"`
	FullName   string              `json:"full_name"`
	Group      string              `json:"group,omitempty"`
	Scores     map[string]*float64 `json:"scores"`
	Graded     int                 `json:"graded"`
	Total      int                 `json:"total"`
	FinalGrade float64             `json:"final_grade"`
}

type jsonReport struct {
	Policy        string          `json:"policy"`
	GeneratedAt   time.Time       `json:"generated_at"`
	Criteria      []jsonCriterion `json:"criteria"`
	Students      []jsonStudent   `json:"students"`
	CourseAverage float64         `json:"course_average"`
}

//--------------------------------------

type JSONExporter struct{}

func NewJSONExporter() *JSONExporter {
	return &JSONExporter{}
}

func (e *JSONExporter) GetFormat() string { return "json" }

func (e *JSONExporter) Export(report application.GradeReport, path string) error {
	data, err := json.MarshalIndent(e.toDTO(report), "", "  ")
	if err != nil {
		return fmt.Errorf("не удалось сформировать JSON: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("не удалось записать файл %s: %w", path, err)
	}
	return nil
}

func (e *JSONExporter) toDTO(report application.GradeReport) jsonReport {
	dto := jsonReport{
		Policy:        report.PolicyName,
		GeneratedAt:   report.GeneratedAt,
		Criteria:      make([]jsonCriterion, 0, len(report.Criteria)),
		Students:      make([]jsonStudent, 0, len(report.Rows)),
		CourseAverage: report.CourseAverage,
	}

	for _, criterion := range report.Criteria {
		dto.Criteria = append(dto.Criteria, jsonCriterion{
			ID:       criterion.ID,
			Name:     criterion.Name,
			MaxScore: criterion.MaxScore,
			Weight:   criterion.Weight,
		})
	}

	for _, row := range report.Rows {
		scores := make(map[string]*float64, len(row.Scores))
		for i, cell := range row.Scores {
			if cell.Graded {
				value := cell.Value
				scores[report.Criteria[i].ID] = &value
			} else {
				scores[report.Criteria[i].ID] = nil
			}
		}

		dto.Students = append(dto.Students, jsonStudent{
			ID:         row.StudentID,
			FullName:   row.FullName,
			Group:      row.Group,
			Scores:     scores,
			Graded:     row.GradedCount,
			Total:      row.TotalCriteria,
			FinalGrade: row.FinalGrade,
		})
	}

	return dto
}
