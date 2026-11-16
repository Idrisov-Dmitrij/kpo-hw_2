package memory

import (
	"fmt"

	"pos-app/internal/domain"
)

type StudentRepository struct {
	students []*domain.Student
	sequence int
}

func NewStudentRepository() *StudentRepository {
	return &StudentRepository{
		students: make([]*domain.Student, 0),
	}
}

func (r *StudentRepository) NextID() domain.StudentID {
	return domain.StudentID(fmt.Sprintf("S-%d", r.sequence+1))
}

func (r *StudentRepository) Save(student *domain.Student) error {
	for i, existing := range r.students {
		if existing.GetID() == student.GetID() {
			r.students[i] = student
			return nil
		}
	}

	r.students = append(r.students, student)
	r.sequence++
	return nil
}

func (r *StudentRepository) FindByID(id domain.StudentID) (*domain.Student, error) {
	for _, student := range r.students {
		if student.GetID() == id {
			return student, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", domain.ErrStudentNotFound, id)
}

func (r *StudentRepository) FindAll() []*domain.Student {
	result := make([]*domain.Student, len(r.students))
	copy(result, r.students)
	return result
}

func (r *StudentRepository) Delete(id domain.StudentID) error {
	for i, student := range r.students {
		if student.GetID() == id {
			r.students = append(r.students[:i], r.students[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w: %s", domain.ErrStudentNotFound, id)
}
