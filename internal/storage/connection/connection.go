package connection

import (
	"GopherChessParty/ent"
	"GopherChessParty/internal/config"
	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

type Repository struct {
	client *ent.Client
}

// MustNewRepository Создание нового подключения.
func MustNewRepository(cfg config.Database) *Repository {
	drv, err := sql.Open("postgres", cfg.DBUrl())
	if err != nil {
		panic(err)
	}

	db := drv.DB()
	db.SetMaxIdleConns(cfg.MaxIdleConns)   // Максимальное количество простаивающих соединений.
	db.SetMaxOpenConns(cfg.MaxOpenConns)   // Максимальное количество открытых соединений.
	db.SetConnMaxLifetime(cfg.MaxTimeLife) // Максимальное время жизни одного соединения.

	return &Repository{client: ent.NewClient(ent.Driver(drv))}
}
