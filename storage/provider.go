package storage

import (
	"context"
	"fmt"
	"github.com/hl540/http-log-proxy/config"
)

type AppStorage interface {
	// AddApp 新增app
	AddApp(ctx context.Context, app *AppModel) error
	// DelApp 删除app
	DelApp(ctx context.Context, id int64) error
	// UpdateApp 修改App
	UpdateApp(ctx context.Context, app *AppModel) error
	// GetAppById 通过id获取App
	GetAppById(ctx context.Context, id int64) (*AppModel, error)
	// GetAppByKey 通过key获取App
	GetAppByKey(ctx context.Context, key string) (*AppModel, error)
	// SearchAppList 获取app列表
	SearchAppList(ctx context.Context, name string, key string) ([]*AppModel, error)
}

type HttpLogStorage interface {
	// AddHttpLog 新增http log
	AddHttpLog(ctx context.Context, log *HttpLogModel) error
	// GetHttpLogByRequestId 通过RequestId获取log详情
	GetHttpLogByRequestId(ctx context.Context, requestId string) (*HttpLogModel, error)
	// GetHttpLogListByAppId 通过appid获取log列表，支持分页
	GetHttpLogListByAppId(ctx context.Context, appId int64, key string, size int64, page int64) ([]*HttpLogModel, error)
}

type Provider interface {
	Init(conf *config.Storage) error
	AppStorage
	HttpLogStorage
}

var storages = make(map[string]Provider)

func Register(name string, provider Provider) {
	storages[name] = provider
}

func Load(conf *config.Config) (Provider, error) {
	provider, ok := storages[conf.Storage.Type]
	if !ok {
		return nil, fmt.Errorf("storage type %s not exist", conf.Storage.Type)
	}
	if err := provider.Init(conf.Storage); err != nil {
		return nil, err
	}
	return provider, nil
}
