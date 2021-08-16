package router

import (
	"github.com/gin-gonic/gin"
	"github.com/luxarts/jsend-go"
	"loquegasto-backend/internal/defines"
	"net/http"
)

func New() *gin.Engine{
	r := gin.Default()

	mapRoutes(r)

	return r
}

func mapRoutes(r *gin.Engine){
	// DB connectors, rest clients, and other stuff init

	// Repositories init

	// Services init

	// Controllers init

	// Endpoints

	// Health check endpoint
	r.GET(defines.EndpointPing, healthCheck)
}

func healthCheck(ctx *gin.Context){
	ctx.JSON(http.StatusOK, jsend.NewSuccess("pong"))
}