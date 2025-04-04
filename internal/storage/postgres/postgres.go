package postgres

import (
	"GopherChessParty/internal/config"
	"GopherChessParty/internal/logger"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	Db  *sqlx.DB
	Log *logger.Logger
}

func MustNew(log *logger.Logger, cfg config.Database) *Postgres {
	db := sqlx.MustConnect("pgx", cfg.DBUrl())
	db.SetMaxOpenConns(cfg.MaxOpenConns)   // максимальное количество открытых соединений
	db.SetMaxIdleConns(cfg.MaxIdleConns)   // максимальное количество простаивающих соединений
	db.SetConnMaxLifetime(cfg.MaxTimeLife) // время жизни каждого соединения

	return &Postgres{Db: db, Log: log}
}
