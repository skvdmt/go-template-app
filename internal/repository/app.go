package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/skvdmt/go-template-app/internal/model"
)

const (
	postgres          = "postgres"
	DB_PASSWORD       = "DB_PASSWORD"
	POSTGRES_PASSWORD = "POSTGRES_PASSWORD"
)

// App Репоситорный слой.
type App struct {
	// Соединение с базой данных.
	db *sql.DB
}

// NewApp Конструктор.
func NewApp(ctx context.Context) (*App, error) {
	model.Logs.Info.Info("repository layer creating")
	a := &App{}
	var err error
	// Соединение с базой данных.
	a.db, err = a.openDB()
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Start Запуск.
func (a *App) Start(ctx context.Context) error {
	model.Logs.Info.Info("repository layer starting")
	return nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {

	_, err := a.db.QueryContext(ctx, `SELECT pg_sleep(3);`)
	if err != nil {
		// Игнорируем ошибки прерывания контекста
		if errors.Is(err, context.Canceled) {
			return nil
		}
		if isPostgresErrorCode(err, pqerror.QueryCanceled) {
			return nil
		}
		return err
	}

	// Закрытие соединения с базой данных.
	if err := a.db.Close(); err != nil {
		return err
	}
	model.Logs.Info.Info("disconnect from database")
	model.Logs.Info.Info("repository layer stopped")
	return nil
}

// openDB Устанавливает соединение с базой данных.
func (a *App) openDB() (*sql.DB, error) {
	penv := DB_PASSWORD
	mode, ok := os.LookupEnv(model.MODE)
	if ok && mode == model.Dev {
		penv = POSTGRES_PASSWORD
	}
	pwd, ok := os.LookupEnv(penv)
	if !ok {
		return nil, fmt.Errorf("env %s unset", penv)
	}
	if !ok || mode != model.Dev {
		var err error
		pwd, err = url.QueryUnescape(pwd)
		if err != nil {
			return nil, err
		}
	}
	db, err := sql.Open(postgres, fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		model.Config.Postgres.Host,
		model.Config.Postgres.Port,
		model.Config.Postgres.User, pwd,
		model.Config.Postgres.Database,
	))
	if err != nil {
		return nil, err
	}
	model.Logs.Info.Info("postgres connection success")
	return db, nil
}

// isPostgresErrorCode Поиск кода в ошибках postgres.
func isPostgresErrorCode(err error, errcode pqerror.Code) bool {
	if pgerr, ok := err.(*pq.Error); ok {
		return pgerr.Code == errcode
	}
	return false
}
