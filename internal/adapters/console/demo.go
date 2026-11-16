package console

import (
	"fmt"

	"pos-app/internal/application"
)

const (
	demoCSVPath  = "pos_demo_report.csv"
	demoJSONPath = "pos_demo_report.json"
)

func (c *Console) RunDemo() {
	if !c.facade.IsEmpty() {
		fmt.Fprintln(c.out, "Демо-сценарий запускается только на пустом модуле (без студентов и критериев).")
		return
	}

	fmt.Fprintln(c.out, "--- Демо-сценарий ПОС ---")

	fmt.Fprintln(c.out, "\n1. Заводим студентов:")
	ivanov := c.demoStudent("Иванов Иван", "БПИ-241")
	petrova := c.demoStudent("Петрова Анна", "БПИ-241")
	sidorov := c.demoStudent("Сидоров Пётр", "БПИ-242")
	c.demoStudent("   ", "БПИ-242")
	c.showStudents()

	fmt.Fprintln(c.out, "\n2. Собираем рубрикатор:")
	homework := c.demoCriterion("Домашние задания", 10, 0.4)
	test := c.demoCriterion("Контрольная работа", 10, 0.2)
	exam := c.demoCriterion("Экзамен", 100, 0.4)
	c.demoCriterion("Посещаемость", 0, 0.1)
	c.demoCriterion("Активность", 10, -1)
	c.demoCriterion("экзамен", 100, 0.4)
	c.showRubric()

	fmt.Fprintln(c.out, "\n3. Выставляем оценки:")
	c.demoAssess(ivanov, homework, 9, "Все ДЗ в срок")
	c.demoAssess(ivanov, test, 8, "")
	c.demoAssess(ivanov, exam, 85, "Уверенный ответ")
	c.demoAssess(petrova, homework, 10, "")
	c.demoAssess(petrova, test, 6, "Ошибка в задаче 3")
	typo := c.demoAssess(petrova, exam, 59, "Опечатка преподавателя")
	c.demoAssess(sidorov, homework, 7, "Одно ДЗ просрочено")

	fmt.Fprintln(c.out, "\n4. Пытаемся нарушить правила домена (ожидаются ошибки):")
	c.demoAssess(sidorov, test, 12, "12 баллов из 10")
	c.demoAssess(sidorov, test, -1, "Отрицательный балл")
	c.demoAssess(sidorov, "C-99", 5, "Критерия нет в рубрикаторе")
	c.demoAssess("S-99", test, 5, "Студента не существует")
	c.demoAssess(ivanov, homework, 10, "Повторная оценка по тому же критерию")

	fmt.Fprintln(c.out, "\n5. Исправляем и отменяем оценки:")
	c.demoCorrect(typo, 69, "Пересчитано после показа работы")
	c.demoCorrect(typo, 120, "Балл больше максимума")
	extra := c.demoAssess(sidorov, test, 5, "Поставлено по ошибке")
	c.demoCancel(extra)
	c.showAllEvaluations()

	fmt.Fprintln(c.out, "\n6. Получаем итог:")
	c.showGradeReport()

	fmt.Fprintln(c.out, "\n7. Меняем правило расчёта итога (Strategy) и округление (Decorator):")
	for i := 1; i < len(c.facade.GradingPolicies()); i++ {
		if err := c.facade.SelectGradingPolicy(i); err != nil {
			c.printError(err)
			continue
		}
		fmt.Fprintln(c.out)
		c.showGradeReport()
	}
	if err := c.facade.SelectGradingPolicy(0); err != nil {
		c.printError(err)
	}
	fmt.Fprintf(c.out, "\nВозвращаем правило по умолчанию: %s\n", c.facade.CurrentGradingPolicy())

	fmt.Fprintln(c.out, "\n8. Экспорт итога:")
	c.demoExport("csv", demoCSVPath)
	c.demoExport("json", demoJSONPath)

	fmt.Fprintln(c.out, "\nДемо-сценарий завершён. Данные остались в модуле — их можно смотреть через меню.")
}

func (c *Console) demoStudent(name, group string) string {
	student, err := c.facade.RegisterStudent(name, group)
	if err != nil {
		fmt.Fprintf(c.out, "  [отказ] студент «%s»: %v\n", name, err)
		return ""
	}
	fmt.Fprintf(c.out, "  [ок] %s — %s (%s)\n", student.ID, student.FullName, student.Group)
	return student.ID
}

func (c *Console) demoCriterion(name string, maxScore, weight float64) string {
	criterion, err := c.facade.AddCriterion(name, maxScore, weight)
	if err != nil {
		fmt.Fprintf(c.out, "  [отказ] критерий «%s» (max %g, вес %g): %v\n", name, maxScore, weight, err)
		return ""
	}
	fmt.Fprintf(c.out, "  [ок] %s — %s (max %g, вес %g)\n", criterion.ID, criterion.Name, criterion.MaxScore, criterion.Weight)
	return criterion.ID
}

func (c *Console) demoAssess(studentID, criterionID string, score float64, comment string) string {
	evaluation, err := c.facade.Assess(studentID, criterionID, score, comment)
	if err != nil {
		fmt.Fprintf(c.out, "  [отказ] %s / %s = %g: %v\n", studentID, criterionID, score, err)
		return ""
	}
	c.printDemoEvaluation("ок", evaluation)
	return evaluation.ID
}

func (c *Console) demoCorrect(evaluationID string, score float64, comment string) {
	evaluation, err := c.facade.CorrectEvaluation(evaluationID, score, comment)
	if err != nil {
		fmt.Fprintf(c.out, "  [отказ] исправление %s на %g: %v\n", evaluationID, score, err)
		return
	}
	c.printDemoEvaluation("исправлено", evaluation)
}

func (c *Console) demoCancel(evaluationID string) {
	if err := c.facade.CancelEvaluation(evaluationID); err != nil {
		fmt.Fprintf(c.out, "  [отказ] отмена %s: %v\n", evaluationID, err)
		return
	}
	fmt.Fprintf(c.out, "  [отменено] %s\n", evaluationID)
}

func (c *Console) demoExport(format, path string) {
	if err := c.facade.ExportReport(format, path); err != nil {
		fmt.Fprintf(c.out, "  [отказ] экспорт %s: %v\n", format, err)
		return
	}
	fmt.Fprintf(c.out, "  [ок] итог выгружен в %s\n", path)
}

func (c *Console) printDemoEvaluation(status string, evaluation application.EvaluationView) {
	fmt.Fprintf(c.out, "  [%s] %s: %s — «%s» %g/%g\n",
		status, evaluation.ID, evaluation.StudentName, evaluation.CriterionName, evaluation.Score, evaluation.MaxScore)
}
