package storage

import (
	"context"
	"fmt"
	"github.com/hl540/http-log-proxy/configs"
)

type AppStorage interface {
	// AddApp 新增app
	AddApp(ctx context.Context, app *AppModel) error
	// DelApp 删除app
	DelApp(ctx context.Context, id string) error
	// UpdateApp 修改App
	UpdateApp(ctx context.Context, app *AppModel) error
	// GetAppById 通过id获取App
	GetAppById(ctx context.Context, id string) (*AppModel, error)
	// SearchAppList 获取app列表
	SearchAppList(ctx context.Context, name string, id string) ([]*AppModel, error)
}

type SearchHttpLogListParam struct {
	Keyword   string
	StartTime int64
	EndTime   int64
	Size      int64
	Page      int64
}

type HttpLogStorage interface {
	// AddHttpLog 新增http log
	AddHttpLog(ctx context.Context, log *HttpLogModel) error
	// GetHttpLogByRequestId 通过RequestId获取log详情
	GetHttpLogByRequestId(ctx context.Context, requestId string) (*HttpLogModel, error)
	// SearchHttpLogList 通过appid获取log列表，支持分页
	SearchHttpLogList(ctx context.Context, appId string, param *SearchHttpLogListParam) (int64, []*HttpLogModel, error)
}

type Provider interface {
	Init(conf *configs.Storage) error
	AppStorage
	HttpLogStorage
}

var storages = make(map[string]Provider)

func Register(name string, provider Provider) {
	storages[name] = provider
}

func Load(conf *configs.Config) (Provider, error) {
	provider, ok := storages[conf.Storage.Type]
	if !ok {
		return nil, fmt.Errorf("storage type %s not exist", conf.Storage.Type)
	}
	if err := provider.Init(conf.Storage); err != nil {
		return nil, err
	}
	return provider, nil
}
