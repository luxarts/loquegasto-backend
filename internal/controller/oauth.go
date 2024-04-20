package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/service"
	"net/http"
)

type OAuthController interface {
	GetLoginURL(ctx *gin.Context)
	Callback(ctx *gin.Context)
}

type oAuthController struct {
	svc service.OAuthService
}

func NewOAuthController(svc service.OAuthService) OAuthController {
	return &oAuthController{svc: svc}
}

func (ctrl *oAuthController) GetLoginURL(ctx *gin.Context) {
	userID := ctx.GetInt64(defines.ParamUserID)

	loginURL, err := ctrl.svc.GetLoginURL(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"url": loginURL})
}

func (ctrl *oAuthController) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	if code == "" || state == "" {
		log.Println("Missing params code or state")
		ctx.Redirect(http.StatusFound, "https://t.me/LoQueGastoTestBot")
		return
	}

	err := ctrl.svc.Callback(code, state)
	if err != nil {
		log.Printf("Error callback: %+v\n", err)
		ctx.Redirect(http.StatusFound, "https://t.me/LoQueGastoTestBot")
		return
	}

	ctx.Redirect(http.StatusFound, "https://t.me/LoQueGastoTestBot")
}
