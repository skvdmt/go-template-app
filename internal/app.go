package internal

import (
	"context"
	"errors"
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
	// Приложение запущено.
	started bool
	// Ошибки приложения.
	eg []error
	// Приложение уже закрывается.
	stopping bool
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Logs.Info.Info(fmt.Sprintf("%s creating", model.APP_NAME))
	// Загрузка конфигурации.
	if err := model.LoadConfig(); err != nil {
		return nil, err
	}
	// Создаем глобальный канал ошибок.
	model.Errors = make(chan error)
	a := &App{
		interrupt: make(chan os.Signal),
	}
	// Создание контекста.
	a.ctx, a.cancel = context.WithCancel(context.Background())
	// Создание транспортного слоя из которого по
	// цепочки создаются остальные слои приложения.
	var err error
	a.delivery, err = delivery.NewApp(a.ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Start Запуск приложения.
func (a *App) Start() error {
	model.Logs.Info.Info(fmt.Sprintf("%s starting", model.APP_NAME))
	go a.errorHandler()
	go func() {
		// Запуск слоев приложения по цепочке.
		if err := a.delivery.Start(a.ctx); err != nil {
			model.Errors <- err
		}
		a.started = true
	}()
	return a.interruptHandler()
}

// errorHandle Обработчик глобального канала ошибок.
func (a *App) errorHandler() {
	model.Logs.Info.Info("error handler starting")
	for {
		err := <-model.Errors
		if err == nil {
			return
		}
		if errors.Is(err, context.Canceled) {
			continue
		}
		a.eg = append(a.eg, err)
		if !a.stopping {
			a.interrupt <- syscall.SIGTERM
		}
	}
}

// interruptHandler Обработчик сигналов остановки приложения.
func (a *App) interruptHandler() error {
	model.Logs.Info.Info("error handler starting")
	signal.Notify(a.interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-a.interrupt
	// close sources
	if err := a.stop(); err != nil {
		model.Errors <- err
	}
	model.Errors <- nil
	close(model.Errors)
	if len(a.eg) > 0 {
		return fmt.Errorf("%v", a.eg)
	}
	return nil
}

// stop Остановка приложения.
func (a *App) stop() error {
	a.stopping = true

	// Отмена контекста.
	a.cancel()
	model.Logs.Info.Info("context canceled")

	// Дождаться запуска приложения.
	for {
		if a.started {
			break
		}
	}

	// Остановка транспортного слоя из которо по цепочке
	// останавливаются все остальные слои.
	if err := a.delivery.Stop(a.ctx); err != nil {
		return err
	}
	// Закрытие канала отслеживающего сигналы
	// прерывания операционной системы.
	close(a.interrupt)
	model.Logs.Info.Info(fmt.Sprintf("%s stopped", model.APP_NAME))
	return nil
}
