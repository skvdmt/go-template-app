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
func NewApp(ctx context.Context) (*App, error) {
	model.Logs.Info.Info("usecase layer creating")
	a := &App{}
	var err error
	// Создание репозиторного слоя.
	a.repository, err = repository.NewApp(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	if err := a.repository.Stop(ctx); err != nil {
		return err
	}
	model.Logs.Info.Info("usecase layer stopped")
	return nil
}
