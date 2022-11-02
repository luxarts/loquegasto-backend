package controller

import (
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/service"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/luxarts/jsend-go"

	"github.com/gin-gonic/gin"
)

type AuthorizerController interface {
	Register(ctx *gin.Context)
}

type authorizerController struct {
	oAuthSrv   service.OAuthService
	usersSrv   service.UsersService
	walletsSrv service.WalletsService
}

func NewAuthorizerController(oAuthSrv service.OAuthService, usersSrv service.UsersService, walletsSrv service.WalletsService) AuthorizerController {
	return &authorizerController{
		oAuthSrv:   oAuthSrv,
		usersSrv:   usersSrv,
		walletsSrv: walletsSrv,
	}
}
func (ctrl *authorizerController) Register(ctx *gin.Context) {
	userIDstr := ctx.Param(defines.ParamUserID)
	code := ctx.Query(defines.QueryCode)

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		return
	}

	// Exchange code for token
	token, err := ctrl.oAuthSrv.GetToken(userIDstr, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, jsend.NewError("failed to GetToken", err))
		return
	}

	now := time.Now()

	userDTO := &domain.UserDTO{
		ID:           userID,
		CreatedAt:    &now,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       &token.Expiry,
	}

	// Create user with token and create the default wallet
	userDTO, err = ctrl.usersSrv.Create(userDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, jsend.NewError("failed to create user", err))
		return
	}

	// Create default wallet
	// Create spreadsheet
	// Store spreadsheetID in DB

	// Redirect to Telegram
	location := url.URL{
		Scheme: "https",
		Host:   "t.me",
		Path:   "LoQueGastoTestBot",
	}

	ctx.Redirect(http.StatusFound, location.String())
}
