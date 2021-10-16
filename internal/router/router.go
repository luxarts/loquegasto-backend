package router

import (
	"fmt"
	"loquegasto-backend/internal/controller"
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/middleware"
	"loquegasto-backend/internal/repository"
	"loquegasto-backend/internal/service"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"

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
	db, err := sqlx.Open("postgres", os.Getenv(defines.EnvPostgreSQLDBURI))
	if err != nil {
		panic(fmt.Sprintf("Fail to connect to database: %v", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Fail to ping to database: %v", err))
	}

	// Repositories init
	txnRepo := repository.NewTransactionsRepository(db)
	usersRepo := repository.NewUsersRepository(db)
	walletsRepo := repository.NewWalletRepository(db)

	// Services init
	txnSrv := service.NewTransactionsService(txnRepo, walletsRepo)
	usersSrv := service.NewUsersService(usersRepo)
	walletsSrv := service.NewWalletsService(walletsRepo)

	// Controllers init
	txnCtrl := controller.NewTransactionsController(txnSrv)
	usersCtrl := controller.NewUsersController(usersSrv)
	walletsCtrl := controller.NewWalletsController(walletsSrv)

	// Middleware
	authMw := middleware.NewAuthMiddleware()

	// Endpoints
	// Transactions
	r.POST(defines.EndpointTransactionsCreate, authMw.Check, txnCtrl.Create)
	r.PUT(defines.EndpointTransactionsUpdateByMsgID, authMw.Check, txnCtrl.UpdateByMsgID)
	r.GET(defines.EndpointTransactionsGetAllByUserID, authMw.Check, txnCtrl.GetAllByUserID)
	// Users
	r.POST(defines.EndpointUsersCreate, authMw.Check, usersCtrl.Create)
	// Wallets
	r.POST(defines.EndpointWalletsCreate, authMw.Check, walletsCtrl.Create)
	r.GET(defines.EndpointWalletsGetAll, authMw.Check, walletsCtrl.GetAll)
	r.GET(defines.EndpointWalletsGetByID, authMw.Check, walletsCtrl.GetByID)
	r.PUT(defines.EndpointWalletsUpdateByID, authMw.Check, walletsCtrl.UpdateByID)
	r.DELETE(defines.EndpointWalletsDeleteByID, authMw.Check, walletsCtrl.DeleteByID)

	// Health check endpoint
	r.GET(defines.EndpointPing, healthCheck)
}

func healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, jsend.NewSuccess("pong"))
}
