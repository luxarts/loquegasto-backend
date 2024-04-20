package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/luxarts/jsend-go"
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/service"
	"net/http"
)

type UsersController interface {
	Create(ctx *gin.Context)
	AuthWithTelegram(ctx *gin.Context)
}
type userController struct {
	svc service.UsersService
}

func NewUsersController(svc service.UsersService) UsersController {
	return &userController{
		svc: svc,
	}
}
func (c *userController) Create(ctx *gin.Context) {
	var body domain.UserCreateRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	if !body.IsValid() {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	response, err := c.svc.Create(&body)

	var jsendErr *jsend.Body
	if errors.As(err, &jsendErr) && jsendErr != nil {
		ctx.JSON(*jsendErr.Code, jsendErr)
		return
	}

	ctx.JSON(http.StatusCreated, jsend.NewSuccess(response))
}
func (c *userController) AuthWithTelegram(ctx *gin.Context) {
	var body domain.UserAuthWithTelegramRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	if !body.IsValid() {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	response, err := c.svc.AuthWithTelegram(&body)

	var jsendErr *jsend.Body
	if errors.As(err, &jsendErr) && jsendErr != nil {
		ctx.JSON(*jsendErr.Code, jsendErr)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
