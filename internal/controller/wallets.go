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

	response, err := c.srv.GetByID(userID, id)
	var jsendErr *jsend.Body
	if errors.As(err, &jsendErr) && jsendErr != nil {
		ctx.JSON(*jsendErr.Code, jsendErr)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *walletsController) GetAll(ctx *gin.Context) {
	userID := ctx.GetString(middleware.CtxKeyUserID)

	var response interface{}
	var err error

	if name, exists := ctx.GetQuery(defines.ParamName); exists {
		response, err = c.srv.GetByName(userID, name)
	} else {
		response, err = c.srv.GetAll(userID)
	}

	var jsendErr *jsend.Body
	if errors.As(err, &jsendErr) && jsendErr != nil {
		ctx.JSON(*jsendErr.Code, jsendErr)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *walletsController) UpdateByID(ctx *gin.Context) {
	var body domain.WalletCreateRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}
	//id := ctx.Param(defines.ParamID)

	userID := ctx.GetString(middleware.CtxKeyUserID)

	if !body.IsValid() {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	// Check if name already used
	wallet, err := c.srv.GetByName(userID, body.Name)
	if err, isError := err.(*jsend.Body); isError && err != nil && *err.Code != http.StatusNotFound {
		ctx.JSON(*err.Code, err)
		return
	}
	if wallet != nil {
		ctx.JSON(http.StatusConflict, defines.ErrNameAlreadyExists)
		return
	}

	response, err := c.srv.UpdateByID(&body)

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
