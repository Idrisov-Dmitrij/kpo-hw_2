package application

import (
	"fmt"
	"strings"

	"pos-app/internal/domain"
)

type ReportService struct {
	students    domain.StudentRepository
	rubrics     domain.RubricRepository
	evaluations domain.EvaluationRepository
	calculator  *domain.GradeCalculator
	policies    []domain.GradingPolicy
	exporters   []ReportExporter
	clock       Clock
}

func NewReportService(
	students domain.StudentRepository,
	rubrics domain.RubricRepository,
	evaluations domain.EvaluationRepository,
	calculator *domain.GradeCalculator,
	policies []domain.GradingPolicy,
	exporters []ReportExporter,
	clock Clock,
) *ReportService {
	return &ReportService{
		students:    students,
		rubrics:     rubrics,
		evaluations: evaluations,
		calculator:  calculator,
		policies:    policies,
		exporters:   exporters,
		clock:       clock,
	}
}

func (s *ReportService) StudentGrade(studentID string) (StudentGradeView, error) {
	student, err := s.students.FindByID(domain.StudentID(studentID))
	if err != nil {
		return StudentGradeView{}, err
	}
	return s.gradeStudent(student, s.rubrics.Get())
}

func (s *ReportService) BuildReport() (GradeReport, error) {
	rubric := s.rubrics.Get()
	report := GradeReport{
		PolicyName:  s.calculator.GetPolicy().GetName(),
		GeneratedAt: s.clock(),
		Criteria:    make([]CriterionView, 0),
		Rows:        make([]StudentGradeView, 0),
	}

	for _, criterion := range rubric.GetCriteria() {
		report.Criteria = append(report.Criteria, toCriterionView(criterion))
	}

	var total float64
	for _, student := range s.students.FindAll() {
		row, err := s.gradeStudent(student, rubric)
		if err != nil {
			return GradeReport{}, err
		}
		report.Rows = append(report.Rows, row)
		total += row.FinalGrade
	}

	if len(report.Rows) > 0 {
		report.CourseAverage = total / float64(len(report.Rows))
	}
	return report, nil
}

func (s *ReportService) Export(format, path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return ErrEmptyExportPath
	}

	for _, exporter := range s.exporters {
		if strings.EqualFold(exporter.GetFormat(), strings.TrimSpace(format)) {
			report, err := s.BuildReport()
			if err != nil {
				return err
			}
			return exporter.Export(report, path)
		}
	}
	return fmt.Errorf("%w: %s", ErrUnknownExportFormat, format)
}

func (s *ReportService) ExportFormats() []string {
	result := make([]string, 0, len(s.exporters))
	for _, exporter := range s.exporters {
		result = append(result, exporter.GetFormat())
	}
	return result
}

func (s *ReportService) Policies() []string {
	result := make([]string, 0, len(s.policies))
	for _, policy := range s.policies {
		result = append(result, policy.GetName())
	}
	return result
}

func (s *ReportService) CurrentPolicy() string {
	return s.calculator.GetPolicy().GetName()
}

func (s *ReportService) SelectPolicy(index int) error {
	if index < 0 || index >= len(s.policies) {
		return fmt.Errorf("%w: %d", ErrUnknownPolicy, index+1)
	}
	return s.calculator.ChangePolicy(s.policies[index])
}

func (s *ReportService) gradeStudent(student *domain.Student, rubric *domain.Rubric) (StudentGradeView, error) {
	evaluations := s.evaluations.FindByStudent(student.GetID())
	grade, err := s.calculator.Calculate(student, rubric, evaluations)
	if err != nil {
		return StudentGradeView{}, err
	}

	byCriterion := make(map[domain.CriterionID]domain.Score, len(evaluations))
	for _, evaluation := range evaluations {
		byCriterion[evaluation.GetCriterionID()] = evaluation.GetScore()
	}

	criteria := rubric.GetCriteria()
	cells := make([]ScoreCell, 0, len(criteria))
	for _, criterion := range criteria {
		score, ok := byCriterion[criterion.GetID()]
		cells = append(cells, ScoreCell{Value: score.GetValue(), Graded: ok})
	}

	return StudentGradeView{
		StudentID:     string(student.GetID()),
		FullName:      student.GetFullName(),
		Group:         student.GetGroup(),
		Scores:        cells,
		GradedCount:   len(evaluations),
		TotalCriteria: len(criteria),
		FinalGrade:    grade.GetValue(),
	}, nil
}
