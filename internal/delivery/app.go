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
	model.Log.Info.Info("delivery layer creating")
	// Создание сервисного слоя.
	u, err := usecase.NewApp()
	if err != nil {
		return nil, err
	}
	return &App{
		usecase: u,
	}, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	if err := a.usecase.Stop(ctx); err != nil {
		return err
	}
	model.Log.Info.Info("delivery layer stopped")
	return nil
}
