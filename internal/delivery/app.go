package delivery

import (
	"context"

	"github.com/skvdmt/go-template-app/internal/model"
	"github.com/skvdmt/go-template-app/internal/usecase"
)

// App Транспортный слой.
type App struct {
	// Сервисный слой.
	usecase Usecase
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Logs.Info.Info("delivery layer creating")
	a := &App{}
	var err error
	// Создание транспортного слоя.
	a.usecase, err = usecase.NewApp()
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Start Запуск.
func (a *App) Start(ctx context.Context) error {
	model.Logs.Info.Info("delivery layer starting")
	return nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	if err := a.usecase.Stop(ctx); err != nil {
		return err
	}
	model.Logs.Info.Info("delivery layer stopped")
	return nil
}
