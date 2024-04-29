package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/luxarts/jsend-go"
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/middleware"
	"loquegasto-backend/internal/service"
	"net/http"
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
	var body domain.WalletCreateRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	if !body.IsValid() {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	userID := ctx.GetString(middleware.CtxKeyUserID)

	response, err := c.srv.Create(&body, userID)

	var jsendErr *jsend.Body
	if errors.As(err, &jsendErr) && jsendErr != nil {
		ctx.JSON(*jsendErr.Code, jsendErr)
		return
	}

	ctx.JSON(http.StatusCreated, jsend.NewSuccess(response))
}
func (c *walletsController) GetByID(ctx *gin.Context) {
	userID := ctx.GetString(middleware.CtxKeyUserID)
	id := ctx.Param(defines.ParamID)

	response, err := c.srv.GetByID(id, userID)
	var jsendErr *jsend.Body
	if errors.As(err, &jsendErr) && jsendErr != nil {
		ctx.JSON(*jsendErr.Code, jsendErr)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *walletsController) GetAll(ctx *gin.Context) {
	userID := ctx.GetString(middleware.CtxKeyUserID)

	response, err := c.srv.GetAll(userID)

	var jsendErr *jsend.Body
	if errors.As(err, &jsendErr) && jsendErr != nil {
		ctx.JSON(*jsendErr.Code, jsendErr)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *walletsController) UpdateByID(ctx *gin.Context) {
	var body domain.WalletUpdateRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	if !body.IsValid() {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	id := ctx.Param(defines.ParamID)
	userID := ctx.GetString(middleware.CtxKeyUserID)

	response, err := c.srv.UpdateByID(&body, id, userID)

	var jsendErr *jsend.Body
	if errors.As(err, &jsendErr) && jsendErr != nil {
		ctx.JSON(*jsendErr.Code, jsendErr)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *walletsController) DeleteByID(ctx *gin.Context) {
	id := ctx.Param(defines.ParamID)

	userID := ctx.GetString(middleware.CtxKeyUserID)

	err := c.srv.DeleteByID(id, userID)

	var jsendErr *jsend.Body
	if errors.As(err, &jsendErr) && jsendErr != nil {
		ctx.JSON(*jsendErr.Code, jsendErr)
		return
	}

	ctx.JSON(http.StatusNoContent, jsend.NewSuccess(nil))
}
