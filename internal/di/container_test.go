package di

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestContainer_ReturnsSingletons(t *testing.T) {
	container := NewContainer()

	first, err := container.BuildFacade()
	if err != nil {
		t.Fatalf("Не ожидалась ошибка сборки, получено: %v", err)
	}
	second, _ := container.BuildFacade()

	if first != second {
		t.Error("Контейнер должен возвращать один и тот же фасад")
	}
	if container.StudentRepository() != container.StudentRepository() {
		t.Error("Контейнер должен возвращать один и тот же репозиторий")
	}
}

func TestContainer_ConsoleScriptedSession(t *testing.T) {
	input := strings.Join([]string{
		"1", "Иванов Иван", "БПИ-241",
		"5", "ДЗ", "10", "1",
		"9", "s-1", "c-1", "12", "",
		"9", "s-1", "c-1", "9,5", "хорошо",
		"14", "S-1",
		"0",
	}, "\n")
	var out bytes.Buffer

	app, err := NewContainer().BuildConsole(strings.NewReader(input), &out)
	if err != nil {
		t.Fatalf("Не ожидалась ошибка сборки, получено: %v", err)
	}
	app.Run()

	output := out.String()
	for _, expected := range []string{
		"добавлен с ID S-1",
		"добавлен с ID C-1",
		"балл выходит за пределы критерия",
		"Оценка E-1 выставлена",
		"итог 9.50 из 10",
		"Выход из программы.",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("В выводе консоли нет %q", expected)
		}
	}
}

func TestContainer_ConsoleStopsOnEOF(t *testing.T) {
	var out bytes.Buffer
	app, _ := NewContainer().BuildConsole(strings.NewReader("5\nДЗ\nабв"), &out)

	app.Run()

	if !strings.Contains(out.String(), "Ввод завершён") {
		t.Error("При окончании ввода консоль должна корректно завершаться")
	}
}

func TestContainer_DemoScenario(t *testing.T) {
	workDir, _ := os.Getwd()
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("Не удалось сменить рабочую папку: %v", err)
	}
	defer os.Chdir(workDir)

	var out bytes.Buffer
	app, _ := NewContainer().BuildConsole(strings.NewReader(""), &out)
	app.RunDemo()

	output := out.String()
	for _, expected := range []string{
		"12.00 > 10.00",
		"критерия нет в рубрикаторе",
		"[исправлено] E-6",
		"[отменено] E-8",
		"итог выгружен в pos_demo_report.csv",
		"итог выгружен в pos_demo_report.json",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("В выводе демо нет %q", expected)
		}
	}
	for _, file := range []string{"pos_demo_report.csv", "pos_demo_report.json"} {
		if _, err := os.Stat(file); err != nil {
			t.Errorf("Файл экспорта %s не создан", file)
		}
	}

	out.Reset()
	app.RunDemo()
	if !strings.Contains(out.String(), "только на пустом модуле") {
		t.Error("Повторный запуск демо на заполненном модуле должен быть отклонён")
	}
}
