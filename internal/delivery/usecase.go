package delivery

import "context"

// Usecase Интерфейс сервисного слоя.
type Usecase interface {
	// Запуск.
	Start(ctx context.Context) error
	// Остановка.
	Stop(ctx context.Context) error
}
