package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"

	_ "github.com/lib/pq"
	"github.com/skvdmt/go-template-app/internal/model"
)

const (
	DB_PASSWORD     = "DB_PASSWORD"
	POSTGRES_DRIVER = "postgres"
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

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
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
	pwd, ok := os.LookupEnv(DB_PASSWORD)
	if !ok {
		return nil, fmt.Errorf("env %v not set", DB_PASSWORD)
	}
	pwd, err := url.QueryUnescape(pwd)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(POSTGRES_DRIVER, fmt.Sprintf(
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
