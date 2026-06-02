package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/skvdmt/go-template-app/internal/delivery"
	"github.com/skvdmt/go-template-app/internal/model"
)

// App Основная структура приложения.
type App struct {
	// Канал сигналов операционной системы для отслеживания
	// сигналов прерывания работы приложения.
	interrupt chan os.Signal
	// Контекст приложения.
	ctx context.Context
	// Функция отмены контекста всего приложения.
	cancel context.CancelFunc
	// Транспортный слой.
	delivery Delivery
	// Приложение можно останавливать.
	canStop bool
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Logs.Info.Info(fmt.Sprintf("%s creating", model.APP_NAME))
	a := &App{
		interrupt: make(chan os.Signal),
	}
	var err error
	// Загрузка конфигурации.
	if err = model.LoadConfig(); err != nil {
		return nil, err
	}
	// Создание контекста.
	// a.ctx, a.cancel = context.WithTimeout(context.Background(), time.Nanosecond*1)
	a.ctx, a.cancel = context.WithCancel(context.Background())
	// Создание транспортного слоя из которого по
	// цепочки создаются остальные слои приложения.
	a.delivery, err = delivery.NewApp(a.ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Start Запуск приложения.
func (a *App) Start() error {
	model.Logs.Info.Info(fmt.Sprintf("%s starting", model.APP_NAME))
	// Создане глобального канала ошибок для всего приложения.
	model.Errors = make(chan error)
	// Начало работы ресурса приложения.
	go func() {
		var err error
		if err = a.delivery.Start(a.ctx); err != nil {
			model.Errors <- err
		}
		a.canStop = true
	}()
	go a.signalHandle()
	return a.errorHanle()
}

// signalHandle Отслеживание сигналов операционной системы.
func (a *App) signalHandle() {
	signal.Notify(a.interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-a.interrupt
	model.Errors <- nil
}

// errorHandle Обработка канала ошибок.
func (a *App) errorHanle() error {
	err := <-model.Errors
	// Завершение работы ресурсов приложения.
	if err2 := a.stop(); err2 != nil {
		return fmt.Errorf("%w; %v", err, err2)
	}
	return err
}

// stop Остановка приложения.
func (a *App) stop() error {
	for {
		if a.canStop {
			break
		}
	}

	// Отмена контекста.
	a.cancel()
	model.Logs.Info.Info("context canceled")

	// Остановка транспортного слоя из которо по цепочке
	// останавливаются все остальные слои.
	if err := a.delivery.Stop(a.ctx); err != nil {
		return err
	}
	// Закрытие канала отслеживающего сигналы
	// прерывания операционной системы.
	close(a.interrupt)
	// Закрытие канала ошибок.
	close(model.Errors)
	model.Logs.Info.Info(fmt.Sprintf("%s stopped", model.APP_NAME))
	return nil
}
