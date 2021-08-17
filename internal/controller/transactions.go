package controller

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/service"
	"loquegasto-backend/internal/utils/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luxarts/jsend-go"
)

type TransactionsController interface {
	Create(ctx *gin.Context)
	GetTotal(ctx *gin.Context)
}

type transactionsController struct {
	srv service.TransactionsService
}

func NewTransactionsController(srv service.TransactionsService) TransactionsController {
	return &transactionsController{
		srv: srv,
	}
}

func (c *transactionsController) Create(ctx *gin.Context) {
	var body domain.TransactionDTO

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, jsend.NewError("shouldbindjson-error", err))
		return
	}

	if !body.IsValid() {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid body"))
		return
	}

	// Get userID from token
	bearerToken := ctx.GetHeader("Authorization")
	userID, err := jwt.GetSubject(bearerToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	body.UserID = userID

	response, err := c.srv.Create(&body)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, jsend.NewSuccess(response))
}
func (c *transactionsController) GetTotal(ctx *gin.Context) {
	// Get userID from token
	bearerToken := ctx.GetHeader("Authorization")
	userID, err := jwt.GetSubject(bearerToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	response, err := c.srv.GetTotal(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
