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

type AccountsController interface {
	Create(ctx *gin.Context)
	GetByName(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	UpdateByID(ctx *gin.Context)
}
type accountsController struct {
	srv service.AccountsService
}

func NewAccountsController(srv service.AccountsService) AccountsController {
	return &accountsController{
		srv: srv,
	}
}
func (c *accountsController) Create(ctx *gin.Context) {
	var body domain.AccountDTO

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, jsend.NewError("shouldbindjson-error", err))
		return
	}

	body.UserID = ctx.GetInt(defines.ParamUserID)

	if !body.IsValid() {
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
func (c *accountsController) GetByName(ctx *gin.Context) {
	userID := ctx.GetInt(defines.ParamUserID)
	name := ctx.Query(defines.ParamName)

	response, err := c.srv.GetByName(userID, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *accountsController) GetByID(ctx *gin.Context) {
	userID := ctx.GetInt(defines.ParamUserID)
	idStr := ctx.Param(defines.ParamID)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, jsend.NewError("atoi-error", err))
		return
	}

	response, err := c.srv.GetByID(userID, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *accountsController) UpdateByID(ctx *gin.Context) {
	var body domain.AccountDTO

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, jsend.NewError("shouldbindjson-error", err))
		return
	}
	idStr := ctx.Param(defines.ParamID)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, jsend.NewError("atoi-error", err))
		return
	}
	body.ID = id
	body.UserID = ctx.GetInt(defines.ParamUserID)

	response, err := c.srv.UpdateByID(&body)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
