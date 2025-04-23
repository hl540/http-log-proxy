package dashboard

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hl540/http-log-proxy/form"
	"github.com/hl540/http-log-proxy/storage"
	"net/http"
	"time"
)

func (h *Handler) HttpLogInfoHandler(ctx *gin.Context) {
	requestId := ctx.Param("request_id")
	if requestId == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("missing request_id"))
		return
	}

	httpLog, err := h.StorageProvider.GetHttpLogByRequestId(ctx, requestId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, httpLog)
}

func (h *Handler) HttpLogListHandler(ctx *gin.Context) {
	var req form.HttpLogListReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	total, list, err := h.StorageProvider.SearchHttpLogList(ctx, req.AppId, &storage.SearchHttpLogListParam{
		Keyword:   req.Keyword,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Size:      req.Size,
		Page:      req.Page,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	data := &form.HttpLogListResp{
		Total: total,
		Data:  make([]*form.HttpLogListRespDataItem, 0, len(list)),
	}
	for _, item := range list {
		data.Data = append(data.Data, &form.HttpLogListRespDataItem{
			CreateAt:      time.Unix(item.CreateAt, 0).Format("2006-01-02 15:04:05"),
			RequestId:     item.RequestId,
			RequestUrl:    item.RequestUrl,
			RequestMethod: item.RequestMethod,
			ResponseCode:  item.ResponseCode,
		})
	}
	ctx.JSON(http.StatusOK, data)
}
