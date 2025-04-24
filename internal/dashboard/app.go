package dashboard

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hl540/http-log-proxy/forms"
	"github.com/hl540/http-log-proxy/storage"
	"net/http"
	"time"
)

func (h *Handler) AppListHandler(ctx *gin.Context) {
	var req forms.AppListReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	list, err := h.StorageProvider.SearchAppList(ctx, req.Name, req.Id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	data := forms.AppListResp{
		Data: make([]*forms.AppListRespDataItem, 0, len(list)),
	}
	for _, item := range list {
		data.Data = append(data.Data, &forms.AppListRespDataItem{
			Id:       item.Id,
			Name:     item.Name,
			Target:   item.Target,
			CreateAt: time.Unix(item.CreateAt, 0).Format("2006-01-02 15:04:05"),
			UpdateAt: time.Unix(item.UpdateAt, 0).Format("2006-01-02 15:04:05"),
		})
	}
	ctx.JSON(http.StatusOK, data)
}

func (h *Handler) NewAppHandler(ctx *gin.Context) {
	var req forms.NewAppReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	app := &storage.AppModel{
		Id:       h.makeAppId(ctx, req.Target),
		Name:     req.Name,
		Target:   req.Target,
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
	}
	err := h.StorageProvider.AddApp(ctx, app)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	data := forms.NewAppResp{
		Id:       app.Id,
		Name:     app.Name,
		Target:   app.Name,
		CreateAt: time.Unix(app.CreateAt, 0).Format("2006-01-02 15:04:05"),
		UpdateAt: time.Unix(app.UpdateAt, 0).Format("2006-01-02 15:04:05"),
	}
	ctx.JSON(http.StatusOK, data)
}

func (h *Handler) makeAppId(ctx context.Context, name string) string {
	hash := md5.New()
	hash.Write([]byte(name))
	hash.Write([]byte(time.Now().String()))
	hash.Sum(nil)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (h *Handler) DelAppHandler(ctx *gin.Context) {
	appId := ctx.Param("app_id")
	if appId == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("missing app_id"))
		return
	}
	err := h.StorageProvider.DelApp(ctx, appId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"count": 1,
	})
}
