package domain

import (
	"errors"
	"testing"
)

func mustMax(t *testing.T, value float64) MaxScore {
	t.Helper()
	max, err := NewMaxScore(value)
	if err != nil {
		t.Fatalf("Некорректный максимум в тесте: %v", err)
	}
	return max
}

func mustWeight(t *testing.T, value float64) Weight {
	t.Helper()
	weight, err := NewWeight(value)
	if err != nil {
		t.Fatalf("Некорректный вес в тесте: %v", err)
	}
	return weight
}

func TestRubric_AddCriterion(t *testing.T) {
	rubric := NewRubric()

	criterion, err := rubric.AddCriterion("Домашние задания", mustMax(t, 10), mustWeight(t, 0.5))
	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if criterion.GetID() != "C-1" {
		t.Errorf("Ожидался ID C-1, получено %s", criterion.GetID())
	}
	if len(rubric.GetCriteria()) != 1 {
		t.Errorf("Ожидался 1 критерий, получено %d", len(rubric.GetCriteria()))
	}
}

func TestRubric_RejectsDuplicateNameIgnoringCase(t *testing.T) {
	rubric := NewRubric()
	rubric.AddCriterion("Экзамен", mustMax(t, 100), mustWeight(t, 0.5))

	_, err := rubric.AddCriterion("  экзамен ", mustMax(t, 10), mustWeight(t, 0.1))
	if !errors.Is(err, ErrDuplicateCriterionName) {
		t.Errorf("Ожидалась ошибка ErrDuplicateCriterionName, получено %v", err)
	}
}

func TestRubric_RejectsEmptyName(t *testing.T) {
	rubric := NewRubric()

	if _, err := rubric.AddCriterion(" ", mustMax(t, 10), mustWeight(t, 1)); !errors.Is(err, ErrEmptyCriterionName) {
		t.Errorf("Ожидалась ошибка ErrEmptyCriterionName, получено %v", err)
	}
	if !rubric.IsEmpty() {
		t.Error("Рубрикатор должен остаться пустым")
	}
}

func TestRubric_IDsAreNotReusedAfterRemoval(t *testing.T) {
	rubric := NewRubric()
	first, _ := rubric.AddCriterion("ДЗ", mustMax(t, 10), mustWeight(t, 1))
	rubric.RemoveCriterion(first.GetID())

	second, _ := rubric.AddCriterion("КР", mustMax(t, 10), mustWeight(t, 1))
	if second.GetID() == first.GetID() {
		t.Errorf("ID удалённого критерия не должен переиспользоваться: %s", second.GetID())
	}
}

func TestRubric_GetUnknownCriterion(t *testing.T) {
	rubric := NewRubric()

	if _, err := rubric.GetCriterion("C-42"); !errors.Is(err, ErrCriterionNotInRubric) {
		t.Errorf("Ожидалась ошибка ErrCriterionNotInRubric, получено %v", err)
	}
	if err := rubric.RemoveCriterion("C-42"); !errors.Is(err, ErrCriterionNotInRubric) {
		t.Errorf("Ожидалась ошибка ErrCriterionNotInRubric, получено %v", err)
	}
}

func TestRubric_ChangeWeight(t *testing.T) {
	rubric := NewRubric()
	criterion, _ := rubric.AddCriterion("ДЗ", mustMax(t, 10), mustWeight(t, 1))

	if err := rubric.ChangeWeight(criterion.GetID(), mustWeight(t, 3)); err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if criterion.GetWeight().GetValue() != 3 {
		t.Errorf("Ожидался вес 3, получено %f", criterion.GetWeight().GetValue())
	}
}

func TestRubric_GetCriteriaReturnsCopy(t *testing.T) {
	rubric := NewRubric()
	rubric.AddCriterion("ДЗ", mustMax(t, 10), mustWeight(t, 1))

	criteria := rubric.GetCriteria()
	criteria[0] = nil

	if rubric.GetCriteria()[0] == nil {
		t.Error("Изменение возвращённого среза не должно ломать рубрикатор")
	}
}
