package middleware

import (
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/utils/jwt"
	"net/http"
	"strings"

	"github.com/luxarts/jsend-go"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	Check(c *gin.Context)
}

type authMiddleware struct {
}

func NewAuthMiddleware() AuthMiddleware {
	return &authMiddleware{}
}

func (a *authMiddleware) Check(c *gin.Context) {
	bearerToken := c.GetHeader("Authorization")

	if bearerToken == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsend.NewFail("missing Authorization token"))
		return
	}

	bearerTokenSplit := strings.Split(bearerToken, " ")

	if len(bearerTokenSplit) != 2 || bearerTokenSplit[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsend.NewFail("invalid Authorization token"))
		return
	}

	tokenSplit := strings.Split(bearerTokenSplit[1], ".")

	if len(tokenSplit) != 3 {
		c.AbortWithStatusJSON(http.StatusBadRequest, jsend.NewFail("invalid Authorization token"))
		return
	}

	if !jwt.Verify(tokenSplit[0], tokenSplit[1], tokenSplit[2]) {
		c.AbortWithStatusJSON(http.StatusNotFound, jsend.NewFail("invalid signature"))
		return
	}

	// Set userID in context
	userID, err := jwt.GetSubject(bearerToken)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, jsend.NewError("getsubject-error", err))
		return
	}

	c.Set(defines.ParamUserID, userID)
}
