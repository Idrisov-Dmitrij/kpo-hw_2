package memory

import (
	"fmt"

	"pos-app/internal/domain"
)

type EvaluationRepository struct {
	evaluations []*domain.Evaluation
	sequence    int
}

func NewEvaluationRepository() *EvaluationRepository {
	return &EvaluationRepository{
		evaluations: make([]*domain.Evaluation, 0),
	}
}

func (r *EvaluationRepository) NextID() domain.EvaluationID {
	return domain.EvaluationID(fmt.Sprintf("E-%d", r.sequence+1))
}

func (r *EvaluationRepository) Save(evaluation *domain.Evaluation) error {
	for i, existing := range r.evaluations {
		if existing.GetID() == evaluation.GetID() {
			r.evaluations[i] = evaluation
			return nil
		}
	}

	r.evaluations = append(r.evaluations, evaluation)
	r.sequence++
	return nil
}

func (r *EvaluationRepository) FindByID(id domain.EvaluationID) (*domain.Evaluation, error) {
	for _, evaluation := range r.evaluations {
		if evaluation.GetID() == id {
			return evaluation, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", domain.ErrEvaluationNotFound, id)
}

func (r *EvaluationRepository) FindByStudent(studentID domain.StudentID) []*domain.Evaluation {
	return r.filter(func(e *domain.Evaluation) bool {
		return e.GetStudentID() == studentID
	})
}

func (r *EvaluationRepository) FindByStudentAndCriterion(studentID domain.StudentID, criterionID domain.CriterionID) (*domain.Evaluation, bool) {
	for _, evaluation := range r.evaluations {
		if evaluation.GetStudentID() == studentID && evaluation.GetCriterionID() == criterionID {
			return evaluation, true
		}
	}
	return nil, false
}

func (r *EvaluationRepository) FindAll() []*domain.Evaluation {
	result := make([]*domain.Evaluation, len(r.evaluations))
	copy(result, r.evaluations)
	return result
}

func (r *EvaluationRepository) Delete(id domain.EvaluationID) error {
	for i, evaluation := range r.evaluations {
		if evaluation.GetID() == id {
			r.evaluations = append(r.evaluations[:i], r.evaluations[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w: %s", domain.ErrEvaluationNotFound, id)
}

func (r *EvaluationRepository) DeleteByStudent(studentID domain.StudentID) {
	r.evaluations = r.filter(func(e *domain.Evaluation) bool {
		return e.GetStudentID() != studentID
	})
}

func (r *EvaluationRepository) DeleteByCriterion(criterionID domain.CriterionID) {
	r.evaluations = r.filter(func(e *domain.Evaluation) bool {
		return e.GetCriterionID() != criterionID
	})
}

func (r *EvaluationRepository) filter(keep func(*domain.Evaluation) bool) []*domain.Evaluation {
	result := make([]*domain.Evaluation, 0, len(r.evaluations))
	for _, evaluation := range r.evaluations {
		if keep(evaluation) {
			result = append(result, evaluation)
		}
	}
	return result
}
