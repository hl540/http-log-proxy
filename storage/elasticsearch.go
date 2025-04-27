package storage

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/hl540/http-log-proxy/configs"
	"github.com/olivere/elastic/v7"
)

func init() {
	Register("elasticsearch", new(ElasticsearchStorage))
}

type ElasticsearchStorage struct {
	es *elastic.Client
}

func (s *ElasticsearchStorage) Init(conf *configs.Storage) error {
	client, err := elastic.NewClient(
		elastic.SetURL(conf.Source),
		elastic.SetBasicAuth(conf.User, conf.Pass),
		elastic.SetSniff(false),
	)
	if err != nil {
		return fmt.Errorf("elasticsearch.NewClient: %s", err)
	}
	s.es = client
	return nil
}

//go:embed elasticsearch_http_log_init.json
var elasticsearchHttpLogMapping string

//go:embed elasticsearch_app_init.json
var elasticsearchAppMapping string

func (s *ElasticsearchStorage) Setup(ctx context.Context) error {
	indexMappings := map[string]string{
		AppModelTableName:     elasticsearchAppMapping,
		HttpLogModelTableName: elasticsearchHttpLogMapping,
	}
	for indexName, mapping := range indexMappings {
		exists, err := s.es.IndexExists(indexName).Do(ctx)
		if err != nil {
			return fmt.Errorf("elasticsearch check %s index exists: %s", indexName, err)
		}
		if exists {
			continue
		}
		_, err = s.es.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			return fmt.Errorf("elasticsearch create %s index: %s", indexName, err)
		}
	}
	return nil
}

func (s *ElasticsearchStorage) AddApp(ctx context.Context, app *AppModel) error {
	_, err := s.es.Index().Index(AppModelTableName).BodyJson(app).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *ElasticsearchStorage) DelApp(ctx context.Context, id string) error {
	query := elastic.NewMatchQuery("id", id)
	searchResult, err := s.es.Search().Index(AppModelTableName).Query(query).Size(1).Do(ctx)
	if err != nil {
		return err
	}
	if searchResult.TotalHits() == 0 {
		return fmt.Errorf("app not found")
	}
	_, err = s.es.Delete().Index(AppModelTableName).Id(searchResult.Hits.Hits[0].Id).Refresh("true").Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *ElasticsearchStorage) UpdateApp(ctx context.Context, app *AppModel) error {
	query := elastic.NewMatchQuery("id", app.Id)
	searchResult, err := s.es.Search().Index(AppModelTableName).Query(query).Size(1).Do(ctx)
	if err != nil {
		return err
	}
	if searchResult.TotalHits() == 0 {
		return fmt.Errorf("app not found")
	}

	_, err = s.es.Update().Index(AppModelTableName).Id(searchResult.Hits.Hits[0].Id).Doc(app).Refresh("true").Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *ElasticsearchStorage) GetAppById(ctx context.Context, id string) (*AppModel, error) {
	query := elastic.NewMatchQuery("id", id)
	searchResult, err := s.es.Search().Index(AppModelTableName).Query(query).Size(1).Do(ctx)
	if err != nil {
		return nil, err
	}
	if searchResult.TotalHits() == 0 {
		return nil, fmt.Errorf("app not found")
	}
	var app AppModel
	err = json.Unmarshal(searchResult.Hits.Hits[0].Source, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *ElasticsearchStorage) SearchAppList(ctx context.Context, name string, id string) ([]*AppModel, error) {
	query := elastic.NewBoolQuery()
	if name != "" {
		query.Must(elastic.NewMatchQuery("name", name))
	}
	if id != "" {
		query.Must(elastic.NewWildcardQuery("id", "*"+id+"*"))
	}
	searchResult, err := s.es.Search().
		Index(AppModelTableName).
		Query(query).
		SortBy(
			elastic.NewFieldSort("create_at").Desc(),
			elastic.NewFieldSort("update_at").Desc(),
		).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var apps []*AppModel
	for _, itme := range searchResult.Hits.Hits {
		var app AppModel
		if err = json.Unmarshal(itme.Source, &app); err != nil {
			return nil, err
		}
		apps = append(apps, &app)
	}
	return apps, nil
}

func (s *ElasticsearchStorage) AddHttpLog(ctx context.Context, log *HttpLogModel) error {
	_, err := s.es.Index().Index(HttpLogModelTableName).BodyJson(log).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *ElasticsearchStorage) GetHttpLogByRequestId(ctx context.Context, requestId string) (*HttpLogModel, error) {
	query := elastic.NewMatchQuery("request_id", requestId)
	searchResult, err := s.es.Search().Index(HttpLogModelTableName).Query(query).Size(1).Do(ctx)
	if err != nil {
		return nil, err
	}
	if searchResult.TotalHits() == 0 {
		return nil, fmt.Errorf("http log not found")
	}
	var httpLog HttpLogModel
	err = json.Unmarshal(searchResult.Hits.Hits[0].Source, &httpLog)
	if err != nil {
		return nil, err
	}
	return &httpLog, nil
}

func (s *ElasticsearchStorage) SearchHttpLogList(ctx context.Context, appId string, param *SearchHttpLogListParam) (int64, []*HttpLogModel, error) {
	query := elastic.NewBoolQuery().Filter(elastic.NewTermQuery("app_id", appId))

	createAtRange := elastic.NewRangeQuery("create_at")
	if param.StartTime != 0 {
		createAtRange.Gte(param.StartTime)
	}
	if param.EndTime != 0 {
		createAtRange.Lte(param.EndTime)
	}
	query.Filter(createAtRange)
	if param.Keyword != "" {
		queryString := elastic.NewQueryStringQuery(param.Keyword).
			DefaultOperator("AND").
			Escape(true).
			Field("request_id").
			Field("request_body").
			Field("response_body")
		query.Filter(queryString)
	}

	searchResult, err := s.es.Search().
		Index(HttpLogModelTableName).
		Query(query).From(int((param.Page - 1) * param.Size)).
		SortBy(elastic.NewFieldSort("create_at").Desc()).
		Size(int(param.Size)).
		Do(ctx)
	if err != nil {
		return 0, nil, err
	}
	list := make([]*HttpLogModel, 0, param.Size)
	for _, item := range searchResult.Hits.Hits {
		var httpLog HttpLogModel
		if err = json.Unmarshal(item.Source, &httpLog); err != nil {
			return 0, nil, err
		}
		list = append(list, &httpLog)
	}
	return searchResult.TotalHits(), list, nil
}
