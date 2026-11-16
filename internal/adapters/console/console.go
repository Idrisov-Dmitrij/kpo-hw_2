package console

import (
	"bufio"
	"fmt"
	"io"

	"pos-app/internal/application"
)

type Console struct {
	facade  *application.PosFacade
	scanner *bufio.Scanner
	out     io.Writer
}

func NewConsole(facade *application.PosFacade, in io.Reader, out io.Writer) *Console {
	return &Console{
		facade:  facade,
		scanner: bufio.NewScanner(in),
		out:     out,
	}
}

func (c *Console) Run() {
	for {
		c.printMenu()
		choice, ok := c.readLine()
		if !ok {
			fmt.Fprintln(c.out, "\nВвод завершён. Выход из программы.")
			return
		}
		fmt.Fprintln(c.out)

		switch choice {
		case "1":
			c.addStudent()
		case "2":
			c.showStudents()
		case "3":
			c.renameStudent()
		case "4":
			c.removeStudent()
		case "5":
			c.addCriterion()
		case "6":
			c.showRubric()
		case "7":
			c.changeCriterionWeight()
		case "8":
			c.removeCriterion()
		case "9":
			c.assess()
		case "10":
			c.showAllEvaluations()
		case "11":
			c.showStudentEvaluations()
		case "12":
			c.correctEvaluation()
		case "13":
			c.cancelEvaluation()
		case "14":
			c.showStudentGrade()
		case "15":
			c.showGradeReport()
		case "16":
			c.selectGradingPolicy()
		case "17":
			c.exportReport()
		case "18":
			c.RunDemo()
		case "0":
			fmt.Fprintln(c.out, "Выход из программы.")
			return
		default:
			fmt.Fprintln(c.out, "Неверный выбор. Попробуйте снова.")
		}
		fmt.Fprintln(c.out)
	}
}

func (c *Console) printMenu() {
	fmt.Fprintln(c.out, "--- ПОС: Преподавательская Оценка Студентов ---")
	fmt.Fprintln(c.out, "Студенты:   1. Добавить   2. Список   3. Переименовать   4. Удалить")
	fmt.Fprintln(c.out, "Рубрикатор: 5. Добавить критерий   6. Показать   7. Изменить вес   8. Удалить критерий")
	fmt.Fprintln(c.out, "Оценки:     9. Выставить   10. Все оценки   11. Оценки студента   12. Исправить   13. Отменить")
	fmt.Fprintln(c.out, "Итог:       14. Итог студента   15. Ведомость   16. Правило расчёта   17. Экспорт")
	fmt.Fprintln(c.out, "Прочее:     18. Демо-сценарий   0. Выход")
	fmt.Fprint(c.out, "Выберите действие: ")
}

func (c *Console) printError(err error) {
	fmt.Fprintf(c.out, "Ошибка: %v\n", err)
}
