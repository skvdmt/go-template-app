package delivery

import "context"

// Usecase Интерфейс сервисного слоя.
type Usecase interface {
	// Остановка.
	Stop(ctx context.Context) error
}
