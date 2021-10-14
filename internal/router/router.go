package router

import (
	"context"
	"log"
	"loquegasto-backend/internal/controller"
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/middleware"
	"loquegasto-backend/internal/repository"
	"loquegasto-backend/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luxarts/jsend-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	r.GET(defines.EndpointTransactionsGetTotal, authMw.Check, txnCtrl.GetTotal)
	r.PUT(defines.EndpointTransactionsUpdateByMsgID, authMw.Check, txnCtrl.UpdateByMsgID)

	// Health check endpoint
	r.GET(defines.EndpointPing, healthCheck)
}

func initMongoClient(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalln(err)
	}
	return client
}

func healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, jsend.NewSuccess("pong"))
}
