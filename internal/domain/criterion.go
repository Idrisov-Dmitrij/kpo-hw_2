package domain

import (
	"fmt"
	"strings"
)

type Criterion struct {
	id       CriterionID
	name     string
	maxScore MaxScore
	weight   Weight
}

func newCriterion(id CriterionID, name string, maxScore MaxScore, weight Weight) (*Criterion, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyCriterionName
	}

	return &Criterion{
		id:       id,
		name:     name,
		maxScore: maxScore,
		weight:   weight,
	}, nil
}

func (c *Criterion) GetID() CriterionID    { return c.id }
func (c *Criterion) GetName() string       { return c.name }
func (c *Criterion) GetMaxScore() MaxScore { return c.maxScore }
func (c *Criterion) GetWeight() Weight     { return c.weight }

func (c *Criterion) Check(score Score) error {
	if !c.maxScore.Contains(score) {
		return fmt.Errorf("%w: %.2f > %.2f по критерию «%s»",
			ErrScoreOutOfRange, score.GetValue(), c.maxScore.GetValue(), c.name)
	}
	return nil
}

func (c *Criterion) changeWeight(weight Weight) {
	c.weight = weight
}
