package main

import (
	"os"

	"github.com/skvdmt/go-template-app/internal"
	"github.com/skvdmt/go-template-app/internal/model"
)

// main Точка входа в приложение.
func main() {
	// Создание глобального экземпляра журнала.
	var err error
	if model.Log, err = model.NewLogger(); err != nil {
		panic(err)
	}
	// Создание приложения.
	app, err := internal.NewApp()
	if err != nil {
		model.Log.Error.Error(err.Error())
		os.Exit(1)
	}
	// Запуск приложения.
	if errs := app.Start(); len(errs) > 0 {
		for _, err := range errs {
			model.Log.Error.Error(err.Error())
		}
	}
	// Закрытие журнала.
	if err := model.Log.Close(); err != nil {
		panic(err)
	}
}
