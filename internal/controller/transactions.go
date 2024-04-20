package controller

import (
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/middleware"
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

	body.UserID = ctx.GetString(middleware.CtxKeyUserID)

	response, err := c.srv.Create(&body)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusCreated, jsend.NewSuccess(response))
}
func (c *transactionsController) UpdateByMsgID(ctx *gin.Context) {
	msgIDStr := ctx.Param(defines.ParamMsgID)
	msgID, err := strconv.ParseInt(msgIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidMsgID)
		return
	}

	var body domain.TransactionDTO
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}
	if !body.IsValidForUpdate() {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	body.UserID = ctx.GetString(middleware.CtxKeyUserID)
	body.MsgID = msgID

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
	from, _ := ctx.GetQuery(defines.QueryFrom)
	if from != "" {
		filters[defines.QueryFrom] = from
	}
	to, _ := ctx.GetQuery(defines.QueryTo)
	if to != "" {
		filters[defines.QueryTo] = to
	}

	response, err := c.srv.GetAll(userID, &filters)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
