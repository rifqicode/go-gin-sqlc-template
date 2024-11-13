package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.gift.id/lv2/loyalty/internal/db"
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
func GetHandler(query *db.Queries) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := query.CheckDBConnection(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, PingResponse{
				Alive:       false,
				DBConnected: false,
			})
		}

		r := PingResponse{
			Alive:       true,
			DBConnected: true,
		}
		ctx.JSON(http.StatusOK, r)
	}
}
