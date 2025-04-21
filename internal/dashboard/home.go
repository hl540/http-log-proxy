package dashboard

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *Handler) HomePageHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	h.String(writer, http.StatusOK, "HomePageHandler")
}
