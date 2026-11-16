package application

import (
	"fmt"
	"time"

	"pos-app/internal/domain"
)

var fixedNow = time.Date(2026, 11, 15, 12, 0, 0, 0, time.UTC)

func fixedClock() time.Time { return fixedNow }

//--------------------------------------

type MockStudentRepository struct {
	students   []*domain.Student
	saveCalled int
}

func (m *MockStudentRepository) NextID() domain.StudentID {
	return domain.StudentID(fmt.Sprintf("S-%d", len(m.students)+1))
}

func (m *MockStudentRepository) Save(student *domain.Student) error {
	m.saveCalled++
	for i, existing := range m.students {
		if existing.GetID() == student.GetID() {
			m.students[i] = student
			return nil
		}
	}
	m.students = append(m.students, student)
	return nil
}

func (m *MockStudentRepository) FindByID(id domain.StudentID) (*domain.Student, error) {
	for _, student := range m.students {
		if student.GetID() == id {
			return student, nil
		}
	}
	return nil, domain.ErrStudentNotFound
}

func (m *MockStudentRepository) FindAll() []*domain.Student {
	return m.students
}

func (m *MockStudentRepository) Delete(id domain.StudentID) error {
	for i, student := range m.students {
		if student.GetID() == id {
			m.students = append(m.students[:i], m.students[i+1:]...)
			return nil
		}
	}
	return domain.ErrStudentNotFound
}

//--------------------------------------

type MockRubricRepository struct {
	rubric     *domain.Rubric
	saveCalled int
}

func (m *MockRubricRepository) Get() *domain.Rubric {
	if m.rubric == nil {
		m.rubric = domain.NewRubric()
	}
	return m.rubric
}

func (m *MockRubricRepository) Save(rubric *domain.Rubric) error {
	m.saveCalled++
	m.rubric = rubric
	return nil
}

//--------------------------------------

type MockEvaluationRepository struct {
	evaluations        []*domain.Evaluation
	saveCalled         int
	deletedByStudent   []domain.StudentID
	deletedByCriterion []domain.CriterionID
}

func (m *MockEvaluationRepository) NextID() domain.EvaluationID {
	return domain.EvaluationID(fmt.Sprintf("E-%d", len(m.evaluations)+1))
}

func (m *MockEvaluationRepository) Save(evaluation *domain.Evaluation) error {
	m.saveCalled++
	for i, existing := range m.evaluations {
		if existing.GetID() == evaluation.GetID() {
			m.evaluations[i] = evaluation
			return nil
		}
	}
	m.evaluations = append(m.evaluations, evaluation)
	return nil
}

func (m *MockEvaluationRepository) FindByID(id domain.EvaluationID) (*domain.Evaluation, error) {
	for _, evaluation := range m.evaluations {
		if evaluation.GetID() == id {
			return evaluation, nil
		}
	}
	return nil, domain.ErrEvaluationNotFound
}

func (m *MockEvaluationRepository) FindByStudent(studentID domain.StudentID) []*domain.Evaluation {
	var result []*domain.Evaluation
	for _, evaluation := range m.evaluations {
		if evaluation.GetStudentID() == studentID {
			result = append(result, evaluation)
		}
	}
	return result
}

func (m *MockEvaluationRepository) FindByStudentAndCriterion(studentID domain.StudentID, criterionID domain.CriterionID) (*domain.Evaluation, bool) {
	for _, evaluation := range m.evaluations {
		if evaluation.GetStudentID() == studentID && evaluation.GetCriterionID() == criterionID {
			return evaluation, true
		}
	}
	return nil, false
}

func (m *MockEvaluationRepository) FindAll() []*domain.Evaluation {
	return m.evaluations
}

func (m *MockEvaluationRepository) Delete(id domain.EvaluationID) error {
	for i, evaluation := range m.evaluations {
		if evaluation.GetID() == id {
			m.evaluations = append(m.evaluations[:i], m.evaluations[i+1:]...)
			return nil
		}
	}
	return domain.ErrEvaluationNotFound
}

func (m *MockEvaluationRepository) DeleteByStudent(studentID domain.StudentID) {
	m.deletedByStudent = append(m.deletedByStudent, studentID)
}

func (m *MockEvaluationRepository) DeleteByCriterion(criterionID domain.CriterionID) {
	m.deletedByCriterion = append(m.deletedByCriterion, criterionID)
}

//--------------------------------------

type MockExporter struct {
	format   string
	exported *GradeReport
	path     string
}

func (m *MockExporter) GetFormat() string { return m.format }

func (m *MockExporter) Export(report GradeReport, path string) error {
	m.exported = &report
	m.path = path
	return nil
}
