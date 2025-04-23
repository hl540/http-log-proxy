package dashboard

import (
	"github.com/gin-gonic/gin"
	"github.com/hl540/http-log-proxy/storage"
)

type Handler struct {
	StorageProvider storage.Provider
}

func NewHandler(storage storage.Provider) *Handler {
	return &Handler{StorageProvider: storage}
}

func (h *Handler) Register(router gin.RouterGroup) {
	router.GET("/dashboard/home", h.HomePageHandler)
	router.POST("/dashboard/app/list", h.AppListHandler)
	router.POST("/dashboard/app/new", h.NewAppHandler)
	router.GET("/dashboard/http_log/:request_id", h.HttpLogInfoHandler)
	router.POST("/dashboard/http_log/list", h.HttpLogListHandler)
}
