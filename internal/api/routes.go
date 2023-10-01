package api

import (
	"net/http"

	"github.com/cruffinoni/fizzbuzz/internal/database"
	"github.com/cruffinoni/fizzbuzz/internal/utils"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	db database.RequestHandler
}

func NewRoutes(db database.RequestHandler) *Routes {
	return &Routes{db: db}
}

// Ping godoc
//
//	@Summary		Ping
//	@Description	Returns pong
//	@ID				ping
//	@Produce		json
//	@Success		200	{object}	utils.StatusOKBuilder
//	@Router			/ping [get]
func (r *Routes) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, utils.NewStatusOKBuilder("pong"))
}
