package domain

import (
	"strings"
	"time"
)

type EvaluationFactory struct{}

func NewEvaluationFactory() *EvaluationFactory {
	return &EvaluationFactory{}
}

func (f *EvaluationFactory) Create(
	id EvaluationID,
	rubric *Rubric,
	student *Student,
	criterionID CriterionID,
	rawScore float64,
	comment string,
	at time.Time,
) (*Evaluation, error) {
	if strings.TrimSpace(string(id)) == "" {
		return nil, ErrEmptyID
	}
	if student == nil {
		return nil, ErrStudentRequired
	}
	if rubric == nil {
		return nil, ErrRubricRequired
	}

	criterion, err := rubric.GetCriterion(criterionID)
	if err != nil {
		return nil, err
	}

	score, err := NewScore(rawScore)
	if err != nil {
		return nil, err
	}
	if err := criterion.Check(score); err != nil {
		return nil, err
	}

	return &Evaluation{
		id:          id,
		studentID:   student.GetID(),
		criterionID: criterion.GetID(),
		score:       score,
		comment:     strings.TrimSpace(comment),
		gradedAt:    at,
	}, nil
}
