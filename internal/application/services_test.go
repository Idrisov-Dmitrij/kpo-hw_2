package application

import (
	"errors"
	"math"
	"testing"

	"pos-app/internal/domain"
)

type testEnv struct {
	students    *MockStudentRepository
	rubrics     *MockRubricRepository
	evaluations *MockEvaluationRepository
	exporter    *MockExporter
	studentSvc  *StudentService
	rubricSvc   *RubricService
	grading     *GradingService
	reports     *ReportService
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	env := &testEnv{
		students:    &MockStudentRepository{},
		rubrics:     &MockRubricRepository{},
		evaluations: &MockEvaluationRepository{},
		exporter:    &MockExporter{format: "csv"},
	}

	weighted := domain.NewWeightedAveragePolicy()
	progress := domain.NewCurrentProgressPolicy()
	calculator, err := domain.NewGradeCalculator(weighted)
	if err != nil {
		t.Fatalf("Не удалось создать калькулятор: %v", err)
	}

	env.studentSvc = NewStudentService(env.students, env.evaluations)
	env.rubricSvc = NewRubricService(env.rubrics, env.evaluations)
	env.grading = NewGradingService(env.students, env.rubrics, env.evaluations, domain.NewEvaluationFactory(), fixedClock)
	env.reports = NewReportService(env.students, env.rubrics, env.evaluations, calculator,
		[]domain.GradingPolicy{weighted, progress}, []ReportExporter{env.exporter}, fixedClock)
	return env
}

func (env *testEnv) seed(t *testing.T) (StudentView, CriterionView, CriterionView) {
	t.Helper()
	student, err := env.studentSvc.Register("Иванов Иван", "БПИ-241")
	if err != nil {
		t.Fatalf("Не удалось добавить студента: %v", err)
	}
	homework, err := env.rubricSvc.AddCriterion("ДЗ", 10, 0.5)
	if err != nil {
		t.Fatalf("Не удалось добавить критерий: %v", err)
	}
	exam, err := env.rubricSvc.AddCriterion("Экзамен", 100, 0.5)
	if err != nil {
		t.Fatalf("Не удалось добавить критерий: %v", err)
	}
	return student, homework, exam
}

func TestGradingService_AssessSavesEvaluation(t *testing.T) {
	env := newTestEnv(t)
	student, homework, _ := env.seed(t)

	view, err := env.grading.Assess(student.ID, homework.ID, 9, "хорошо")
	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if env.evaluations.saveCalled != 1 {
		t.Error("Оценка должна быть сохранена в репозиторий")
	}
	if view.StudentName != "Иванов Иван" || view.CriterionName != "ДЗ" || view.MaxScore != 10 {
		t.Errorf("Представление оценки заполнено неверно: %+v", view)
	}
	if !view.GradedAt.Equal(fixedNow) {
		t.Errorf("Дата оценки должна браться из Clock, получено %v", view.GradedAt)
	}
}

func TestGradingService_AssessRejectsScoreAboveMax(t *testing.T) {
	env := newTestEnv(t)
	student, homework, _ := env.seed(t)

	_, err := env.grading.Assess(student.ID, homework.ID, 12, "")
	if !errors.Is(err, domain.ErrScoreOutOfRange) {
		t.Errorf("Ожидалась ошибка ErrScoreOutOfRange, получено %v", err)
	}
	if env.evaluations.saveCalled != 0 {
		t.Error("Недопустимая оценка не должна сохраняться")
	}
}

func TestGradingService_AssessRejectsUnknownCriterionAndStudent(t *testing.T) {
	env := newTestEnv(t)
	student, homework, _ := env.seed(t)

	if _, err := env.grading.Assess(student.ID, "C-99", 5, ""); !errors.Is(err, domain.ErrCriterionNotInRubric) {
		t.Errorf("Ожидалась ошибка ErrCriterionNotInRubric, получено %v", err)
	}
	if _, err := env.grading.Assess("S-99", homework.ID, 5, ""); !errors.Is(err, domain.ErrStudentNotFound) {
		t.Errorf("Ожидалась ошибка ErrStudentNotFound, получено %v", err)
	}
	if env.evaluations.saveCalled != 0 {
		t.Error("Ни одна некорректная оценка не должна сохраниться")
	}
}

func TestGradingService_AssessRejectsSecondEvaluationForSameCriterion(t *testing.T) {
	env := newTestEnv(t)
	student, homework, _ := env.seed(t)
	env.grading.Assess(student.ID, homework.ID, 5, "")

	_, err := env.grading.Assess(student.ID, homework.ID, 7, "")
	if !errors.Is(err, domain.ErrDuplicateEvaluation) {
		t.Errorf("Ожидалась ошибка ErrDuplicateEvaluation, получено %v", err)
	}
}

func TestGradingService_CorrectAndCancel(t *testing.T) {
	env := newTestEnv(t)
	student, _, exam := env.seed(t)
	evaluation, _ := env.grading.Assess(student.ID, exam.ID, 59, "")

	corrected, err := env.grading.Correct(evaluation.ID, 69, "апелляция")
	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if corrected.Score != 69 || corrected.Comment != "апелляция" {
		t.Errorf("Оценка исправлена неверно: %+v", corrected)
	}

	if _, err := env.grading.Correct(evaluation.ID, 101, ""); !errors.Is(err, domain.ErrScoreOutOfRange) {
		t.Errorf("Ожидалась ошибка ErrScoreOutOfRange, получено %v", err)
	}

	if err := env.grading.Cancel(evaluation.ID); err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if len(env.grading.ListAll()) != 0 {
		t.Error("После отмены оценок быть не должно")
	}
	if err := env.grading.Cancel(evaluation.ID); !errors.Is(err, domain.ErrEvaluationNotFound) {
		t.Errorf("Ожидалась ошибка ErrEvaluationNotFound, получено %v", err)
	}
}

func TestStudentService_RemoveCascadesEvaluations(t *testing.T) {
	env := newTestEnv(t)
	student, _, _ := env.seed(t)

	if err := env.studentSvc.Remove(student.ID); err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if len(env.evaluations.deletedByStudent) != 1 || string(env.evaluations.deletedByStudent[0]) != student.ID {
		t.Error("При удалении студента должны удаляться его оценки")
	}
}

func TestRubricService_RemoveCriterionCascadesEvaluations(t *testing.T) {
	env := newTestEnv(t)
	_, homework, _ := env.seed(t)

	if err := env.rubricSvc.RemoveCriterion(homework.ID); err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if len(env.evaluations.deletedByCriterion) != 1 {
		t.Error("При удалении критерия должны удаляться оценки по нему")
	}
	if len(env.rubricSvc.List()) != 1 {
		t.Errorf("Ожидался 1 критерий, получено %d", len(env.rubricSvc.List()))
	}
}

func TestRubricService_AddCriterionValidatesInput(t *testing.T) {
	env := newTestEnv(t)

	if _, err := env.rubricSvc.AddCriterion("ДЗ", 0, 1); !errors.Is(err, domain.ErrInvalidMaxScore) {
		t.Errorf("Ожидалась ошибка ErrInvalidMaxScore, получено %v", err)
	}
	if _, err := env.rubricSvc.AddCriterion("ДЗ", 10, -1); !errors.Is(err, domain.ErrInvalidWeight) {
		t.Errorf("Ожидалась ошибка ErrInvalidWeight, получено %v", err)
	}
	if env.rubrics.saveCalled != 0 {
		t.Error("Некорректный критерий не должен сохраняться")
	}
}

func TestReportService_BuildReport(t *testing.T) {
	env := newTestEnv(t)
	ivanov, homework, exam := env.seed(t)
	petrova, _ := env.studentSvc.Register("Петрова Анна", "")
	env.grading.Assess(ivanov.ID, homework.ID, 10, "")
	env.grading.Assess(ivanov.ID, exam.ID, 60, "")
	env.grading.Assess(petrova.ID, homework.ID, 4, "")

	report, err := env.reports.BuildReport()
	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if len(report.Rows) != 2 || len(report.Criteria) != 2 {
		t.Fatalf("Ожидалось 2 строки и 2 критерия, получено %d и %d", len(report.Rows), len(report.Criteria))
	}
	if math.Abs(report.Rows[0].FinalGrade-8) > 1e-9 {
		t.Errorf("Ожидался итог Иванова 8, получено %f", report.Rows[0].FinalGrade)
	}
	if math.Abs(report.Rows[1].FinalGrade-2) > 1e-9 {
		t.Errorf("Ожидался итог Петровой 2, получено %f", report.Rows[1].FinalGrade)
	}
	if report.Rows[1].Scores[1].Graded {
		t.Error("Неоценённый критерий должен быть помечен в ведомости")
	}
	if math.Abs(report.CourseAverage-5) > 1e-9 {
		t.Errorf("Ожидался средний итог 5, получено %f", report.CourseAverage)
	}
}

func TestReportService_SelectPolicyChangesGrade(t *testing.T) {
	env := newTestEnv(t)
	student, homework, _ := env.seed(t)
	env.grading.Assess(student.ID, homework.ID, 4, "")

	if err := env.reports.SelectPolicy(1); err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	grade, _ := env.reports.StudentGrade(student.ID)
	if math.Abs(grade.FinalGrade-4) > 1e-9 {
		t.Errorf("По текущей успеваемости ожидался итог 4, получено %f", grade.FinalGrade)
	}

	if err := env.reports.SelectPolicy(5); !errors.Is(err, ErrUnknownPolicy) {
		t.Errorf("Ожидалась ошибка ErrUnknownPolicy, получено %v", err)
	}
}

func TestReportService_ExportUsesMatchingExporter(t *testing.T) {
	env := newTestEnv(t)
	env.seed(t)

	if err := env.reports.Export("CSV", "out.csv"); err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if env.exporter.exported == nil || env.exporter.path != "out.csv" {
		t.Error("Экспорт должен быть делегирован адаптеру нужного формата")
	}

	if err := env.reports.Export("xml", "out.xml"); !errors.Is(err, ErrUnknownExportFormat) {
		t.Errorf("Ожидалась ошибка ErrUnknownExportFormat, получено %v", err)
	}
	if err := env.reports.Export("csv", "  "); !errors.Is(err, ErrEmptyExportPath) {
		t.Errorf("Ожидалась ошибка ErrEmptyExportPath, получено %v", err)
	}
}
