package router

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"loquegasto-backend/internal/controller"
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/middleware"
	"loquegasto-backend/internal/repository"
	"loquegasto-backend/internal/service"
	"net/http"
	"os"

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
	postgresURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv(defines.EnvPostgresUser),
		os.Getenv(defines.EnvPostgresPassword),
		os.Getenv(defines.EnvPostgresHost),
		os.Getenv(defines.EnvPostgresPort),
		os.Getenv(defines.EnvPostgresDB),
	)

	db, err := sqlx.Open("postgres", postgresURI)
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
	catRepo := repository.NewCategoriesRepository(db)

	// Services init
	txnSvc := service.NewTransactionsService(txnRepo, walletsRepo, catRepo)
	usersSvc := service.NewUsersService(usersRepo)
	walletsSvc := service.NewWalletsService(walletsRepo)
	catSvc := service.NewCategoriesService(catRepo)

	// Controllers init
	txnCtrl := controller.NewTransactionsController(txnSvc)
	usersCtrl := controller.NewUsersController(usersSvc)
	walletsCtrl := controller.NewWalletsController(walletsSvc)
	catCtrl := controller.NewCategoriesController(catSvc)

	// Middleware
	authMw := middleware.NewAuthMiddleware()
	appAuthMw := middleware.NewAppAuthMiddleware()

	// Endpoints
	// Health check endpoint
	r.GET(defines.EndpointPing, healthCheck)

	// Authorized endpoints
	authorized := r.Group("/")
	authorized.Use(authMw.Check)

	appAuthorizer := r.Group("/")
	appAuthorizer.Use(appAuthMw.Check)

	// Users
	appAuthorizer.POST(defines.EndpointUsersCreate, usersCtrl.Create)
	appAuthorizer.POST(defines.EndpointUserAuthWithTelegram, usersCtrl.AuthWithTelegram)

	// Transactions
	authorized.POST(defines.EndpointTransactionsCreate, txnCtrl.Create)
	authorized.PUT(defines.EndpointTransactionsUpdateByMsgID, txnCtrl.UpdateByMsgID)
	authorized.GET(defines.EndpointTransactionsGetAll, txnCtrl.GetAll)

	// Wallets
	authorized.POST(defines.EndpointWalletsCreate, walletsCtrl.Create)
	authorized.GET(defines.EndpointWalletsGetAll, walletsCtrl.GetAll)
	authorized.GET(defines.EndpointWalletsGetByID, walletsCtrl.GetByID)
	authorized.PUT(defines.EndpointWalletsUpdateByID, walletsCtrl.UpdateByID)
	authorized.DELETE(defines.EndpointWalletsDeleteByID, walletsCtrl.DeleteByID)

	// Categories
	authorized.POST(defines.EndpointCategoriesCreate, catCtrl.Create)
	authorized.GET(defines.EndpointCategoriesGetAll, catCtrl.GetAll)
	authorized.DELETE(defines.EndpointCategoriesDeleteByID, catCtrl.DeleteByID)
	authorized.PUT(defines.EndpointCategoriesUpdateByID, catCtrl.UpdateByID)
	authorized.GET(defines.EndpointCategoriesGetByID, catCtrl.GetByID)
}

func healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, jsend.NewSuccess("pong"))
}
