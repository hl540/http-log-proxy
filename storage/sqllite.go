package storage

import (
	"fmt"
	"github.com/hl540/http-log-proxy/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	Register("local", new(SQLiteStorage))
}

type SQLiteStorage struct {
	db *sqlx.DB
	*MySqlStorage
}

func (s *SQLiteStorage) Init(conf *config.Storage) error {
	db, err := sqlx.Open("sqlite3", conf.Source)
	if err != nil {
		return fmt.Errorf("sql.Open: %s", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("db.Ping: %s", err)
	}
	s.db = db
	s.MySqlStorage = &MySqlStorage{
		db: s.db,
	}
	return nil
}
