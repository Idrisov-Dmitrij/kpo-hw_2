package domain

import (
	"errors"
	"math"
	"testing"
)

var _ GradingPolicy = (*WeightedAveragePolicy)(nil)
var _ GradingPolicy = (*CurrentProgressPolicy)(nil)
var _ GradingPolicy = (*RoundedPolicy)(nil)

type fixedPolicy struct {
	value float64
}

func (p *fixedPolicy) GetName() string { return "fixed" }

func (p *fixedPolicy) Calculate(criteria []*Criterion, scores map[CriterionID]Score) float64 {
	return p.value
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func gradeFor(t *testing.T, policy GradingPolicy, rubric *Rubric, evaluations ...*Evaluation) float64 {
	t.Helper()
	calculator, _ := NewGradeCalculator(policy)
	grade, err := calculator.Calculate(newTestStudent(t, "S-1"), rubric, evaluations)
	if err != nil {
		t.Fatalf("Не ожидалась ошибка расчёта итога: %v", err)
	}
	return grade.GetValue()
}

func TestWeightedAveragePolicy_AllGraded(t *testing.T) {
	rubric, homework, exam := newTestRubric(t)
	student := newTestStudent(t, "S-1")
	factory := NewEvaluationFactory()
	e1, _ := factory.Create("E-1", rubric, student, homework.GetID(), 8, "", testTime)
	e2, _ := factory.Create("E-2", rubric, student, exam.GetID(), 50, "", testTime)

	got := gradeFor(t, NewWeightedAveragePolicy(), rubric, e1, e2)
	expected := (0.6*0.8 + 0.4*0.5) * 10

	if !almostEqual(got, expected) {
		t.Errorf("Ожидался итог %f, получено %f", expected, got)
	}
}

func TestWeightedAveragePolicy_MissingCountsAsZero(t *testing.T) {
	rubric, homework, _ := newTestRubric(t)
	student := newTestStudent(t, "S-1")
	e1, _ := NewEvaluationFactory().Create("E-1", rubric, student, homework.GetID(), 10, "", testTime)

	got := gradeFor(t, NewWeightedAveragePolicy(), rubric, e1)
	if !almostEqual(got, 6) {
		t.Errorf("Ожидался итог 6 (экзамен = 0), получено %f", got)
	}
}

func TestCurrentProgressPolicy_IgnoresMissing(t *testing.T) {
	rubric, homework, _ := newTestRubric(t)
	student := newTestStudent(t, "S-1")
	e1, _ := NewEvaluationFactory().Create("E-1", rubric, student, homework.GetID(), 7, "", testTime)

	got := gradeFor(t, NewCurrentProgressPolicy(), rubric, e1)
	if !almostEqual(got, 7) {
		t.Errorf("Ожидался итог 7 по единственному оценённому критерию, получено %f", got)
	}
}

func TestPolicies_EmptyRubricOrNoEvaluations(t *testing.T) {
	if got := gradeFor(t, NewWeightedAveragePolicy(), NewRubric()); got != 0 {
		t.Errorf("Ожидался итог 0 для пустого рубрикатора, получено %f", got)
	}

	rubric, _, _ := newTestRubric(t)
	if got := gradeFor(t, NewCurrentProgressPolicy(), rubric); got != 0 {
		t.Errorf("Ожидался итог 0 без оценок, получено %f", got)
	}
}

func TestRoundedPolicy_DecoratesInnerPolicy(t *testing.T) {
	cases := []struct {
		value     float64
		precision int
		expected  float64
	}{
		{7.5, 0, 8},
		{7.49, 0, 7},
		{7.956, 2, 7.96},
		{7.954, 2, 7.95},
		{3.45, 1, 3.5},
	}

	for _, tc := range cases {
		policy, err := NewRoundedPolicy(&fixedPolicy{value: tc.value}, tc.precision)
		if err != nil {
			t.Fatalf("Не ожидалась ошибка, получено: %v", err)
		}
		if got := policy.Calculate(nil, nil); !almostEqual(got, tc.expected) {
			t.Errorf("Округление %v до %d знаков: ожидалось %v, получено %v", tc.value, tc.precision, tc.expected, got)
		}
	}
}

func TestRoundedPolicy_CanBeStacked(t *testing.T) {
	inner, _ := NewRoundedPolicy(&fixedPolicy{value: 7.449}, 2)
	outer, _ := NewRoundedPolicy(inner, 1)

	if got := outer.Calculate(nil, nil); !almostEqual(got, 7.5) {
		t.Errorf("Ожидалось 7.5 (7.449 -> 7.45 -> 7.5), получено %v", got)
	}
}

func TestRoundedPolicy_Validation(t *testing.T) {
	if _, err := NewRoundedPolicy(nil, 2); !errors.Is(err, ErrGradingPolicyRequired) {
		t.Errorf("Ожидалась ошибка ErrGradingPolicyRequired, получено %v", err)
	}
	if _, err := NewRoundedPolicy(&fixedPolicy{}, 5); !errors.Is(err, ErrInvalidRoundingPrecision) {
		t.Errorf("Ожидалась ошибка ErrInvalidRoundingPrecision, получено %v", err)
	}
}

func TestGradeCalculator_ChangePolicy(t *testing.T) {
	calculator, _ := NewGradeCalculator(&fixedPolicy{value: 5})
	if err := calculator.ChangePolicy(&fixedPolicy{value: 9}); err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}

	grade, _ := calculator.Calculate(newTestStudent(t, "S-1"), NewRubric(), nil)
	if grade.GetValue() != 9 {
		t.Errorf("Ожидался итог по новой стратегии 9, получено %f", grade.GetValue())
	}
	if err := calculator.ChangePolicy(nil); !errors.Is(err, ErrGradingPolicyRequired) {
		t.Errorf("Ожидалась ошибка ErrGradingPolicyRequired, получено %v", err)
	}
}

func TestGradeCalculator_RejectsForeignEvaluation(t *testing.T) {
	rubric, homework, _ := newTestRubric(t)
	other := newTestStudent(t, "S-2")
	foreign, _ := NewEvaluationFactory().Create("E-1", rubric, other, homework.GetID(), 5, "", testTime)

	calculator, _ := NewGradeCalculator(NewWeightedAveragePolicy())
	_, err := calculator.Calculate(newTestStudent(t, "S-1"), rubric, []*Evaluation{foreign})
	if !errors.Is(err, ErrForeignEvaluation) {
		t.Errorf("Ожидалась ошибка ErrForeignEvaluation, получено %v", err)
	}
}

func TestGradeCalculator_RejectsEvaluationForRemovedCriterion(t *testing.T) {
	rubric, homework, _ := newTestRubric(t)
	student := newTestStudent(t, "S-1")
	evaluation, _ := NewEvaluationFactory().Create("E-1", rubric, student, homework.GetID(), 5, "", testTime)
	rubric.RemoveCriterion(homework.GetID())

	calculator, _ := NewGradeCalculator(NewWeightedAveragePolicy())
	_, err := calculator.Calculate(student, rubric, []*Evaluation{evaluation})
	if !errors.Is(err, ErrCriterionNotInRubric) {
		t.Errorf("Ожидалась ошибка ErrCriterionNotInRubric, получено %v", err)
	}
}

func TestGradeCalculator_RejectsDuplicateEvaluations(t *testing.T) {
	rubric, homework, _ := newTestRubric(t)
	student := newTestStudent(t, "S-1")
	factory := NewEvaluationFactory()
	e1, _ := factory.Create("E-1", rubric, student, homework.GetID(), 5, "", testTime)
	e2, _ := factory.Create("E-2", rubric, student, homework.GetID(), 7, "", testTime)

	calculator, _ := NewGradeCalculator(NewWeightedAveragePolicy())
	_, err := calculator.Calculate(student, rubric, []*Evaluation{e1, e2})
	if !errors.Is(err, ErrDuplicateEvaluation) {
		t.Errorf("Ожидалась ошибка ErrDuplicateEvaluation, получено %v", err)
	}
}

func TestGradeCalculator_RequiresPolicy(t *testing.T) {
	if _, err := NewGradeCalculator(nil); !errors.Is(err, ErrGradingPolicyRequired) {
		t.Errorf("Ожидалась ошибка ErrGradingPolicyRequired, получено %v", err)
	}
}

func TestGradeCalculator_RejectsPolicyResultOutOfScale(t *testing.T) {
	calculator, _ := NewGradeCalculator(&fixedPolicy{value: 12})

	_, err := calculator.Calculate(newTestStudent(t, "S-1"), NewRubric(), nil)
	if !errors.Is(err, ErrInvalidFinalGrade) {
		t.Errorf("Ожидалась ошибка ErrInvalidFinalGrade, получено %v", err)
	}
}
