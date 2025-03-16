package postgres

import (
	"github.com/George-c0de/GopherChessParty/internal/config"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type Postgres struct {
	Db  *sqlx.DB
	Log *slog.Logger
}

func MustNew(log *slog.Logger, cfg config.Database) *Postgres {
	db := sqlx.MustConnect("pgx", cfg.DBUrl())
	db.SetMaxOpenConns(cfg.MaxOpenConns)   // максимальное количество открытых соединений
	db.SetMaxIdleConns(cfg.MaxIdleConns)   // максимальное количество простаивающих соединений
	db.SetConnMaxLifetime(cfg.MaxTimeLife) // время жизни каждого соединения

	return &Postgres{Db: db, Log: log}
}
