package controller

import (
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luxarts/jsend-go"
)

type TransactionsController interface {
	Create(ctx *gin.Context)
	UpdateByMsgID(ctx *gin.Context)
	GetAllByUserID(ctx *gin.Context)
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
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid body"))
		return
	}

	if !body.IsValid() {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid body"))
		return
	}

	body.UserID = ctx.GetInt(defines.ParamUserID)

	response, err := c.srv.Create(&body)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusCreated, jsend.NewSuccess(response))
}
func (c *transactionsController) UpdateByMsgID(ctx *gin.Context) {
	userID := ctx.GetInt(defines.ParamUserID)
	msgIDStr := ctx.Param(defines.ParamMsgID)
	msgID, err := strconv.Atoi(msgIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid message id"))
		return
	}

	var body domain.TransactionDTO
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid body"))
		return
	}
	if !body.IsValidForUpdate() {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid body"))
		return
	}
	if body.MsgID != msgID {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid body: msg_id doesn't match"))
		return
	}

	response, err := c.srv.UpdateByMsgID(userID, &body)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *transactionsController) GetAllByUserID(ctx *gin.Context) {
	userID := ctx.GetInt(defines.ParamUserID)

	response, err := c.srv.GetAllByUserID(userID)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
