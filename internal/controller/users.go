package controller

import (
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luxarts/jsend-go"
)

type UsersController interface {
	Create(ctx *gin.Context)
	Get(ctx *gin.Context)
}
type userController struct {
	srv service.UsersService
}

func NewUsersController(srv service.UsersService) UsersController {
	return &userController{
		srv: srv,
	}
}
func (c *userController) Create(ctx *gin.Context) {
	var body domain.UserDTO

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid body"))
		return
	}

	body.ID = ctx.GetInt(defines.ParamUserID)

	if !body.IsValid() {
		ctx.JSON(http.StatusBadRequest, jsend.NewFail("invalid body"))
		return
	}

	user, err := c.srv.GetByID(body.ID)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	if user != nil {
		ctx.JSON(http.StatusConflict, jsend.NewFail("user already exists"))
		return
	}

	response, err := c.srv.Create(&body)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusCreated, jsend.NewSuccess(response))
}
func (c *userController) Get(ctx *gin.Context) {
	userID := ctx.GetInt(defines.ParamUserID)

	response, err := c.srv.GetByID(userID)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
