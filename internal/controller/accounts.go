package controller

import (
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luxarts/jsend-go"
)

type AccountsController interface {
	Create(ctx *gin.Context)
	GetByName(ctx *gin.Context)
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
