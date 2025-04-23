package storage

import (
	"context"
	"github.com/hl540/http-log-proxy/config"
)

func init() {
	Register("elasticsearch", new(ElasticsearchStorage))
}

type ElasticsearchStorage struct{}

func (s *ElasticsearchStorage) Init(conf *config.Storage) error {
	//TODO implement me
	panic("implement me")
}

func (s *ElasticsearchStorage) AddApp(ctx context.Context, app *AppModel) error {
	//TODO implement me
	panic("implement me")
}

func (s *ElasticsearchStorage) DelApp(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (s *ElasticsearchStorage) UpdateApp(ctx context.Context, app *AppModel) error {
	//TODO implement me
	panic("implement me")
}

func (s *ElasticsearchStorage) GetAppById(ctx context.Context, id int64) (*AppModel, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ElasticsearchStorage) GetAppByKey(ctx context.Context, key string) (*AppModel, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ElasticsearchStorage) SearchAppList(ctx context.Context, name string, key string) ([]*AppModel, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ElasticsearchStorage) AddHttpLog(ctx context.Context, log *HttpLogModel) error {
	//TODO implement me
	panic("implement me")
}

func (s *ElasticsearchStorage) GetHttpLogByRequestId(ctx context.Context, requestId string) (*HttpLogModel, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ElasticsearchStorage) SearchHttpLogList(ctx context.Context, appId int64, param *SearchHttpLogListParam) (int64, []*HttpLogModel, error) {
	//TODO implement me
	panic("implement me")
}
