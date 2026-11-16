package export

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"pos-app/internal/application"
)

var _ application.ReportExporter = (*CSVExporter)(nil)
var _ application.ReportExporter = (*JSONExporter)(nil)

func sampleReport() application.GradeReport {
	return application.GradeReport{
		PolicyName:  "Взвешенное среднее",
		GeneratedAt: time.Date(2026, 11, 15, 12, 0, 0, 0, time.UTC),
		Criteria: []application.CriterionView{
			{ID: "C-1", Name: "ДЗ", MaxScore: 10, Weight: 0.5},
			{ID: "C-2", Name: "Экзамен", MaxScore: 100, Weight: 0.5},
		},
		Rows: []application.StudentGradeView{
			{
				StudentID: "S-1", FullName: "Иванов Иван", Group: "БПИ-241",
				Scores:      []application.ScoreCell{{Value: 9, Graded: true}, {}},
				GradedCount: 1, TotalCriteria: 2, FinalGrade: 4.5,
			},
		},
		CourseAverage: 4.5,
	}
}

func TestCSVExporter_WritesHeaderAndRows(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.csv")

	if err := NewCSVExporter().Export(sampleReport(), path); err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}

	file, _ := os.Open(path)
	defer file.Close()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatalf("CSV должен читаться: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("Ожидалось 2 строки (заголовок + студент), получено %d", len(records))
	}
	row := records[1]
	if row[0] != "S-1" || row[3] != "9" || row[4] != "" || row[5] != "1/2" || row[6] != "4.50" {
		t.Errorf("Строка CSV сформирована неверно: %v", row)
	}
}

func TestJSONExporter_WritesNullForUngraded(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.json")

	if err := NewJSONExporter().Export(sampleReport(), path); err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}

	data, _ := os.ReadFile(path)
	var decoded jsonReport
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("JSON должен читаться: %v", err)
	}

	scores := decoded.Students[0].Scores
	if scores["C-1"] == nil || *scores["C-1"] != 9 {
		t.Error("Оценённый критерий должен содержать балл")
	}
	if scores["C-2"] != nil {
		t.Error("Неоценённый критерий должен быть null")
	}
	if decoded.CourseAverage != 4.5 {
		t.Errorf("Ожидался средний итог 4.5, получено %f", decoded.CourseAverage)
	}
}

func TestExporters_FailOnBadPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing-dir", "report")

	if err := NewCSVExporter().Export(sampleReport(), path); err == nil {
		t.Error("Ожидалась ошибка записи CSV в несуществующую папку")
	}
	if err := NewJSONExporter().Export(sampleReport(), path); err == nil {
		t.Error("Ожидалась ошибка записи JSON в несуществующую папку")
	}
}
