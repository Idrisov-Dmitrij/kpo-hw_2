package di

import (
	"io"
	"time"

	"pos-app/internal/adapters/console"
	"pos-app/internal/adapters/export"
	"pos-app/internal/adapters/memory"
	"pos-app/internal/application"
	"pos-app/internal/domain"
)

type Container struct {
	students    domain.StudentRepository
	rubrics     domain.RubricRepository
	evaluations domain.EvaluationRepository
	calculator  *domain.GradeCalculator
	policies    []domain.GradingPolicy
	facade      *application.PosFacade
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Clock() application.Clock {
	return time.Now
}

func (c *Container) StudentRepository() domain.StudentRepository {
	if c.students == nil {
		c.students = memory.NewStudentRepository()
	}
	return c.students
}

func (c *Container) RubricRepository() domain.RubricRepository {
	if c.rubrics == nil {
		c.rubrics = memory.NewRubricRepository()
	}
	return c.rubrics
}

func (c *Container) EvaluationRepository() domain.EvaluationRepository {
	if c.evaluations == nil {
		c.evaluations = memory.NewEvaluationRepository()
	}
	return c.evaluations
}

func (c *Container) Exporters() []application.ReportExporter {
	return []application.ReportExporter{
		export.NewCSVExporter(),
		export.NewJSONExporter(),
	}
}

func (c *Container) GradingPolicies() ([]domain.GradingPolicy, error) {
	if c.policies != nil {
		return c.policies, nil
	}

	weightedHundredths, err := domain.NewRoundedPolicy(domain.NewWeightedAveragePolicy(), 2)
	if err != nil {
		return nil, err
	}
	weightedInteger, err := domain.NewRoundedPolicy(domain.NewWeightedAveragePolicy(), 0)
	if err != nil {
		return nil, err
	}
	progressHundredths, err := domain.NewRoundedPolicy(domain.NewCurrentProgressPolicy(), 2)
	if err != nil {
		return nil, err
	}

	c.policies = []domain.GradingPolicy{weightedHundredths, weightedInteger, progressHundredths}
	return c.policies, nil
}

func (c *Container) GradeCalculator() (*domain.GradeCalculator, error) {
	if c.calculator != nil {
		return c.calculator, nil
	}

	policies, err := c.GradingPolicies()
	if err != nil {
		return nil, err
	}
	calculator, err := domain.NewGradeCalculator(policies[0])
	if err != nil {
		return nil, err
	}

	c.calculator = calculator
	return c.calculator, nil
}

func (c *Container) BuildFacade() (*application.PosFacade, error) {
	if c.facade != nil {
		return c.facade, nil
	}

	calculator, err := c.GradeCalculator()
	if err != nil {
		return nil, err
	}
	policies, err := c.GradingPolicies()
	if err != nil {
		return nil, err
	}

	students := application.NewStudentService(c.StudentRepository(), c.EvaluationRepository())
	rubric := application.NewRubricService(c.RubricRepository(), c.EvaluationRepository())
	grading := application.NewGradingService(
		c.StudentRepository(),
		c.RubricRepository(),
		c.EvaluationRepository(),
		domain.NewEvaluationFactory(),
		c.Clock(),
	)
	reports := application.NewReportService(
		c.StudentRepository(),
		c.RubricRepository(),
		c.EvaluationRepository(),
		calculator,
		policies,
		c.Exporters(),
		c.Clock(),
	)

	c.facade = application.NewPosFacade(students, rubric, grading, reports)
	return c.facade, nil
}

func (c *Container) BuildConsole(in io.Reader, out io.Writer) (*console.Console, error) {
	facade, err := c.BuildFacade()
	if err != nil {
		return nil, err
	}
	return console.NewConsole(facade, in, out), nil
}
