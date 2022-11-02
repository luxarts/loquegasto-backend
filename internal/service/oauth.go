package service

import (
	"context"
	"log"
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/domain"
	"net/url"
	"os"

	"golang.org/x/oauth2"

	"golang.org/x/oauth2/google"
)

type OAuthService interface {
	GetToken(userID string, code string) (*oauth2.Token, error)
}

type oAuthService struct {
	creds domain.OAuthCredentials
}

func NewOAuthService() OAuthService {
	s := oAuthService{}

	clientID := os.Getenv(defines.EnvGoogleClientID)
	projectID := os.Getenv(defines.EnvGoogleProjectID)
	secret := os.Getenv(defines.EnvGoogleClientSecret)

	s.creds = domain.NewCredentials(clientID, projectID, secret)

	return &s
}

func (s *oAuthService) GetToken(userID string, code string) (*oauth2.Token, error) {
	redirectURI := url.URL{
		Scheme: "http",
		Host:   os.Getenv(defines.EnvBackendBaseURL),
		Path:   defines.APIAuthURL + userID,
	}

	credsBytes := s.creds.AddRedirectURI(redirectURI.String()).Bytes()

	config, err := google.ConfigFromJSON(credsBytes, defines.ScopeSheetsRW)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return config.Exchange(context.Background(), code)
}
