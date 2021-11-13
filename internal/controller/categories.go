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

type CategoriesController interface {
	Create(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	DeleteByID(ctx *gin.Context)
}
type categoriesController struct {
	srv service.CategoriesService
}

func NewCategoriesController(srv service.CategoriesService) CategoriesController {
	return &categoriesController{
		srv: srv,
	}
}

func (c *categoriesController) Create(ctx *gin.Context) {
	var body domain.CategoryDTO

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
func (c *categoriesController) GetAll(ctx *gin.Context) {
	userID := ctx.GetInt(defines.ParamUserID)

	var response interface{}
	var err error

	if emoji, exists := ctx.GetQuery(defines.ParamEmoji); exists {
		response, err = c.srv.GetByEmoji(emoji, userID)
	} else if name, exists := ctx.GetQuery(defines.ParamName); exists {
		response, err = c.srv.GetByName(name, userID)
	} else {
		response, err = c.srv.GetAll(userID)
	}

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *categoriesController) DeleteByID(ctx *gin.Context) {
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
