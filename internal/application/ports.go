package application

import (
	"errors"
	"time"
)

var (
	ErrUnknownExportFormat = errors.New("неизвестный формат экспорта")
	ErrUnknownPolicy       = errors.New("неизвестное правило расчёта итога")
	ErrEmptyExportPath     = errors.New("путь для экспорта не может быть пустым")
)

type ReportExporter interface {
	GetFormat() string
	Export(report GradeReport, path string) error
}

type Clock func() time.Time
