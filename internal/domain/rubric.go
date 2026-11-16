package domain

import (
	"fmt"
	"strings"
)

type Rubric struct {
	criteria []*Criterion
	sequence int
}

func NewRubric() *Rubric {
	return &Rubric{
		criteria: make([]*Criterion, 0),
	}
}

func (r *Rubric) AddCriterion(name string, maxScore MaxScore, weight Weight) (*Criterion, error) {
	if r.hasName(name) {
		return nil, fmt.Errorf("%w: «%s»", ErrDuplicateCriterionName, strings.TrimSpace(name))
	}

	id := CriterionID(fmt.Sprintf("C-%d", r.sequence+1))
	criterion, err := newCriterion(id, name, maxScore, weight)
	if err != nil {
		return nil, err
	}

	r.sequence++
	r.criteria = append(r.criteria, criterion)
	return criterion, nil
}

func (r *Rubric) GetCriterion(id CriterionID) (*Criterion, error) {
	for _, criterion := range r.criteria {
		if criterion.id == id {
			return criterion, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrCriterionNotInRubric, id)
}

func (r *Rubric) GetCriteria() []*Criterion {
	result := make([]*Criterion, len(r.criteria))
	copy(result, r.criteria)
	return result
}

func (r *Rubric) ChangeWeight(id CriterionID, weight Weight) error {
	criterion, err := r.GetCriterion(id)
	if err != nil {
		return err
	}

	criterion.changeWeight(weight)
	return nil
}

func (r *Rubric) RemoveCriterion(id CriterionID) error {
	for i, criterion := range r.criteria {
		if criterion.id == id {
			r.criteria = append(r.criteria[:i], r.criteria[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrCriterionNotInRubric, id)
}

func (r *Rubric) IsEmpty() bool {
	return len(r.criteria) == 0
}

func (r *Rubric) hasName(name string) bool {
	name = strings.TrimSpace(name)
	for _, criterion := range r.criteria {
		if strings.EqualFold(criterion.name, name) {
			return true
		}
	}
	return false
}
