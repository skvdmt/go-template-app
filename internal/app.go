package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/skvdmt/go-template-app/internal/delivery"
	"github.com/skvdmt/go-template-app/internal/model"
)

// App Основная структура приложения.
type App struct {
	// Канал сигналов операционной системы
	// для остановки ресурсов.
	exit chan os.Signal
	// Количество запущеных ресурсов.
	source *sync.WaitGroup
	// Контекст приложения.
	ctx context.Context
	// Функция отмены контекста всего приложения.
	cancel context.CancelFunc
	// Транспортный слой.
	delivery Delivery
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Logs.Info.Info(fmt.Sprintf("%s creating", model.APP_NAME))
	a := &App{
		exit:   make(chan os.Signal),
		source: &sync.WaitGroup{},
	}
	var err error
	// Загрузка конфигурации.
	if err = model.LoadConfig(); err != nil {
		return nil, err
	}
	// Создание контекста.
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
	a.source.Add(1)
	go func() {
		var err error
		if err = a.delivery.Start(a.ctx); err != nil {
			model.Errors <- err
		}
		// Завершение работы ресурса приложения.
		a.source.Done()
	}()

	go a.signalHandle()
	return a.errorHanle()
}

// signalHandle Отслеживание сигналов операционной системы.
func (a *App) signalHandle() {
	signal.Notify(a.exit, syscall.SIGTERM)
	<-a.exit
	model.Errors <- a.stop()
}

// errorHandle Обработка канала ошибок.
func (a *App) errorHanle() error {
	err := <-model.Errors
	close(model.Errors)
	return err
}

// stop Остановка приложения.
func (a *App) stop() error {
	// Остановка транспортного слоя из которо по цепочке
	// останавливаются все остальные слои.
	if err := a.delivery.Stop(a.ctx); err != nil {
		return err
	}
	// Ожидание завершения работы ресурсов приложения.
	a.source.Wait()
	// Закрытие канала отслеживающего сигналы операционной системы.
	close(a.exit)
	// Отмена контекста.
	a.cancel()
	model.Logs.Info.Info(fmt.Sprintf("%s stopped", model.APP_NAME))
	return nil
}
