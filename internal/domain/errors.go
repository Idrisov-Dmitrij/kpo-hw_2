package domain

import "errors"

var (
	ErrEmptyID                  = errors.New("идентификатор не может быть пустым")
	ErrEmptyStudentName         = errors.New("ФИО студента не может быть пустым")
	ErrEmptyCriterionName       = errors.New("название критерия не может быть пустым")
	ErrDuplicateCriterionName   = errors.New("критерий с таким названием уже есть в рубрикаторе")
	ErrInvalidScore             = errors.New("балл должен быть неотрицательным числом")
	ErrInvalidMaxScore          = errors.New("максимальный балл критерия должен быть больше нуля")
	ErrInvalidWeight            = errors.New("вес критерия должен быть больше нуля")
	ErrInvalidFinalGrade        = errors.New("итог должен лежать в диапазоне [0, 10]")
	ErrScoreOutOfRange          = errors.New("балл выходит за пределы критерия")
	ErrCriterionNotInRubric     = errors.New("критерия нет в рубрикаторе")
	ErrCriterionMismatch        = errors.New("оценка относится к другому критерию")
	ErrStudentRequired          = errors.New("оценку можно выставить только существующему студенту")
	ErrRubricRequired           = errors.New("оценку можно выставить только по рубрикатору")
	ErrForeignEvaluation        = errors.New("оценка принадлежит другому студенту")
	ErrDuplicateEvaluation      = errors.New("по критерию уже выставлена оценка этому студенту")
	ErrStudentNotFound          = errors.New("студент не найден")
	ErrEvaluationNotFound       = errors.New("оценка не найдена")
	ErrGradingPolicyRequired    = errors.New("не задано правило расчёта итога")
	ErrInvalidRoundingPrecision = errors.New("точность округления должна быть в диапазоне [0, 4]")
)
