package application

import (
	"fmt"

	"pos-app/internal/domain"
)

type GradingService struct {
	students    domain.StudentRepository
	rubrics     domain.RubricRepository
	evaluations domain.EvaluationRepository
	factory     *domain.EvaluationFactory
	clock       Clock
}

func NewGradingService(
	students domain.StudentRepository,
	rubrics domain.RubricRepository,
	evaluations domain.EvaluationRepository,
	factory *domain.EvaluationFactory,
	clock Clock,
) *GradingService {
	return &GradingService{
		students:    students,
		rubrics:     rubrics,
		evaluations: evaluations,
		factory:     factory,
		clock:       clock,
	}
}

func (s *GradingService) Assess(studentID, criterionID string, score float64, comment string) (EvaluationView, error) {
	student, err := s.students.FindByID(domain.StudentID(studentID))
	if err != nil {
		return EvaluationView{}, err
	}

	if existing, ok := s.evaluations.FindByStudentAndCriterion(student.GetID(), domain.CriterionID(criterionID)); ok {
		return EvaluationView{}, fmt.Errorf("%w (оценка %s) — используйте исправление", domain.ErrDuplicateEvaluation, existing.GetID())
	}

	evaluation, err := s.factory.Create(
		s.evaluations.NextID(),
		s.rubrics.Get(),
		student,
		domain.CriterionID(criterionID),
		score,
		comment,
		s.clock(),
	)
	if err != nil {
		return EvaluationView{}, err
	}
	if err := s.evaluations.Save(evaluation); err != nil {
		return EvaluationView{}, err
	}
	return s.toView(evaluation), nil
}

func (s *GradingService) Correct(evaluationID string, score float64, comment string) (EvaluationView, error) {
	evaluation, err := s.evaluations.FindByID(domain.EvaluationID(evaluationID))
	if err != nil {
		return EvaluationView{}, err
	}

	criterion, err := s.rubrics.Get().GetCriterion(evaluation.GetCriterionID())
	if err != nil {
		return EvaluationView{}, err
	}
	newScore, err := domain.NewScore(score)
	if err != nil {
		return EvaluationView{}, err
	}
	if err := evaluation.Correct(criterion, newScore, comment, s.clock()); err != nil {
		return EvaluationView{}, err
	}
	if err := s.evaluations.Save(evaluation); err != nil {
		return EvaluationView{}, err
	}
	return s.toView(evaluation), nil
}

func (s *GradingService) Cancel(evaluationID string) error {
	return s.evaluations.Delete(domain.EvaluationID(evaluationID))
}

func (s *GradingService) ListByStudent(studentID string) ([]EvaluationView, error) {
	student, err := s.students.FindByID(domain.StudentID(studentID))
	if err != nil {
		return nil, err
	}
	return s.toViews(s.evaluations.FindByStudent(student.GetID())), nil
}

func (s *GradingService) ListAll() []EvaluationView {
	return s.toViews(s.evaluations.FindAll())
}

func (s *GradingService) toViews(evaluations []*domain.Evaluation) []EvaluationView {
	result := make([]EvaluationView, 0, len(evaluations))
	for _, evaluation := range evaluations {
		result = append(result, s.toView(evaluation))
	}
	return result
}

func (s *GradingService) toView(evaluation *domain.Evaluation) EvaluationView {
	view := EvaluationView{
		ID:          string(evaluation.GetID()),
		StudentID:   string(evaluation.GetStudentID()),
		CriterionID: string(evaluation.GetCriterionID()),
		Score:       evaluation.GetScore().GetValue(),
		Comment:     evaluation.GetComment(),
		GradedAt:    evaluation.GetGradedAt(),
	}

	if student, err := s.students.FindByID(evaluation.GetStudentID()); err == nil {
		view.StudentName = student.GetFullName()
	}
	if criterion, err := s.rubrics.Get().GetCriterion(evaluation.GetCriterionID()); err == nil {
		view.CriterionName = criterion.GetName()
		view.MaxScore = criterion.GetMaxScore().GetValue()
	}
	return view
}
