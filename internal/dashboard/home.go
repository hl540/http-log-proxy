package dashboard

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) HomePageHandler(ctx *gin.Context) {
	key := h.makeAppId(ctx, "Default")
	fmt.Println(key)
	ctx.HTML(http.StatusOK, "home_page.html", nil)
}
