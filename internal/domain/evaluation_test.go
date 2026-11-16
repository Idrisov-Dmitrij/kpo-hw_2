package domain

import (
	"errors"
	"testing"
	"time"
)

var testTime = time.Date(2026, 11, 15, 12, 0, 0, 0, time.UTC)

func newTestRubric(t *testing.T) (*Rubric, *Criterion, *Criterion) {
	t.Helper()
	rubric := NewRubric()
	homework, _ := rubric.AddCriterion("ДЗ", mustMax(t, 10), mustWeight(t, 0.6))
	exam, _ := rubric.AddCriterion("Экзамен", mustMax(t, 100), mustWeight(t, 0.4))
	return rubric, homework, exam
}

func newTestStudent(t *testing.T, id StudentID) *Student {
	t.Helper()
	student, err := NewStudent(id, "Иванов Иван", "БПИ-241")
	if err != nil {
		t.Fatalf("Некорректный студент в тесте: %v", err)
	}
	return student
}

func TestEvaluationFactory_CreatesValidEvaluation(t *testing.T) {
	rubric, homework, _ := newTestRubric(t)
	student := newTestStudent(t, "S-1")

	evaluation, err := NewEvaluationFactory().Create("E-1", rubric, student, homework.GetID(), 10, " отлично ", testTime)
	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if evaluation.GetScore().GetValue() != 10 {
		t.Errorf("Ожидался балл 10, получено %f", evaluation.GetScore().GetValue())
	}
	if evaluation.GetComment() != "отлично" {
		t.Errorf("Ожидался обрезанный комментарий, получено %q", evaluation.GetComment())
	}
	if evaluation.GetStudentID() != "S-1" || evaluation.GetCriterionID() != homework.GetID() {
		t.Error("Оценка должна ссылаться на студента и критерий")
	}
}

func TestEvaluationFactory_RejectsScoreAboveMax(t *testing.T) {
	rubric, homework, _ := newTestRubric(t)
	student := newTestStudent(t, "S-1")

	_, err := NewEvaluationFactory().Create("E-1", rubric, student, homework.GetID(), 12, "", testTime)
	if !errors.Is(err, ErrScoreOutOfRange) {
		t.Errorf("Ожидалась ошибка ErrScoreOutOfRange для 12 из 10, получено %v", err)
	}
}

func TestEvaluationFactory_RejectsNegativeScore(t *testing.T) {
	rubric, homework, _ := newTestRubric(t)
	student := newTestStudent(t, "S-1")

	_, err := NewEvaluationFactory().Create("E-1", rubric, student, homework.GetID(), -1, "", testTime)
	if !errors.Is(err, ErrInvalidScore) {
		t.Errorf("Ожидалась ошибка ErrInvalidScore, получено %v", err)
	}
}

func TestEvaluationFactory_RejectsCriterionOutsideRubric(t *testing.T) {
	rubric, _, _ := newTestRubric(t)
	student := newTestStudent(t, "S-1")

	_, err := NewEvaluationFactory().Create("E-1", rubric, student, "C-99", 5, "", testTime)
	if !errors.Is(err, ErrCriterionNotInRubric) {
		t.Errorf("Ожидалась ошибка ErrCriterionNotInRubric, получено %v", err)
	}
}

func TestEvaluationFactory_RequiresStudentAndRubric(t *testing.T) {
	rubric, homework, _ := newTestRubric(t)
	student := newTestStudent(t, "S-1")
	factory := NewEvaluationFactory()

	if _, err := factory.Create("E-1", rubric, nil, homework.GetID(), 5, "", testTime); !errors.Is(err, ErrStudentRequired) {
		t.Errorf("Ожидалась ошибка ErrStudentRequired, получено %v", err)
	}
	if _, err := factory.Create("E-1", nil, student, homework.GetID(), 5, "", testTime); !errors.Is(err, ErrRubricRequired) {
		t.Errorf("Ожидалась ошибка ErrRubricRequired, получено %v", err)
	}
}

func TestEvaluation_CorrectWithinRange(t *testing.T) {
	rubric, _, exam := newTestRubric(t)
	student := newTestStudent(t, "S-1")
	evaluation, _ := NewEvaluationFactory().Create("E-1", rubric, student, exam.GetID(), 59, "", testTime)

	newScore, _ := NewScore(69)
	later := testTime.Add(time.Hour)
	if err := evaluation.Correct(exam, newScore, "апелляция", later); err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if evaluation.GetScore().GetValue() != 69 || !evaluation.GetGradedAt().Equal(later) {
		t.Error("Балл и дата оценки должны обновиться")
	}
}

func TestEvaluation_CorrectCannotBreakRange(t *testing.T) {
	rubric, _, exam := newTestRubric(t)
	student := newTestStudent(t, "S-1")
	evaluation, _ := NewEvaluationFactory().Create("E-1", rubric, student, exam.GetID(), 59, "", testTime)

	tooMuch, _ := NewScore(120)
	if err := evaluation.Correct(exam, tooMuch, "", testTime); !errors.Is(err, ErrScoreOutOfRange) {
		t.Errorf("Ожидалась ошибка ErrScoreOutOfRange, получено %v", err)
	}
	if evaluation.GetScore().GetValue() != 59 {
		t.Errorf("Балл не должен измениться после отказа, получено %f", evaluation.GetScore().GetValue())
	}
}

func TestEvaluation_CorrectRequiresSameCriterion(t *testing.T) {
	rubric, homework, exam := newTestRubric(t)
	student := newTestStudent(t, "S-1")
	evaluation, _ := NewEvaluationFactory().Create("E-1", rubric, student, exam.GetID(), 59, "", testTime)

	score, _ := NewScore(5)
	if err := evaluation.Correct(homework, score, "", testTime); !errors.Is(err, ErrCriterionMismatch) {
		t.Errorf("Ожидалась ошибка ErrCriterionMismatch, получено %v", err)
	}
}
