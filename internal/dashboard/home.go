package dashboard

import (
	"github.com/hl540/http-log-proxy/storage"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func (h *Handler) HomePageHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	//appKey := request.Form.Get("app_key")
	appList, err := h.StorageProvider.SearchAppList(request.Context(), "", "")
	if err != nil {
		log.Printf("Error searching app list: %v", err)
		appList = make([]*storage.AppModel, 0)
	}

	logTotal, logList, err := h.StorageProvider.SearchHttpLogList(request.Context(), appList[0].Id, "", 100, 1)
	if err != nil {
		log.Printf("Error searching app list: %v", err)
		logList = make([]*storage.HttpLogModel, 0)
	}
	h.View(writer, "./template/home_page.html", map[string]any{
		"AppList":  appList,
		"LogList":  logList,
		"LogTotal": logTotal,
	})
}
