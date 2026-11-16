package application

import "pos-app/internal/domain"

type RubricService struct {
	rubrics     domain.RubricRepository
	evaluations domain.EvaluationRepository
}

func NewRubricService(rubrics domain.RubricRepository, evaluations domain.EvaluationRepository) *RubricService {
	return &RubricService{
		rubrics:     rubrics,
		evaluations: evaluations,
	}
}

func (s *RubricService) AddCriterion(name string, maxScore, weight float64) (CriterionView, error) {
	max, err := domain.NewMaxScore(maxScore)
	if err != nil {
		return CriterionView{}, err
	}
	w, err := domain.NewWeight(weight)
	if err != nil {
		return CriterionView{}, err
	}

	rubric := s.rubrics.Get()
	criterion, err := rubric.AddCriterion(name, max, w)
	if err != nil {
		return CriterionView{}, err
	}
	if err := s.rubrics.Save(rubric); err != nil {
		return CriterionView{}, err
	}
	return toCriterionView(criterion), nil
}

func (s *RubricService) ChangeWeight(id string, weight float64) (CriterionView, error) {
	w, err := domain.NewWeight(weight)
	if err != nil {
		return CriterionView{}, err
	}

	rubric := s.rubrics.Get()
	criterionID := domain.CriterionID(id)
	if err := rubric.ChangeWeight(criterionID, w); err != nil {
		return CriterionView{}, err
	}
	if err := s.rubrics.Save(rubric); err != nil {
		return CriterionView{}, err
	}

	criterion, err := rubric.GetCriterion(criterionID)
	if err != nil {
		return CriterionView{}, err
	}
	return toCriterionView(criterion), nil
}

func (s *RubricService) RemoveCriterion(id string) error {
	rubric := s.rubrics.Get()
	criterionID := domain.CriterionID(id)
	if err := rubric.RemoveCriterion(criterionID); err != nil {
		return err
	}
	if err := s.rubrics.Save(rubric); err != nil {
		return err
	}

	s.evaluations.DeleteByCriterion(criterionID)
	return nil
}

func (s *RubricService) List() []CriterionView {
	criteria := s.rubrics.Get().GetCriteria()
	result := make([]CriterionView, 0, len(criteria))
	for _, criterion := range criteria {
		result = append(result, toCriterionView(criterion))
	}
	return result
}
