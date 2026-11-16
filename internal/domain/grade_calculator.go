package domain

import "fmt"

type GradeCalculator struct {
	policy GradingPolicy
}

func NewGradeCalculator(policy GradingPolicy) (*GradeCalculator, error) {
	if policy == nil {
		return nil, ErrGradingPolicyRequired
	}
	return &GradeCalculator{policy: policy}, nil
}

func (c *GradeCalculator) GetPolicy() GradingPolicy { return c.policy }

func (c *GradeCalculator) ChangePolicy(policy GradingPolicy) error {
	if policy == nil {
		return ErrGradingPolicyRequired
	}

	c.policy = policy
	return nil
}

func (c *GradeCalculator) Calculate(student *Student, rubric *Rubric, evaluations []*Evaluation) (FinalGrade, error) {
	if student == nil {
		return FinalGrade{}, ErrStudentRequired
	}
	if rubric == nil {
		return FinalGrade{}, ErrRubricRequired
	}

	scores := make(map[CriterionID]Score, len(evaluations))
	for _, evaluation := range evaluations {
		if evaluation.GetStudentID() != student.GetID() {
			return FinalGrade{}, fmt.Errorf("%w: %s", ErrForeignEvaluation, evaluation.GetID())
		}

		criterion, err := rubric.GetCriterion(evaluation.GetCriterionID())
		if err != nil {
			return FinalGrade{}, err
		}
		if err := criterion.Check(evaluation.GetScore()); err != nil {
			return FinalGrade{}, err
		}
		if _, exists := scores[criterion.GetID()]; exists {
			return FinalGrade{}, fmt.Errorf("%w: %s", ErrDuplicateEvaluation, criterion.GetName())
		}

		scores[criterion.GetID()] = evaluation.GetScore()
	}

	return NewFinalGrade(c.policy.Calculate(rubric.GetCriteria(), scores))
}
