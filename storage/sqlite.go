package storage

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/hl540/http-log-proxy/configs"
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

func (s *SQLiteStorage) Init(conf *configs.Storage) error {
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

//go:embed sqlite_init.sql
var sqliteInitSql string

func (s *SQLiteStorage) Setup(ctx context.Context) error {
	query := fmt.Sprint("SELECT COUNT(*) FROM sqlite_master WHERE `type` = 'table' AND `name` = ? ")
	var count int
	err := s.db.GetContext(ctx, &count, query, AppModelTableName)
	if err != nil {
		return fmt.Errorf("query table exists: %s", err)
	}
	if count > 0 {
		return nil
	}

	_, err = s.db.ExecContext(ctx, sqliteInitSql)
	if err != nil {
		return fmt.Errorf("sqlite init: %s", err)
	}
	return nil
}
