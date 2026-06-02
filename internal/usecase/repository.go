package usecase

import "context"

// Repository Интерфейс репозиторного слоя.
type Repository interface {
	// Запуск.
	Start(ctx context.Context) error
	// Остановка.
	Stop(ctx context.Context) error
}
