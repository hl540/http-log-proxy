package storage

const AppModelTableName = "tb_app"

type AppModel struct {
	Id       int64  `db:"id" json:"id"`
	Key      string `db:"key" json:"key"`
	Name     string `db:"name" json:"name"`
	Target   string `db:"target" json:"target"`
	CreateAt int64  `db:"create_at" json:"create_at"`
	UpdateAt int64  `db:"update_at" json:"update_at"`
}

const HttpLogModelTableName = "tb_http_log"

type HttpLogModel struct {
	RequestId      string `db:"request_id" json:"request_id"`
	AppId          int64  `db:"app_id" json:"app_id"`
	AppKey         string `db:"app_key" json:"app_key"`
	RequestUrl     string `db:"request_url" json:"request_url"`
	RequestMethod  string `db:"request_method" json:"request_method"`
	RequestHeader  string `db:"request_header" json:"request_header"`
	RequestBody    string `db:"request_body" json:"request_body"`
	ResponseCode   int    `db:"response_code" json:"response_code"`
	ResponseHeader string `db:"response_header" json:"response_header"`
	ResponseBody   string `db:"response_body" json:"response_body"`
	CreateAt       int64  `db:"create_at" json:"create_at"`
}
