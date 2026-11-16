package application

import "pos-app/internal/domain"

type StudentService struct {
	students    domain.StudentRepository
	evaluations domain.EvaluationRepository
}

func NewStudentService(students domain.StudentRepository, evaluations domain.EvaluationRepository) *StudentService {
	return &StudentService{
		students:    students,
		evaluations: evaluations,
	}
}

func (s *StudentService) Register(fullName, group string) (StudentView, error) {
	student, err := domain.NewStudent(s.students.NextID(), fullName, group)
	if err != nil {
		return StudentView{}, err
	}
	if err := s.students.Save(student); err != nil {
		return StudentView{}, err
	}
	return toStudentView(student), nil
}

func (s *StudentService) Rename(id, fullName string) (StudentView, error) {
	student, err := s.students.FindByID(domain.StudentID(id))
	if err != nil {
		return StudentView{}, err
	}
	if err := student.Rename(fullName); err != nil {
		return StudentView{}, err
	}
	if err := s.students.Save(student); err != nil {
		return StudentView{}, err
	}
	return toStudentView(student), nil
}

func (s *StudentService) Remove(id string) error {
	studentID := domain.StudentID(id)
	if err := s.students.Delete(studentID); err != nil {
		return err
	}

	s.evaluations.DeleteByStudent(studentID)
	return nil
}

func (s *StudentService) List() []StudentView {
	students := s.students.FindAll()
	result := make([]StudentView, 0, len(students))
	for _, student := range students {
		result = append(result, toStudentView(student))
	}
	return result
}
