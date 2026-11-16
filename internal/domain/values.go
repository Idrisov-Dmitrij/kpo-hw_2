package domain

import (
	"fmt"
	"math"
)

type StudentID string
type CriterionID string
type EvaluationID string

//--------------------------------------

type Score struct {
	value float64
}

func NewScore(value float64) (Score, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
		return Score{}, fmt.Errorf("%w: %v", ErrInvalidScore, value)
	}
	return Score{value: value}, nil
}

func (s Score) GetValue() float64 { return s.value }

//--------------------------------------

type MaxScore struct {
	value float64
}

func NewMaxScore(value float64) (MaxScore, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || value <= 0 {
		return MaxScore{}, fmt.Errorf("%w: %v", ErrInvalidMaxScore, value)
	}
	return MaxScore{value: value}, nil
}

func (m MaxScore) GetValue() float64 { return m.value }

func (m MaxScore) Contains(score Score) bool {
	return score.value <= m.value
}

func (m MaxScore) Normalize(score Score) float64 {
	return score.value / m.value
}

//--------------------------------------

type Weight struct {
	value float64
}

func NewWeight(value float64) (Weight, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || value <= 0 {
		return Weight{}, fmt.Errorf("%w: %v", ErrInvalidWeight, value)
	}
	return Weight{value: value}, nil
}

func (w Weight) GetValue() float64 { return w.value }

//--------------------------------------

const (
	MinFinalGrade = 0.0
	MaxFinalGrade = 10.0
	gradeEpsilon  = 1e-9
)

type FinalGrade struct {
	value float64
}

func NewFinalGrade(value float64) (FinalGrade, error) {
	if math.IsNaN(value) || value < MinFinalGrade-gradeEpsilon || value > MaxFinalGrade+gradeEpsilon {
		return FinalGrade{}, fmt.Errorf("%w: %v", ErrInvalidFinalGrade, value)
	}
	return FinalGrade{value: math.Min(math.Max(value, MinFinalGrade), MaxFinalGrade)}, nil
}

func (g FinalGrade) GetValue() float64 { return g.value }
