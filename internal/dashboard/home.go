package dashboard

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) HomePageHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "home_page.html", nil)
}
