package router

import (
	"loquegasto-backend/internal/controller"
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/middleware"
	"loquegasto-backend/internal/repository"
	"loquegasto-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luxarts/jsend-go"
)

func New() *gin.Engine {
	r := gin.Default()

	mapRoutes(r)

	return r
}

func mapRoutes(r *gin.Engine) {
	// DB connectors, rest clients, and other stuff init

	// Repositories init
	txnRepo := repository.NewTransactionsRepository()

	// Services init
	txnSrv := service.NewTransactionsService(txnRepo)

	// Controllers init
	txnCtrl := controller.NewTransactionsController(txnSrv)

	// Middleware
	authMw := middleware.NewAuthMiddleware()

	// Endpoints
	r.POST(defines.EndpointTransactionsCreate, authMw.Check, txnCtrl.Create)
	r.PUT(defines.EndpointTransactionsUpdateByMsgID, authMw.Check, txnCtrl.UpdateByMsgID)
	r.GET(defines.EndpointTransactionsGetAllByUserID, authMw.Check, txnCtrl.GetAllByUserID)

	// Health check endpoint
	r.GET(defines.EndpointPing, healthCheck)
}

func healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, jsend.NewSuccess("pong"))
}
