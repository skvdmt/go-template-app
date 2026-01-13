package internal

import "context"

// Delivery Интерфейс транспортного слоя.
type Delivery interface {
	// Запуск.
	Start(ctx context.Context) error
	// Остановка.
	Stop(ctx context.Context) error
}
