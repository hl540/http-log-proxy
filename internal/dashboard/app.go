package dashboard

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hl540/http-log-proxy/form"
	"github.com/hl540/http-log-proxy/storage"
	"net/http"
	"time"
)

func (h *Handler) AppListHandler(ctx *gin.Context) {
	var req form.AppListReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	list, err := h.StorageProvider.SearchAppList(ctx, req.Name, req.Key)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	data := form.AppListResp{
		Data: make([]*form.AppListRespDataItem, 0, len(list)),
	}
	for _, item := range list {
		data.Data = append(data.Data, &form.AppListRespDataItem{
			Id:       item.Id,
			Name:     item.Name,
			Key:      item.Key,
			CreateAt: time.Unix(item.CreateAt, 0).Format("2006-01-02 15:04:05"),
			UpdateAt: time.Unix(item.UpdateAt, 0).Format("2006-01-02 15:04:05"),
		})
	}
	ctx.JSON(http.StatusOK, list)
}

func (h *Handler) NewAppHandler(ctx *gin.Context) {
	var req form.NewAppReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	app := &storage.AppModel{
		Key:      h.makeAppKey(ctx, req.Target),
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

	data := form.NewAppResp{
		Id:       app.Id,
		Key:      app.Key,
		Name:     app.Name,
		Target:   app.Name,
		CreateAt: time.Unix(app.CreateAt, 0).Format("2006-01-02 15:04:05"),
		UpdateAt: time.Unix(app.UpdateAt, 0).Format("2006-01-02 15:04:05"),
	}
	ctx.JSON(http.StatusOK, data)
}

func (h *Handler) makeAppKey(ctx context.Context, name string) string {
	hash := md5.New()
	hash.Write([]byte(name))
	hash.Write([]byte(time.Now().String()))
	hash.Sum(nil)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
