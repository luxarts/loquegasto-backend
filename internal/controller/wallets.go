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
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	body.UserID = ctx.GetInt(defines.ParamUserID)

	if !body.IsValid() {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	// Check if name already used
	wallet, err := c.srv.GetByName(body.UserID, body.Name)
	if err, isError := err.(*jsend.Body); isError && err != nil && *err.Code != http.StatusNotFound {
		ctx.JSON(*err.Code, err)
		return
	}
	if wallet != nil {
		ctx.JSON(http.StatusConflict, defines.ErrNameAlreadyExists)
		return
	}

	response, err := c.srv.Create(&body)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusCreated, jsend.NewSuccess(response))
}
func (c *walletsController) GetByID(ctx *gin.Context) {
	userID := ctx.GetInt(defines.ParamUserID)
	idStr := ctx.Param(defines.ParamID)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidID)
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

	var response interface{}
	var err error

	if name, exists := ctx.GetQuery(defines.ParamName); exists {
		response, err = c.srv.GetByName(userID, name)
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
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}
	idStr := ctx.Param(defines.ParamID)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidID)
		return
	}
	body.ID = id
	body.UserID = ctx.GetInt(defines.ParamUserID)

	if !body.IsValid() {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	// Check if name already used
	wallet, err := c.srv.GetByName(body.UserID, body.Name)
	if err, isError := err.(*jsend.Body); isError && err != nil && *err.Code != http.StatusNotFound {
		ctx.JSON(*err.Code, err)
		return
	}
	if wallet != nil {
		ctx.JSON(http.StatusConflict, defines.ErrNameAlreadyExists)
		return
	}

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
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidID)
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
