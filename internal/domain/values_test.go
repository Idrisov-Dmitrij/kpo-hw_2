package domain

import (
	"errors"
	"math"
	"testing"
)

func TestNewScore_RejectsNegativeAndNaN(t *testing.T) {
	for _, value := range []float64{-0.01, math.NaN(), math.Inf(1)} {
		if _, err := NewScore(value); !errors.Is(err, ErrInvalidScore) {
			t.Errorf("Ожидалась ошибка ErrInvalidScore для %v, получено %v", value, err)
		}
	}
}

func TestNewScore_AcceptsZero(t *testing.T) {
	score, err := NewScore(0)
	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if score.GetValue() != 0 {
		t.Errorf("Ожидался балл 0, получено %f", score.GetValue())
	}
}

func TestNewMaxScore_MustBePositive(t *testing.T) {
	for _, value := range []float64{0, -5} {
		if _, err := NewMaxScore(value); !errors.Is(err, ErrInvalidMaxScore) {
			t.Errorf("Ожидалась ошибка ErrInvalidMaxScore для %v, получено %v", value, err)
		}
	}
}

func TestNewWeight_MustBePositive(t *testing.T) {
	for _, value := range []float64{0, -0.3} {
		if _, err := NewWeight(value); !errors.Is(err, ErrInvalidWeight) {
			t.Errorf("Ожидалась ошибка ErrInvalidWeight для %v, получено %v", value, err)
		}
	}
}

func TestMaxScore_ContainsBoundary(t *testing.T) {
	max, _ := NewMaxScore(10)
	ten, _ := NewScore(10)
	tooMuch, _ := NewScore(10.01)

	if !max.Contains(ten) {
		t.Error("Балл, равный максимуму, должен быть допустим")
	}
	if max.Contains(tooMuch) {
		t.Error("Балл больше максимума должен быть недопустим")
	}
}

func TestNewFinalGrade_Range(t *testing.T) {
	if _, err := NewFinalGrade(10.5); !errors.Is(err, ErrInvalidFinalGrade) {
		t.Errorf("Ожидалась ошибка ErrInvalidFinalGrade, получено %v", err)
	}
	if _, err := NewFinalGrade(-1); !errors.Is(err, ErrInvalidFinalGrade) {
		t.Errorf("Ожидалась ошибка ErrInvalidFinalGrade, получено %v", err)
	}
	if _, err := NewFinalGrade(7.25); err != nil {
		t.Errorf("Не ожидалась ошибка, получено: %v", err)
	}
}

func TestNewStudent_RequiresName(t *testing.T) {
	if _, err := NewStudent("S-1", "   ", "БПИ"); !errors.Is(err, ErrEmptyStudentName) {
		t.Errorf("Ожидалась ошибка ErrEmptyStudentName, получено %v", err)
	}
	if _, err := NewStudent("", "Иванов", "БПИ"); !errors.Is(err, ErrEmptyID) {
		t.Errorf("Ожидалась ошибка ErrEmptyID, получено %v", err)
	}
}

func TestStudent_RenameKeepsInvariant(t *testing.T) {
	student, _ := NewStudent("S-1", "Иванов", "БПИ")

	if err := student.Rename(""); !errors.Is(err, ErrEmptyStudentName) {
		t.Errorf("Ожидалась ошибка ErrEmptyStudentName, получено %v", err)
	}
	if student.GetFullName() != "Иванов" {
		t.Errorf("ФИО не должно измениться после неудачного переименования, получено %s", student.GetFullName())
	}
}
