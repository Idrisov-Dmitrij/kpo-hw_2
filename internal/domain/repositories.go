package domain

type StudentRepository interface {
	NextID() StudentID
	Save(student *Student) error
	FindByID(id StudentID) (*Student, error)
	FindAll() []*Student
	Delete(id StudentID) error
}

type RubricRepository interface {
	Get() *Rubric
	Save(rubric *Rubric) error
}

type EvaluationRepository interface {
	NextID() EvaluationID
	Save(evaluation *Evaluation) error
	FindByID(id EvaluationID) (*Evaluation, error)
	FindByStudent(studentID StudentID) []*Evaluation
	FindByStudentAndCriterion(studentID StudentID, criterionID CriterionID) (*Evaluation, bool)
	FindAll() []*Evaluation
	Delete(id EvaluationID) error
	DeleteByStudent(studentID StudentID)
	DeleteByCriterion(criterionID CriterionID)
}
