package internal

import (
	"context"
	"errors"
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
	// Канал сигналов прерывания.
	interrupt chan os.Signal
	// Контекст.
	ctx context.Context
	// Отмена контекста.
	cancel context.CancelFunc
	// Транспортный слой.
	delivery Delivery
	// Ошибки приложения.
	errors []error
	// Приложение закрывается.
	stopping bool
	// Корректное завершение горутин.
	wg *sync.WaitGroup
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Log.Info.Info(fmt.Sprintf("app %s creating", model.NAME))
	var err error
	// Загрузка конфигурации.
	if model.Conf, err = model.NewConfig(); err != nil {
		return nil, err
	}
	// Создание глобального канала ошибок.
	model.Errors = make(chan error)
	// Создание транспортного слоя.
	d, err := delivery.NewApp()
	if err != nil {
		return nil, err
	}
	// Создание контекста.
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:       ctx,
		cancel:    cancel,
		interrupt: make(chan os.Signal),
		delivery:  d,
		wg:        &sync.WaitGroup{},
	}, nil
}

// Start Запуск.
func (a *App) Start() []error {
	model.Log.Info.Info(fmt.Sprintf("app %s starting", model.NAME))

	a.wg.Add(1)
	go a.errorHandler()

	return a.interruptHandler()
}

// errorHandle Обработчик глобального канала ошибок.
func (a *App) errorHandler() {
	model.Log.Info.Info("error handler starting")
	defer func() {
		model.Log.Info.Info("error handler stopped")
		a.wg.Done()
	}()
	for {
		err := <-model.Errors
		if err == nil {
			return
		}
		if errors.Is(err, context.Canceled) {
			continue
		}
		a.errors = append(a.errors, err)
		if !a.stopping {
			a.interrupt <- syscall.SIGTERM
		}
	}
}

// interruptHandler Обработчик канала сигналов прерывания.
func (a *App) interruptHandler() []error {
	model.Log.Info.Info("waiting for interrupt signals")
	signal.Notify(a.interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-a.interrupt
	a.stop()
	close(model.Errors)
	return a.errors
}

// stop Остановка.
func (a *App) stop() {
	defer model.Log.Info.Info(fmt.Sprintf("app %s stopped", model.NAME))
	a.stopping = true
	// Отмена контекста.
	a.cancel()
	model.Log.Info.Info("context canceled")
	// Остановка транспортного слоя.
	if err := a.delivery.Stop(a.ctx); err != nil {
		model.Errors <- err
	}
	// Закрытие канала сигналов прерывания.
	close(a.interrupt)
}
