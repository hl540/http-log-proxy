package storage

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/hl540/http-log-proxy/configs"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

func init() {
	Register("mysql", new(MySqlStorage))
}

type MySqlStorage struct {
	db *sqlx.DB
}

func (s *MySqlStorage) Init(conf *configs.Storage) error {
	//TODO implement me
	panic("implement me")
}

//go:embed mysql_init.sql
var mysqlInitSql string

func (s *MySqlStorage) Setup(ctx context.Context) error {
	query := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE `table_schema` = DATABASE() AND `table_name` = ?")
	var count int
	err := s.db.GetContext(ctx, &count, query, AppModelTableName)
	if err != nil {
		return fmt.Errorf("query table exists: %s", err)
	}
	if count > 0 {
		return nil
	}

	_, err = s.db.ExecContext(ctx, mysqlInitSql)
	if err != nil {
		return fmt.Errorf("mysql init: %s", err)
	}
	return nil
}

func (s *MySqlStorage) AddApp(ctx context.Context, app *AppModel) error {
	insertSql := fmt.Sprintf("INSERT INTO %s (`id`, `name`, `target`, `create_at`, `update_at`) VALUES (?, ?, ?, ?, ?)", AppModelTableName)
	_, err := s.db.ExecContext(ctx, insertSql, app.Id, app.Name, app.Target, app.CreateAt, app.UpdateAt)
	return err
}

func (s *MySqlStorage) DelApp(ctx context.Context, id string) error {
	deleteSql := fmt.Sprintf("DELETE FROM %s WHERE `id` = ?", AppModelTableName)
	_, err := s.db.ExecContext(ctx, deleteSql, id)
	return err
}

func (s *MySqlStorage) UpdateApp(ctx context.Context, app *AppModel) error {
	updateSql := fmt.Sprintf("UPDATE %s SET `name` = ?, `target` = ?, `update_at` = ? WHERE `id` = ?", AppModelTableName)
	_, err := s.db.ExecContext(ctx, updateSql, app.Name, app.Target, app.UpdateAt, app.Id)
	return err
}

func (s *MySqlStorage) GetAppById(ctx context.Context, id string) (*AppModel, error) {
	selectSql := fmt.Sprintf("SELECT * FROM %s WHERE `id` = ? LIMIT 1", AppModelTableName)
	var app AppModel
	err := s.db.GetContext(ctx, &app, selectSql, id)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *MySqlStorage) SearchAppList(ctx context.Context, name string, id string) ([]*AppModel, error) {
	selectSql := fmt.Sprintf("SELECT * FROM %s WHERE 1 = 1", AppModelTableName)
	var args []any
	if name != "" {
		selectSql = fmt.Sprintf("%s AND `name` LIKE ?", selectSql)
		args = append(args, "%"+name+"%")
	}
	if id != "" {
		selectSql = fmt.Sprintf("%s AND `id` LIKE ?", selectSql)
		args = append(args, "%"+id+"%")
	}
	var apps []*AppModel
	err := s.db.SelectContext(ctx, &apps, selectSql, args...)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *MySqlStorage) AddHttpLog(ctx context.Context, log *HttpLogModel) error {
	insertSql := fmt.Sprintf("INSERT INTO %s (`request_id`, `app_id`, `request_url`, `request_method`, `request_header`, `request_body`, `response_code`, `response_header`, `response_body`, `create_at`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", HttpLogModelTableName)
	_, err := s.db.ExecContext(ctx, insertSql,
		log.RequestId,
		log.AppId,
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
	selectSql := fmt.Sprintf("SELECT * FROM %s WHERE `request_id` = ?", HttpLogModelTableName)
	var log HttpLogModel
	err := s.db.GetContext(ctx, &log, selectSql, requestId)
	return &log, err
}

func (s *MySqlStorage) SearchHttpLogList(ctx context.Context, appId string, param *SearchHttpLogListParam) (int64, []*HttpLogModel, error) {
	selectSql := fmt.Sprintf("SELECT * FROM %s WHERE `app_id` = ?", HttpLogModelTableName)
	args := []any{appId}
	if param.Keyword != "" {
		selectSql = fmt.Sprintf("%s AND `request_body` LIKE ?", selectSql)
		selectSql = fmt.Sprintf("%s OR `request_body` LIKE ?", selectSql)
		selectSql = fmt.Sprintf("%s OR `response_body` LIKE ?", selectSql)
		selectSql = fmt.Sprintf("%s OR `response_body` LIKE ?", selectSql)
		args = append(args, "%"+param.Keyword+"%")
		args = append(args, "%"+UnicodeForMySQLLike(param.Keyword)+"%")
		args = append(args, "%"+param.Keyword+"%")
		args = append(args, "%"+UnicodeForMySQLLike(param.Keyword)+"%")
	}
	if param.StartTime != 0 {
		selectSql = fmt.Sprintf("%s AND `create_at` >= ?", selectSql)
		args = append(args, param.StartTime)
	}
	if param.EndTime != 0 {
		selectSql = fmt.Sprintf("%s AND `create_at` <= ?", selectSql)
		args = append(args, param.EndTime)
	}

	// 查询总数
	var count int64
	selectCountSql := strings.Replace(selectSql, "*", "COUNT(1)", 1)
	err := s.db.GetContext(ctx, &count, selectCountSql, args...)
	if err != nil {
		return 0, nil, err
	}

	// 查询列表
	list := make([]*HttpLogModel, 0, param.Size)
	selectSql = fmt.Sprintf("%s ORDER BY `create_at` DESC LIMIT %d, %d", selectSql, (param.Page-1)*param.Size, param.Size)
	err = s.db.SelectContext(ctx, &list, selectSql, args...)
	return count, list, err
}

func UnicodeForMySQLLike(s string) string {
	ascii := strconv.QuoteToASCII(s) // 结果例如: "\u4e2d\u6587"
	return ascii[1 : len(ascii)-1]   // 去掉前后引号
}
