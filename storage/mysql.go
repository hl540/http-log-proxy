package storage

import (
	"context"
	"fmt"
	"github.com/hl540/http-log-proxy/config"
	"github.com/jmoiron/sqlx"
)

func init() {
	Register("mysql", new(MySqlStorage))
}

type MySqlStorage struct {
	db *sqlx.DB
}

func (s *MySqlStorage) Init(conf *config.Storage) error {
	//TODO implement me
	panic("implement me")
}

func (s *MySqlStorage) AddApp(ctx context.Context, app *AppModel) error {
	insertSql := fmt.Sprintf("INSERT INTO %s (`key`, `name`, `target`, `create_at`, `update_at`) VALUES (?, ?, ?, ?, ?)", AppModelTableName)
	_, err := s.db.ExecContext(ctx, insertSql, app.Key, app.Name, app.Target, app.CreateAt, app.UpdateAt)
	return err
}

func (s *MySqlStorage) DelApp(ctx context.Context, id int64) error {
	deleteSql := fmt.Sprintf("DELETE FROM %s WHERE `id` = ?", AppModelTableName)
	_, err := s.db.ExecContext(ctx, deleteSql, id)
	return err
}

func (s *MySqlStorage) UpdateApp(ctx context.Context, app *AppModel) error {
	updateSql := fmt.Sprintf("UPDATE %s SET `name` = ?, `target` = ?, `update_at` = ? WHERE `id` = ?", AppModelTableName)
	_, err := s.db.ExecContext(ctx, updateSql, app.Name, app.Target, app.UpdateAt, app.Id)
	return err
}

func (s *MySqlStorage) GetAppById(ctx context.Context, id int64) (*AppModel, error) {
	selectSql := fmt.Sprintf("SELECT * FROM %s WHERE `id` = ? LIMIT 1", AppModelTableName)
	var app AppModel
	err := s.db.GetContext(ctx, &app, selectSql, id)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *MySqlStorage) GetAppByKey(ctx context.Context, key string) (*AppModel, error) {
	selectSql := fmt.Sprintf("SELECT * FROM %s WHERE `key` = ? LIMIT 1", AppModelTableName)
	var app AppModel
	err := s.db.GetContext(ctx, &app, selectSql, key)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *MySqlStorage) SearchAppList(ctx context.Context, name string, key string) ([]*AppModel, error) {
	selectSql := fmt.Sprintf("SELECT * FROM %s WHERE 1 = 1", AppModelTableName)
	var args []any
	if name != "" {
		selectSql = fmt.Sprintf("%s AND `name` LIKE ?", selectSql)
		args = append(args, "%"+name+"%")
	}
	if key != "" {
		selectSql = fmt.Sprintf("%s AND `key` LIKE ?", selectSql)
		args = append(args, "%"+key+"%")
	}
	var apps []*AppModel
	err := s.db.SelectContext(ctx, &apps, selectSql, args...)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *MySqlStorage) AddHttpLog(ctx context.Context, log *HttpLogModel) error {
	insertSql := fmt.Sprintf("INSERT INTO %s (`request_id`, `app_id`, `app_key`, `request_url`, `request_method`, `request_header`, `request_body`, `response_code`, `response_header`, `response_body`, `create_at`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", HttpLogModelTableName)
	_, err := s.db.ExecContext(ctx, insertSql,
		log.RequestId,
		log.AppId,
		log.AppKey,
		log.RequestUrl,
		log.RequestMethod,
		log.RequestHeader,
		log.RequestBody,
		log.ResponseCode,
		log.ResponseHeader,
		log.ResponseBody,
		log.CreateAt,
	)
	return err
}

func (s *MySqlStorage) GetHttpLogByRequestId(ctx context.Context, requestId string) (*HttpLogModel, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MySqlStorage) GetHttpLogListByAppId(ctx context.Context, appId int64, key string, size int64, page int64) ([]*HttpLogModel, error) {
	//TODO implement me
	panic("implement me")
}
