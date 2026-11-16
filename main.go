package main

import (
	"flag"
	"fmt"
	"os"

	"pos-app/internal/di"
)

func main() {
	demo := flag.Bool("demo", false, "запустить сквозной демо-сценарий и выйти")
	flag.Parse()

	container := di.NewContainer()
	app, err := container.BuildConsole(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка запуска: %v\n", err)
		os.Exit(1)
	}

	if *demo {
		app.RunDemo()
		return
	}
	app.Run()
}
