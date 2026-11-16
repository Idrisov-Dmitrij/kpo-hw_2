package memory

import (
	"errors"
	"testing"
	"time"

	"pos-app/internal/domain"
)

var _ domain.StudentRepository = (*StudentRepository)(nil)
var _ domain.RubricRepository = (*RubricRepository)(nil)
var _ domain.EvaluationRepository = (*EvaluationRepository)(nil)

func TestStudentRepository_NextIDIsReservedOnlyBySave(t *testing.T) {
	repo := NewStudentRepository()

	if repo.NextID() != repo.NextID() {
		t.Error("NextID без сохранения не должен сдвигать последовательность")
	}

	student, _ := domain.NewStudent(repo.NextID(), "Иванов", "")
	repo.Save(student)

	if repo.NextID() != "S-2" {
		t.Errorf("Ожидался следующий ID S-2, получено %s", repo.NextID())
	}
}

func TestStudentRepository_CRUD(t *testing.T) {
	repo := NewStudentRepository()
	student, _ := domain.NewStudent(repo.NextID(), "Иванов", "")
	repo.Save(student)
	repo.Save(student)

	if len(repo.FindAll()) != 1 {
		t.Errorf("Повторное сохранение не должно дублировать студента, получено %d", len(repo.FindAll()))
	}
	if _, err := repo.FindByID("S-1"); err != nil {
		t.Errorf("Не ожидалась ошибка, получено: %v", err)
	}
	if err := repo.Delete("S-1"); err != nil {
		t.Errorf("Не ожидалась ошибка, получено: %v", err)
	}
	if _, err := repo.FindByID("S-1"); !errors.Is(err, domain.ErrStudentNotFound) {
		t.Errorf("Ожидалась ошибка ErrStudentNotFound, получено %v", err)
	}
}

func TestEvaluationRepository_FiltersAndCascades(t *testing.T) {
	repo := NewEvaluationRepository()
	rubric := domain.NewRubric()
	max, _ := domain.NewMaxScore(10)
	weight, _ := domain.NewWeight(1)
	homework, _ := rubric.AddCriterion("ДЗ", max, weight)
	exam, _ := rubric.AddCriterion("Экзамен", max, weight)
	ivanov, _ := domain.NewStudent("S-1", "Иванов", "")
	petrova, _ := domain.NewStudent("S-2", "Петрова", "")
	factory := domain.NewEvaluationFactory()

	for _, pair := range []struct {
		student   *domain.Student
		criterion domain.CriterionID
	}{
		{ivanov, homework.GetID()},
		{ivanov, exam.GetID()},
		{petrova, homework.GetID()},
	} {
		evaluation, _ := factory.Create(repo.NextID(), rubric, pair.student, pair.criterion, 5, "", time.Now())
		repo.Save(evaluation)
	}

	if len(repo.FindByStudent("S-1")) != 2 {
		t.Errorf("Ожидалось 2 оценки Иванова, получено %d", len(repo.FindByStudent("S-1")))
	}
	if _, ok := repo.FindByStudentAndCriterion("S-2", exam.GetID()); ok {
		t.Error("У Петровой нет оценки за экзамен")
	}

	repo.DeleteByCriterion(homework.GetID())
	if len(repo.FindAll()) != 1 {
		t.Errorf("После удаления по критерию ожидалась 1 оценка, получено %d", len(repo.FindAll()))
	}

	repo.DeleteByStudent("S-1")
	if len(repo.FindAll()) != 0 {
		t.Errorf("После удаления по студенту оценок быть не должно, получено %d", len(repo.FindAll()))
	}
}
