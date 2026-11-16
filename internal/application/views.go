package application

import (
	"time"

	"pos-app/internal/domain"
)

type StudentView struct {
	ID       string
	FullName string
	Group    string
}

func toStudentView(student *domain.Student) StudentView {
	return StudentView{
		ID:       string(student.GetID()),
		FullName: student.GetFullName(),
		Group:    student.GetGroup(),
	}
}

//--------------------------------------

type CriterionView struct {
	ID       string
	Name     string
	MaxScore float64
	Weight   float64
}

func toCriterionView(criterion *domain.Criterion) CriterionView {
	return CriterionView{
		ID:       string(criterion.GetID()),
		Name:     criterion.GetName(),
		MaxScore: criterion.GetMaxScore().GetValue(),
		Weight:   criterion.GetWeight().GetValue(),
	}
}

//--------------------------------------

type EvaluationView struct {
	ID            string
	StudentID     string
	StudentName   string
	CriterionID   string
	CriterionName string
	Score         float64
	MaxScore      float64
	Comment       string
	GradedAt      time.Time
}

//--------------------------------------

type ScoreCell struct {
	Value  float64
	Graded bool
}

type StudentGradeView struct {
	StudentID     string
	FullName      string
	Group         string
	Scores        []ScoreCell
	GradedCount   int
	TotalCriteria int
	FinalGrade    float64
}

type GradeReport struct {
	PolicyName    string
	GeneratedAt   time.Time
	Criteria      []CriterionView
	Rows          []StudentGradeView
	CourseAverage float64
}
