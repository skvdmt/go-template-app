package repository

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/skvdmt/go-template-app/internal/model"
)

// App Репоситорный слой.
type App struct {
	// Соединение с базой данных.
	db *sql.DB
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	model.Log.Info.Info("repository layer creating")
	a := &App{}
	var err error
	// Соединение с базой данных.
	a.db, err = OpenPostgresDB()
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	if err := a.closeDB(); err != nil {
		return err
	}
	model.Log.Info.Info("repository layer stopped")
	return nil
}

// closeDB Закрывает соединение с базой данных.
func (a *App) closeDB() error {
	defer model.Log.Info.Info("database connection closed")
	return a.db.Close()
}
