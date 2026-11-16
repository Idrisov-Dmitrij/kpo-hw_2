package console

import (
	"fmt"
	"strings"

	"pos-app/internal/application"
)

func (c *Console) addStudent() {
	name, ok := c.readString("ФИО студента")
	if !ok {
		return
	}
	group, ok := c.readString("Группа (можно пропустить)")
	if !ok {
		return
	}

	student, err := c.facade.RegisterStudent(name, group)
	if err != nil {
		c.printError(err)
		return
	}
	fmt.Fprintf(c.out, "Студент %s добавлен с ID %s.\n", student.FullName, student.ID)
}

func (c *Console) showStudents() {
	students := c.facade.ListStudents()
	if len(students) == 0 {
		fmt.Fprintln(c.out, "Студентов пока нет.")
		return
	}

	fmt.Fprintf(c.out, "%-6s %-30s %-12s\n", "ID", "ФИО", "Группа")
	fmt.Fprintln(c.out, strings.Repeat("-", 50))
	for _, student := range students {
		fmt.Fprintf(c.out, "%-6s %-30s %-12s\n", student.ID, student.FullName, student.Group)
	}
}

func (c *Console) renameStudent() {
	id, ok := c.readID("ID студента")
	if !ok {
		return
	}
	name, ok := c.readString("Новое ФИО")
	if !ok {
		return
	}

	student, err := c.facade.RenameStudent(id, name)
	if err != nil {
		c.printError(err)
		return
	}
	fmt.Fprintf(c.out, "Студент %s теперь: %s.\n", student.ID, student.FullName)
}

func (c *Console) removeStudent() {
	id, ok := c.readID("ID студента")
	if !ok {
		return
	}

	if err := c.facade.RemoveStudent(id); err != nil {
		c.printError(err)
		return
	}
	fmt.Fprintf(c.out, "Студент %s и его оценки удалены.\n", id)
}

//--------------------------------------

func (c *Console) addCriterion() {
	name, ok := c.readString("Название критерия")
	if !ok {
		return
	}
	maxScore, ok := c.readFloat("Максимальный балл")
	if !ok {
		return
	}
	weight, ok := c.readFloat("Вес в итоге")
	if !ok {
		return
	}

	criterion, err := c.facade.AddCriterion(name, maxScore, weight)
	if err != nil {
		c.printError(err)
		return
	}
	fmt.Fprintf(c.out, "Критерий «%s» добавлен с ID %s.\n", criterion.Name, criterion.ID)
}

func (c *Console) showRubric() {
	criteria := c.facade.ListCriteria()
	if len(criteria) == 0 {
		fmt.Fprintln(c.out, "Рубрикатор пуст.")
		return
	}

	var totalWeight float64
	for _, criterion := range criteria {
		totalWeight += criterion.Weight
	}

	fmt.Fprintf(c.out, "%-6s %-28s %-10s %-8s %-8s\n", "ID", "Критерий", "Макс.", "Вес", "Доля")
	fmt.Fprintln(c.out, strings.Repeat("-", 64))
	for _, criterion := range criteria {
		fmt.Fprintf(c.out, "%-6s %-28s %-10g %-8g %5.1f%%\n",
			criterion.ID, criterion.Name, criterion.MaxScore, criterion.Weight, criterion.Weight/totalWeight*100)
	}
}

func (c *Console) changeCriterionWeight() {
	id, ok := c.readID("ID критерия")
	if !ok {
		return
	}
	weight, ok := c.readFloat("Новый вес")
	if !ok {
		return
	}

	criterion, err := c.facade.ChangeCriterionWeight(id, weight)
	if err != nil {
		c.printError(err)
		return
	}
	fmt.Fprintf(c.out, "Вес критерия «%s» теперь %g.\n", criterion.Name, criterion.Weight)
}

func (c *Console) removeCriterion() {
	id, ok := c.readID("ID критерия")
	if !ok {
		return
	}

	if err := c.facade.RemoveCriterion(id); err != nil {
		c.printError(err)
		return
	}
	fmt.Fprintf(c.out, "Критерий %s и оценки по нему удалены.\n", id)
}

//--------------------------------------

func (c *Console) assess() {
	studentID, ok := c.readID("ID студента")
	if !ok {
		return
	}
	criterionID, ok := c.readID("ID критерия")
	if !ok {
		return
	}
	score, ok := c.readFloat("Балл")
	if !ok {
		return
	}
	comment, ok := c.readString("Комментарий (можно пропустить)")
	if !ok {
		return
	}

	evaluation, err := c.facade.Assess(studentID, criterionID, score, comment)
	if err != nil {
		c.printError(err)
		return
	}
	fmt.Fprintf(c.out, "Оценка %s выставлена: %s — «%s»: %g/%g.\n",
		evaluation.ID, evaluation.StudentName, evaluation.CriterionName, evaluation.Score, evaluation.MaxScore)
}

func (c *Console) showAllEvaluations() {
	c.printEvaluations(c.facade.ListEvaluations())
}

func (c *Console) showStudentEvaluations() {
	studentID, ok := c.readID("ID студента")
	if !ok {
		return
	}

	evaluations, err := c.facade.ListStudentEvaluations(studentID)
	if err != nil {
		c.printError(err)
		return
	}
	c.printEvaluations(evaluations)
}

func (c *Console) printEvaluations(evaluations []application.EvaluationView) {
	if len(evaluations) == 0 {
		fmt.Fprintln(c.out, "Оценок пока нет.")
		return
	}

	fmt.Fprintf(c.out, "%-6s %-22s %-24s %-10s %-17s %s\n", "ID", "Студент", "Критерий", "Балл", "Дата", "Комментарий")
	fmt.Fprintln(c.out, strings.Repeat("-", 96))
	for _, e := range evaluations {
		fmt.Fprintf(c.out, "%-6s %-22s %-24s %-10s %-17s %s\n",
			e.ID, e.StudentName, e.CriterionName,
			fmt.Sprintf("%g/%g", e.Score, e.MaxScore),
			e.GradedAt.Format("02.01.2006 15:04"), e.Comment)
	}
}

func (c *Console) correctEvaluation() {
	id, ok := c.readID("ID оценки")
	if !ok {
		return
	}
	score, ok := c.readFloat("Новый балл")
	if !ok {
		return
	}
	comment, ok := c.readString("Комментарий (можно пропустить)")
	if !ok {
		return
	}

	evaluation, err := c.facade.CorrectEvaluation(id, score, comment)
	if err != nil {
		c.printError(err)
		return
	}
	fmt.Fprintf(c.out, "Оценка %s исправлена: %g/%g.\n", evaluation.ID, evaluation.Score, evaluation.MaxScore)
}

func (c *Console) cancelEvaluation() {
	id, ok := c.readID("ID оценки")
	if !ok {
		return
	}

	if err := c.facade.CancelEvaluation(id); err != nil {
		c.printError(err)
		return
	}
	fmt.Fprintf(c.out, "Оценка %s отменена.\n", id)
}

//--------------------------------------

func (c *Console) showStudentGrade() {
	studentID, ok := c.readID("ID студента")
	if !ok {
		return
	}

	grade, err := c.facade.StudentGrade(studentID)
	if err != nil {
		c.printError(err)
		return
	}
	fmt.Fprintf(c.out, "Правило: %s\n", c.facade.CurrentGradingPolicy())
	fmt.Fprintf(c.out, "%s (%s): итог %.2f из 10, оценено критериев %d/%d.\n",
		grade.FullName, grade.StudentID, grade.FinalGrade, grade.GradedCount, grade.TotalCriteria)
}

func (c *Console) showGradeReport() {
	report, err := c.facade.GradeReport()
	if err != nil {
		c.printError(err)
		return
	}
	if len(report.Rows) == 0 {
		fmt.Fprintln(c.out, "Студентов пока нет — ведомость пуста.")
		return
	}

	fmt.Fprintf(c.out, "Правило: %s\n", report.PolicyName)
	fmt.Fprintf(c.out, "%-6s %-30s %-12s %-9s %-6s\n", "ID", "ФИО", "Группа", "Оценено", "Итог")
	fmt.Fprintln(c.out, strings.Repeat("-", 67))
	for _, row := range report.Rows {
		fmt.Fprintf(c.out, "%-6s %-30s %-12s %-9s %6.2f\n",
			row.StudentID, row.FullName, row.Group,
			fmt.Sprintf("%d/%d", row.GradedCount, row.TotalCriteria), row.FinalGrade)
	}
	fmt.Fprintln(c.out, strings.Repeat("-", 67))
	fmt.Fprintf(c.out, "Средний итог по курсу: %.2f\n", report.CourseAverage)
}

func (c *Console) selectGradingPolicy() {
	current := c.facade.CurrentGradingPolicy()
	for i, name := range c.facade.GradingPolicies() {
		marker := " "
		if name == current {
			marker = "*"
		}
		fmt.Fprintf(c.out, "%s %d. %s\n", marker, i+1, name)
	}

	choice, ok := c.readInt("Номер правила")
	if !ok {
		return
	}
	if err := c.facade.SelectGradingPolicy(choice - 1); err != nil {
		c.printError(err)
		return
	}
	fmt.Fprintf(c.out, "Текущее правило: %s\n", c.facade.CurrentGradingPolicy())
}

func (c *Console) exportReport() {
	formats := strings.Join(c.facade.ExportFormats(), "/")
	format, ok := c.readString(fmt.Sprintf("Формат (%s)", formats))
	if !ok {
		return
	}
	path, ok := c.readString("Путь к файлу")
	if !ok {
		return
	}

	if err := c.facade.ExportReport(format, path); err != nil {
		c.printError(err)
		return
	}
	fmt.Fprintf(c.out, "Итог выгружен в %s.\n", path)
}
