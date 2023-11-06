package repository

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"os"
)

type OAuthRepository interface {
	GetLoginURL(state string, codeChallenge string) string
	Exchange(code string, codeVerifier string) (*oauth2.Token, error)
}

type oAuthRepository struct {
	config *oauth2.Config
}

func NewOAuthRepository() OAuthRepository {
	/*credentials, err := os.ReadFile("googlecredentials.json")
	if err != nil {
		log.Fatalf("Error reading credentials:\n%+v\n", err)
	}*/
	credentials := os.Getenv("GOOGLE_CREDENTIALS")
	config, err := google.ConfigFromJSON([]byte(credentials), "https://www.googleapis.com/auth/drive.file")
	if err != nil {
		log.Fatalf("Error google ConfigFromJSON: %+v\n", err)
	}

	return &oAuthRepository{
		config: config,
	}
}

func (repo *oAuthRepository) GetLoginURL(state string, codeChallenge string) string {
	return repo.config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
}
func (repo *oAuthRepository) Exchange(code string, codeVerifier string) (*oauth2.Token, error) {
	ctx := context.Background()
	return repo.config.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
}
