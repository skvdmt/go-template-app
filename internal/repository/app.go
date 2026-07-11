package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/skvdmt/go-template-app/internal/model"
)

const (
	POSTGRES_PASSWORD = "POSTGRES_PASSWORD"
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
	a.db, err = a.openDB()
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
	defer model.Log.Info.Info("postgres database connection closed")
	return a.db.Close()
}

// openDB Открывает соединение с базой данных.
func (a *App) openDB() (*sql.DB, error) {
	model.Log.Info.Info("postgres database connection opening")
	pwd, o := os.LookupEnv(POSTGRES_PASSWORD)
	if !o {
		return nil, fmt.Errorf("env %s is not set", POSTGRES_PASSWORD)
	}
	c, err := pq.NewConnectorConfig(pq.Config{
		User:     model.Conf.Postgres.User,
		Password: pwd,
		Host:     model.Conf.Postgres.Host,
		Port:     model.Conf.Postgres.Port,
		Database: model.Conf.Postgres.Database,
		SSLMode:  pq.SSLModeDisable,
	})
	if err != nil {
		return nil, err
	}
	db := sql.OpenDB(c)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	model.Log.Info.Info("postgres database connection success")
	return db, nil
}

// isPostgresErrorCode Поиск кода в ошибках postgres.
func (a *App) isPostgresErrorCode(err error, errcode pqerror.Code) bool {
	if pgerr, ok := err.(*pq.Error); ok {
		return pgerr.Code == errcode
	}
	return false
}
