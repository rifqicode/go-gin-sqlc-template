package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.gift.id/lv2/loyalty/configs/shared"
)

type PingResponse struct {
	Alive       bool `json:"alive"`
	DBConnected bool `json:"db_connected"`
}

// GetHandler godoc
// @Summary ping server and database
// @Schemes http http
// @Description do ping and check database connection
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {object} PingResponse
// @Router /ping [get]
func GetHandler(meta *shared.SharedMeta) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbConnected := meta.DB.Ping(ctx) == nil
		r := PingResponse{
			Alive:       true,
			DBConnected: dbConnected,
		}
		ctx.JSON(http.StatusOK, r)
	}
}
