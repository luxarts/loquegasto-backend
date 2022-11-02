package controller

import (
	"fmt"
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/service"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type AuthorizerController interface {
	Login(ctx *gin.Context)
}

type authorizerController struct {
	oAuthSrv service.OAuthService
}

func NewAuthorizerController(oAuthSrv service.OAuthService) AuthorizerController {
	return &authorizerController{oAuthSrv: oAuthSrv}
}
func (ctrl *authorizerController) Login(ctx *gin.Context) {
	userID := ctx.Param(defines.ParamUserID)
	code := ctx.Query(defines.QueryCode)

	// Exchange code for token
	token, err := ctrl.oAuthSrv.GetToken(userID, code)
	if err != nil {
		return
	}

	// Create user with token and create the default wallet
	// Create default wallet
	// Create spreadsheet
	// Store spreadsheetID in DB

	fmt.Printf("UserID: %s Token: %+v\n", userID, token)

	// Redirect to Telegram
	location := url.URL{
		Scheme: "https",
		Host:   "t.me",
		Path:   "LoQueGastoTestBot",
	}

	ctx.Redirect(http.StatusFound, location.String())
}
