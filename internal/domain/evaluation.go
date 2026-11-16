package domain

import (
	"strings"
	"time"
)

type Evaluation struct {
	id          EvaluationID
	studentID   StudentID
	criterionID CriterionID
	score       Score
	comment     string
	gradedAt    time.Time
}

func (e *Evaluation) GetID() EvaluationID         { return e.id }
func (e *Evaluation) GetStudentID() StudentID     { return e.studentID }
func (e *Evaluation) GetCriterionID() CriterionID { return e.criterionID }
func (e *Evaluation) GetScore() Score             { return e.score }
func (e *Evaluation) GetComment() string          { return e.comment }
func (e *Evaluation) GetGradedAt() time.Time      { return e.gradedAt }

func (e *Evaluation) Correct(criterion *Criterion, score Score, comment string, at time.Time) error {
	if criterion == nil || criterion.id != e.criterionID {
		return ErrCriterionMismatch
	}
	if err := criterion.Check(score); err != nil {
		return err
	}

	e.score = score
	e.comment = strings.TrimSpace(comment)
	e.gradedAt = at
	return nil
}
