package middleware

import (
	"github.com/gin-gonic/gin"
	"loquegasto-backend/internal/defines"
	"net/http"
	"os"
)

type AppAuthMiddleware interface {
	Check(ctx *gin.Context)
}
type appAuthMiddleware struct {
}

func NewAppAuthMiddleware() AppAuthMiddleware {
	return &appAuthMiddleware{}
}

func (m *appAuthMiddleware) Check(ctx *gin.Context) {
	clientID := ctx.GetHeader("client_id")
	clientSecret := ctx.GetHeader("client_secret")

	if clientID != os.Getenv("CLIENT_ID") || clientSecret != os.Getenv("CLIENT_SECRET") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, defines.ErrUnauthorized)
		return
	}
}
