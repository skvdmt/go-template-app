package internal

import "context"

// Delivery Интерфейс транспортного слоя.
type Delivery interface {
	// Остановка.
	Stop(ctx context.Context) error
}
