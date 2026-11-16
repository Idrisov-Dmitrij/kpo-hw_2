package memory

import "pos-app/internal/domain"

type RubricRepository struct {
	rubric *domain.Rubric
}

func NewRubricRepository() *RubricRepository {
	return &RubricRepository{
		rubric: domain.NewRubric(),
	}
}

func (r *RubricRepository) Get() *domain.Rubric {
	return r.rubric
}

func (r *RubricRepository) Save(rubric *domain.Rubric) error {
	r.rubric = rubric
	return nil
}
