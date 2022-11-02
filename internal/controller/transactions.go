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
	GetAll(ctx *gin.Context)
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
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	if !body.IsValid() {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	body.UserID = ctx.GetInt64(defines.ParamUserID)

	response, err := c.srv.Create(&body)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusCreated, jsend.NewSuccess(response))
}
func (c *transactionsController) UpdateByMsgID(ctx *gin.Context) {
	msgIDStr := ctx.Param(defines.ParamMsgID)
	msgID, err := strconv.Atoi(msgIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidMsgID)
		return
	}

	var body domain.TransactionDTO
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}
	if !body.IsValidForUpdate() || body.MsgID != msgID {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	body.UserID = ctx.GetInt64(defines.ParamUserID)

	response, err := c.srv.UpdateByMsgID(&body)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *transactionsController) GetAll(ctx *gin.Context) {
	userID := ctx.GetInt64(defines.ParamUserID)

	filters := make(domain.TransactionFilters)

	walletID, _ := ctx.GetQuery(defines.QueryWalletID)
	if walletID != "" {
		filters[defines.QueryWalletID] = walletID
	}
	categoryID, _ := ctx.GetQuery(defines.QueryCategoryID)
	if categoryID != "" {
		filters[defines.QueryCategoryID] = categoryID
	}

	response, err := c.srv.GetAll(userID, &filters)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
