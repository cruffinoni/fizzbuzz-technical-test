package api

import (
	"net/http"

	"github.com/cruffinoni/fizzbuzz/internal/database"
	"github.com/cruffinoni/fizzbuzz/internal/utils"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	db database.Database
}

func NewRoutes(db *database.DB) *Routes {
	return &Routes{db: db}
}
func (r *Routes) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, utils.NewStatusOKBuilder("pong"))
}
