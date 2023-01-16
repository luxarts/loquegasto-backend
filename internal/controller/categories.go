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
	UpdateByID(ctx *gin.Context)
	GetByID(ctx *gin.Context)
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
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	body.UserID = ctx.GetInt64(defines.ParamUserID)

	if !body.IsValid() {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	category, err := c.srv.GetByName(body.Name, body.UserID)
	if err, isError := err.(*jsend.Body); isError && err != nil && *err.Code != http.StatusNotFound {
		ctx.JSON(*err.Code, err)
		return
	}
	if category != nil {
		ctx.JSON(http.StatusConflict, defines.ErrNameAlreadyExists)
		return
	}

	category, err = c.srv.GetByEmoji(body.Emoji, body.UserID)
	if err, isError := err.(*jsend.Body); isError && err != nil && *err.Code != http.StatusNotFound {
		ctx.JSON(*err.Code, err)
		return
	}
	if category != nil {
		ctx.JSON(http.StatusConflict, defines.ErrEmojiAlreadyExists)
		return
	}

	response, err := c.srv.Create(&body)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusCreated, jsend.NewSuccess(response))
}
func (c *categoriesController) GetAll(ctx *gin.Context) {
	userID := ctx.GetInt64(defines.ParamUserID)

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
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidID)
		return
	}
	userID := ctx.GetInt64(defines.ParamUserID)

	err = c.srv.DeleteByID(id, userID)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusNoContent, jsend.NewSuccess(nil))
}
func (c *categoriesController) UpdateByID(ctx *gin.Context) {
	var body domain.CategoryDTO

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}
	idStr := ctx.Param(defines.ParamID)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidID)
		return
	}
	body.ID = id
	body.UserID = ctx.GetInt64(defines.ParamUserID)

	if !body.IsValid() {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidBody)
		return
	}

	// Check if name already used
	category, err := c.srv.GetByName(body.Name, body.UserID)
	if err, isError := err.(*jsend.Body); isError && err != nil && *err.Code != http.StatusNotFound {
		ctx.JSON(*err.Code, err)
		return
	}
	if category != nil && category.ID != body.ID {
		ctx.JSON(http.StatusConflict, defines.ErrNameAlreadyExists)
		return
	}

	// Check if emoji already used
	category, err = c.srv.GetByEmoji(body.Emoji, body.UserID)
	if err, isError := err.(*jsend.Body); isError && err != nil && *err.Code != http.StatusNotFound {
		ctx.JSON(*err.Code, err)
		return
	}
	if category != nil && category.ID != body.ID {
		ctx.JSON(http.StatusConflict, defines.ErrEmojiAlreadyExists)
		return
	}

	response, err := c.srv.UpdateByID(&body)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
func (c *categoriesController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param(defines.ParamID)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, defines.ErrInvalidID)
		return
	}
	userID := ctx.GetInt64(defines.ParamUserID)

	response, err := c.srv.GetByID(id, userID)

	if err, isError := err.(*jsend.Body); isError && err != nil {
		ctx.JSON(*err.Code, err)
		return
	}

	ctx.JSON(http.StatusOK, jsend.NewSuccess(response))
}
