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

type WalletsController interface {
	Create(ctx *gin.Context)
	GetByName(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	UpdateByID(ctx *gin.Context)
	DeleteByID(ctx *gin.Context)
}
type walletsController struct {
	srv service.WalletsService
}

func NewWalletsController(srv service.WalletsService) WalletsController {
	return &walletsController{
		srv: srv,
	}
}
func (c *walletsController) Create(ctx *gin.Context) {
	var body domain.WalletDTO

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid body"))
		return
	}

	body.UserID = ctx.GetInt(defines.ParamUserID)

	if !body.IsValid() {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid body"))
		return
	}

	response, err := c.srv.Create(&body)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusCreated, jsend.NewSuccess(response))
}
func (c *walletsController) GetByName(ctx *gin.Context) {
	userID := ctx.GetInt(defines.ParamUserID)
	name := ctx.Query(defines.ParamName)

	if name == "" {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("missing name query param"))
		return
	}

	response, err := c.srv.GetByName(userID, name)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *walletsController) GetByID(ctx *gin.Context) {
	userID := ctx.GetInt(defines.ParamUserID)
	idStr := ctx.Param(defines.ParamID)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid id"))
		return
	}

	response, err := c.srv.GetByID(userID, id)
	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *walletsController) GetAll(ctx *gin.Context) {
	userID := ctx.GetInt(defines.ParamUserID)
	filterName := ctx.Query(defines.ParamName)

	var response *[]domain.WalletDTO
	var err error

	if filterName != "" {
		response, err = c.srv.GetByName(userID, filterName)
	} else {
		response, err = c.srv.GetAll(userID)
	}

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *walletsController) UpdateByID(ctx *gin.Context) {
	var body domain.WalletDTO

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid body"))
		return
	}
	idStr := ctx.Param(defines.ParamID)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid id"))
		return
	}
	body.ID = id
	body.UserID = ctx.GetInt(defines.ParamUserID)

	response, err := c.srv.UpdateByID(&body)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *walletsController) DeleteByID(ctx *gin.Context) {
	idStr := ctx.Param(defines.ParamID)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid id"))
		return
	}
	userID := ctx.GetInt(defines.ParamUserID)

	err = c.srv.DeleteByID(id, userID)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusNoContent, jsend.NewSuccess(nil))
}
