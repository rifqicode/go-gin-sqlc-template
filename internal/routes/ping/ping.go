package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingResponse struct {
	Alive bool `json:"alive"`
}

// GetHandler godoc
// @Summary ping server
// @Schemes http http
// @Description do ping
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {object} PingResponse
// @Router /ping [get]
func GetHandler(ctx *gin.Context) {
	r := PingResponse{
		Alive: true,
	}
	ctx.JSON(http.StatusOK, r)
	return
}
