package repository

import (
	"GopherChessParty/ent"
	"GopherChessParty/internal/dto"
	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

type Connection struct {
	client *ent.Client
}

// MustNewConnection Создание нового подключения.
func MustNewConnection(cfg dto.Database) *Connection {
	drv, err := sql.Open("postgres", cfg.Url())
	if err != nil {
		panic(err)
	}

	db := drv.DB()
	db.SetMaxIdleConns(cfg.MaxIdleConns)   // Максимальное количество простаивающих соединений.
	db.SetMaxOpenConns(cfg.MaxOpenConns)   // Максимальное количество открытых соединений.
	db.SetConnMaxLifetime(cfg.MaxTimeLife) // Максимальное время жизни одного соединения.

	return &Connection{client: ent.NewClient(ent.Driver(drv))}
}
