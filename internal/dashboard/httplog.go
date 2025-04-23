package dashboard

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func (h *Handler) HttpLogInfoHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	requestId := params.ByName("request_id")
	httpLog, err := h.StorageProvider.GetHttpLogByRequestId(request.Context(), requestId)
	if err != nil {
		h.Error(writer, err)
		return
	}
	h.Json(writer, http.StatusOK, httpLog)
}

type HttpLogListRespDataItem struct {
	CreateAt      string `json:"create_at"`
	RequestId     string `json:"request_id"`
	RequestUrl    string `json:"request_url"`
	RequestMethod string `json:"request_method"`
	ResponseCode  int    `json:"response_code"`
}

func (h *Handler) HttpLogListHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	count, list, err := h.StorageProvider.SearchHttpLogList(request.Context(), 10, "", 15, 1)
	if err != nil {
		h.Error(writer, err)
		return
	}
	data := make([]*HttpLogListRespDataItem, 0, len(list))
	for _, item := range list {
		data = append(data, &HttpLogListRespDataItem{
			CreateAt:      time.Unix(item.CreateAt, 0).Format("2006-01-02 15:04:05"),
			RequestId:     item.RequestId,
			RequestUrl:    item.RequestUrl,
			RequestMethod: item.RequestMethod,
			ResponseCode:  item.ResponseCode,
		})
	}
	h.Json(writer, http.StatusOK, map[string]any{
		"total": count,
		"data":  data,
	})
}
