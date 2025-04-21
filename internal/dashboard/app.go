package dashboard

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/hl540/http-log-proxy/storage"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func (h *Handler) AppListHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	list, err := h.StorageProvider.SearchAppList(request.Context(), "", "")
	if err != nil {
		h.Error(writer, err)
		return
	}
	h.Json(writer, http.StatusOK, list)
}

func (h *Handler) NewAppHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	app := &storage.AppModel{
		Key:      h.makeAppKey(request.Context(), "xxx"),
		Name:     "测试" + time.Now().String(),
		Target:   "http://127.0.0.1:8000",
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
	}
	err := h.StorageProvider.AddApp(request.Context(), app)
	if err != nil {
		h.Json(writer, http.StatusInternalServerError, err.Error())
		return
	}
	h.Json(writer, http.StatusOK, app)
}

func (h *Handler) makeAppKey(ctx context.Context, name string) string {
	hash := md5.New()
	hash.Write([]byte(name))
	hash.Write([]byte(time.Now().String()))
	hash.Sum(nil)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
