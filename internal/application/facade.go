package application

type PosFacade struct {
	students *StudentService
	rubric   *RubricService
	grading  *GradingService
	reports  *ReportService
}

func NewPosFacade(students *StudentService, rubric *RubricService, grading *GradingService, reports *ReportService) *PosFacade {
	return &PosFacade{
		students: students,
		rubric:   rubric,
		grading:  grading,
		reports:  reports,
	}
}

func (f *PosFacade) IsEmpty() bool {
	return len(f.students.List()) == 0 && len(f.rubric.List()) == 0
}

func (f *PosFacade) RegisterStudent(fullName, group string) (StudentView, error) {
	return f.students.Register(fullName, group)
}

func (f *PosFacade) RenameStudent(id, fullName string) (StudentView, error) {
	return f.students.Rename(id, fullName)
}

func (f *PosFacade) RemoveStudent(id string) error {
	return f.students.Remove(id)
}

func (f *PosFacade) ListStudents() []StudentView {
	return f.students.List()
}

func (f *PosFacade) AddCriterion(name string, maxScore, weight float64) (CriterionView, error) {
	return f.rubric.AddCriterion(name, maxScore, weight)
}

func (f *PosFacade) ChangeCriterionWeight(id string, weight float64) (CriterionView, error) {
	return f.rubric.ChangeWeight(id, weight)
}

func (f *PosFacade) RemoveCriterion(id string) error {
	return f.rubric.RemoveCriterion(id)
}

func (f *PosFacade) ListCriteria() []CriterionView {
	return f.rubric.List()
}

func (f *PosFacade) Assess(studentID, criterionID string, score float64, comment string) (EvaluationView, error) {
	return f.grading.Assess(studentID, criterionID, score, comment)
}

func (f *PosFacade) CorrectEvaluation(evaluationID string, score float64, comment string) (EvaluationView, error) {
	return f.grading.Correct(evaluationID, score, comment)
}

func (f *PosFacade) CancelEvaluation(evaluationID string) error {
	return f.grading.Cancel(evaluationID)
}

func (f *PosFacade) ListEvaluations() []EvaluationView {
	return f.grading.ListAll()
}

func (f *PosFacade) ListStudentEvaluations(studentID string) ([]EvaluationView, error) {
	return f.grading.ListByStudent(studentID)
}

func (f *PosFacade) StudentGrade(studentID string) (StudentGradeView, error) {
	return f.reports.StudentGrade(studentID)
}

func (f *PosFacade) GradeReport() (GradeReport, error) {
	return f.reports.BuildReport()
}

func (f *PosFacade) ExportReport(format, path string) error {
	return f.reports.Export(format, path)
}

func (f *PosFacade) ExportFormats() []string {
	return f.reports.ExportFormats()
}

func (f *PosFacade) GradingPolicies() []string {
	return f.reports.Policies()
}

func (f *PosFacade) CurrentGradingPolicy() string {
	return f.reports.CurrentPolicy()
}

func (f *PosFacade) SelectGradingPolicy(index int) error {
	return f.reports.SelectPolicy(index)
}
