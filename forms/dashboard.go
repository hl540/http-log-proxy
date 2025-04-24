package forms

type AppListReq struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type AppListResp struct {
	Data []*AppListRespDataItem `json:"data"`
}

type AppListRespDataItem struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Target   string `json:"target"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
}

type NewAppReq struct {
	Name   string `form:"name" binding:"required,max=50"`
	Target string `form:"target" binding:"required"`
}

type NewAppResp struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Target   string `json:"target"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
}

type HttpLogListReq struct {
	AppId     string `json:"app_id" binding:"required"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Keyword   string `json:"keyword"`
	Page      int64  `json:"page" binding:"min=1"`
	Size      int64  `json:"size" binding:"min=1,max=100"`
}

type HttpLogListResp struct {
	Total int64                      `json:"total"`
	Data  []*HttpLogListRespDataItem `json:"data"`
}

type HttpLogListRespDataItem struct {
	CreateAt      string `json:"create_at"`
	RequestId     string `json:"request_id"`
	RequestUrl    string `json:"request_url"`
	RequestMethod string `json:"request_method"`
	ResponseCode  int    `json:"response_code"`
}
