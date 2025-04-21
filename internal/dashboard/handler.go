package dashboard

import (
	"encoding/json"
	"fmt"
	"github.com/hl540/http-log-proxy/storage"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
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
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	_ = json.NewEncoder(writer).Encode(obj)
}

func (h *Handler) View(writer http.ResponseWriter, filename string, data any) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		h.Error(writer, err)
		return
	}
	if err := tmpl.Execute(writer, data); err != nil {
		h.Error(writer, err)
		return
	}
}
