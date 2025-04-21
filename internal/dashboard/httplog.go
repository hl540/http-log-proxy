package dashboard

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
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

func (h *Handler) HttpLogListHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	h.String(writer, http.StatusOK, "HttpLogListHandler")
}
