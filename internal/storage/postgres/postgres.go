package postgres

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"time"
)

type Postgres struct {
	Db  *sqlx.DB
	Log *slog.Logger
}

func New(log *slog.Logger) (*Postgres, error) {
	dsn := "user=postgres dbname=gopher_chess sslmode=disable"
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Error("Ошибка подключения:", err)
	}
	// Настройка пула соединений
	db.SetMaxOpenConns(25)           // максимальное количество открытых соединений
	db.SetMaxIdleConns(25)           // максимальное количество простаивающих соединений
	db.SetConnMaxLifetime(time.Hour) // время жизни каждого соединения

	return &Postgres{Db: db, Log: log}, nil
}
