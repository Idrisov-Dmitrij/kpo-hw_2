package domain

import "strings"

type Student struct {
	id       StudentID
	fullName string
	group    string
}

func NewStudent(id StudentID, fullName, group string) (*Student, error) {
	if strings.TrimSpace(string(id)) == "" {
		return nil, ErrEmptyID
	}

	student := &Student{id: id, group: strings.TrimSpace(group)}
	if err := student.Rename(fullName); err != nil {
		return nil, err
	}

	return student, nil
}

func (s *Student) GetID() StudentID    { return s.id }
func (s *Student) GetFullName() string { return s.fullName }
func (s *Student) GetGroup() string    { return s.group }

func (s *Student) Rename(fullName string) error {
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return ErrEmptyStudentName
	}

	s.fullName = fullName
	return nil
}
