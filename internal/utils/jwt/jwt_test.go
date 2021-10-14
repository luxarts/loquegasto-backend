package jwt

import (
	"loquegasto-backend/internal/defines"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var prevJwtSecret string

func Test_GenerateToken(t *testing.T) {
	setupTest()
	header := Header{
		Algorithm: "HS256",
		Type:      "JWT",
	}
	payload := Payload{
		Subject: "Cosme Fulanito",
	}

	token := GenerateToken(&header, &payload)

	require.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJDb3NtZSBGdWxhbml0byJ9.lxYGUGjD0Lgb0OHseRwlhXfhHczcQJB7DVX3o1eVOOo", token)
	teardownTest()
}
func Test_generateSignature(t *testing.T) {
	setupTest()
	header := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9`
	payload := `eyJzdWIiOiJDb3NtZSBGdWxhbml0byJ9`

	signature := generateSignature(header, payload)

	require.Equal(t, "lxYGUGjD0Lgb0OHseRwlhXfhHczcQJB7DVX3o1eVOOo", signature)
	teardownTest()
}

func Test_VerifyOK(t *testing.T) {
	setupTest()
	jwt := strings.Split("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJDb3NtZSBGdWxhbml0byJ9.lxYGUGjD0Lgb0OHseRwlhXfhHczcQJB7DVX3o1eVOOo", ".")

	out := Verify(jwt[0], jwt[1], jwt[2])

	require.True(t, out)
	teardownTest()
}

func Test_VerifyFail(t *testing.T) {
	setupTest()
	jwt := strings.Split("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJDb3NtZSBGdWxhbml0byJ9.lxYGUGjD0Lgb0OHseRwlhXfhHczcQJB7DVX3o1eVOO1", ".")

	out := Verify(jwt[0], jwt[1], jwt[2])

	require.False(t, out)
	teardownTest()
}

func setupTest() {
	prevJwtSecret = os.Getenv(defines.EnvJWTSecret)
	os.Setenv(defines.EnvJWTSecret, "dGVzdGluZ2tleTEyMzQ")
}
func teardownTest() {
	os.Setenv(defines.EnvJWTSecret, prevJwtSecret)
	prevJwtSecret = ""
}
