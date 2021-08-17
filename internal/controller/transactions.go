package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/luxarts/jsend-go"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/service"
	"net/http"
)

type TransactionsController interface {
	Create(ctx *gin.Context)
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

	if !body.IsValid(){
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid body"))
		return
	}

	response, err := c.srv.Create(&body)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, jsend.NewSuccess(response))
}
