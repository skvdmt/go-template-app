package usecase

import (
	"context"

	"github.com/skvdmt/go-template-app/internal/model"
	"github.com/skvdmt/go-template-app/internal/repository"
)

// App Сервисный слой.
type App struct {
	// Репозиторный слой.
	repository Repository
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Log.Info.Info("usecase layer creating")
	// Создание репозиторного слоя.
	r, err := repository.NewApp()
	if err != nil {
		return nil, err
	}
	return &App{
		repository: r,
	}, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	if err := a.repository.Stop(ctx); err != nil {
		return err
	}
	model.Log.Info.Info("usecase layer stopped")
	return nil
}
