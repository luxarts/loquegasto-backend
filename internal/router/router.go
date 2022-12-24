package router

import (
	"fmt"
	"loquegasto-backend/internal/controller"
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/middleware"
	"loquegasto-backend/internal/repository"
	"loquegasto-backend/internal/service"
	"loquegasto-backend/internal/utils/jwt"
	"net/http"
	"os"
	"strconv"

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
	catRepo := repository.NewCategoriesRepository(db)

	// Services init
	txnSrv := service.NewTransactionsService(txnRepo, walletsRepo, catRepo)
	usersSrv := service.NewUsersService(usersRepo)
	walletsSrv := service.NewWalletsService(walletsRepo)
	catSrv := service.NewCategoriesService(catRepo)

	// Controllers init
	txnCtrl := controller.NewTransactionsController(txnSrv)
	usersCtrl := controller.NewUsersController(usersSrv)
	walletsCtrl := controller.NewWalletsController(walletsSrv)
	catCtrl := controller.NewCategoriesController(catSrv)

	// Middleware
	authMw := middleware.NewAuthMiddleware()

	// Endpoints
	r.GET("/token/:userID", generateToken)

	// Health check endpoint
	r.GET(defines.EndpointPing, healthCheck)

	// Authorized endpoints
	authorized := r.Group("/")
	authorized.Use(authMw.Check)

	// Transactions
	authorized.POST(defines.EndpointTransactionsCreate, txnCtrl.Create)
	authorized.PUT(defines.EndpointTransactionsUpdateByMsgID, txnCtrl.UpdateByMsgID)
	authorized.GET(defines.EndpointTransactionsGetAll, txnCtrl.GetAll)

	// Users
	authorized.POST(defines.EndpointUsersCreate, usersCtrl.Create)
	authorized.GET(defines.EndpointUsersGet, usersCtrl.Get)
	authorized.PUT(defines.EndpointUsersUpdate, usersCtrl.Update)
	authorized.DELETE(defines.EndpointUsersDelete, usersCtrl.Delete)

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
}

func healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, jsend.NewSuccess("pong"))
}
func generateToken(ctx *gin.Context) {
	userID := ctx.Param("userID")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidID)
		return
	}
	token := jwt.GenerateToken(nil, &jwt.Payload{Subject: userIDInt})

	ctx.String(http.StatusOK, token)
}
