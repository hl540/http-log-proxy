package dashboard

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hl540/http-log-proxy/storage"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type Handler struct {
	StorageProvider storage.Provider
}

func NewHandler(storage storage.Provider) *Handler {
	return &Handler{StorageProvider: storage}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.GET("/dashboard/home", h.HomePageHandler)
	router.POST("/dashboard/app/list", h.AppListHandler)
	router.POST("/dashboard/app/new", h.NewAppHandler)
	router.GET("/dashboard/http_log/:request_id", h.HttpLogInfoHandler)
	router.POST("/dashboard/http_log/list", h.HttpLogListHandler)
}
func (h *Handler) Error(writer http.ResponseWriter, err error) {
	http.Error(writer, err.Error(), http.StatusInternalServerError)
}

func (h *Handler) Bytes(writer http.ResponseWriter, code int, buf []byte) {
	writer.WriteHeader(code)
	_, _ = writer.Write(buf)
}

func (h *Handler) String(writer http.ResponseWriter, code int, format string, values ...any) {
	writer.WriteHeader(code)
	_, _ = fmt.Fprintf(writer, format, values...)
}

func (h *Handler) Json(writer http.ResponseWriter, code int, obj interface{}) {
	writer.WriteHeader(code)
	_ = json.NewEncoder(writer).Encode(obj)
}

func (h *Handler) HomePageHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	h.String(writer, http.StatusOK, "HomePageHandler")
}

func (h *Handler) AppListHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	list, err := h.StorageProvider.SearchAppList(request.Context(), "", "")
	if err != nil {
		h.Error(writer, err)
		return
	}
	h.Json(writer, http.StatusOK, list)
}

func (h *Handler) NewAppHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	h.StorageProvider.AddApp(request.Context(), &storage.AppModel{
		Key:      uuid.NewString(),
		Name:     "测试" + time.Now().String(),
		Target:   "http:/127.0.0.1:8080",
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
	})
	h.Json(writer, http.StatusOK, "success")
}

func (h *Handler) HttpLogInfoHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	requestId := params.ByName("request_id")
	httpLog, err := h.StorageProvider.GetHttpLogByRequestId(request.Context(), requestId)
	if err != nil {
		h.Error(writer, err)
		return
	}
	h.Json(writer, http.StatusOK, httpLog)
}

func (h *Handler) HttpLogListHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	h.String(writer, http.StatusOK, "HttpLogListHandler")
}
