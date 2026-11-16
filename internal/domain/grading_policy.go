package domain

import (
	"fmt"
	"math"
)

type GradingPolicy interface {
	GetName() string
	Calculate(criteria []*Criterion, scores map[CriterionID]Score) float64
}

//--------------------------------------

type WeightedAveragePolicy struct{}

func NewWeightedAveragePolicy() *WeightedAveragePolicy {
	return &WeightedAveragePolicy{}
}

func (p *WeightedAveragePolicy) GetName() string {
	return "Взвешенное среднее по всему рубрикатору (нет оценки = 0)"
}

func (p *WeightedAveragePolicy) Calculate(criteria []*Criterion, scores map[CriterionID]Score) float64 {
	var weighted, totalWeight float64

	for _, criterion := range criteria {
		weight := criterion.GetWeight().GetValue()
		totalWeight += weight

		if score, ok := scores[criterion.GetID()]; ok {
			weighted += weight * criterion.GetMaxScore().Normalize(score)
		}
	}

	if totalWeight == 0 {
		return MinFinalGrade
	}
	return weighted / totalWeight * MaxFinalGrade
}

//--------------------------------------

type CurrentProgressPolicy struct{}

func NewCurrentProgressPolicy() *CurrentProgressPolicy {
	return &CurrentProgressPolicy{}
}

func (p *CurrentProgressPolicy) GetName() string {
	return "Текущая успеваемость (только оценённые критерии)"
}

func (p *CurrentProgressPolicy) Calculate(criteria []*Criterion, scores map[CriterionID]Score) float64 {
	var weighted, gradedWeight float64

	for _, criterion := range criteria {
		score, ok := scores[criterion.GetID()]
		if !ok {
			continue
		}

		weight := criterion.GetWeight().GetValue()
		gradedWeight += weight
		weighted += weight * criterion.GetMaxScore().Normalize(score)
	}

	if gradedWeight == 0 {
		return MinFinalGrade
	}
	return weighted / gradedWeight * MaxFinalGrade
}

//--------------------------------------

const maxRoundingPrecision = 4

type RoundedPolicy struct {
	inner     GradingPolicy
	precision int
}

func NewRoundedPolicy(inner GradingPolicy, precision int) (*RoundedPolicy, error) {
	if inner == nil {
		return nil, ErrGradingPolicyRequired
	}
	if precision < 0 || precision > maxRoundingPrecision {
		return nil, fmt.Errorf("%w: %d", ErrInvalidRoundingPrecision, precision)
	}

	return &RoundedPolicy{inner: inner, precision: precision}, nil
}

func (p *RoundedPolicy) GetName() string {
	switch p.precision {
	case 0:
		return p.inner.GetName() + ", округление до целого"
	case 1:
		return p.inner.GetName() + ", округление до десятых"
	case 2:
		return p.inner.GetName() + ", округление до сотых"
	default:
		return fmt.Sprintf("%s, округление до %d знаков после запятой", p.inner.GetName(), p.precision)
	}
}

func (p *RoundedPolicy) Calculate(criteria []*Criterion, scores map[CriterionID]Score) float64 {
	factor := math.Pow(10, float64(p.precision))
	return math.Floor(p.inner.Calculate(criteria, scores)*factor+0.5+gradeEpsilon) / factor
}
