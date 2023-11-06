package service

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"log"
	"loquegasto-backend/internal/repository"
	"math/rand"
	"sync"
	"time"
)

type OAuthService interface {
	GetLoginURL(userID int64) (string, error)
	Callback(code string, state string) error
}

type oAuthService struct {
	statesVerifiers      map[string]string
	statesVerifiersMutex sync.Mutex
	repo                 repository.OAuthRepository
}

func NewOAuthService(repo repository.OAuthRepository) OAuthService {
	return &oAuthService{
		repo:            repo,
		statesVerifiers: make(map[string]string),
	}
}

func (svc *oAuthService) GetLoginURL(userID int64) (string, error) {
	state := uuid.NewString()
	codeVerifier, codeChallenge, err := svc.generatePKCE()
	if err != nil {
		return "", err
	}

	svc.statesVerifiersMutex.Lock()
	svc.statesVerifiers[state] = codeVerifier
	svc.statesVerifiersMutex.Unlock()

	log.Printf("code_verifier: %s\n", codeVerifier)

	loginURL := svc.repo.GetLoginURL(state, codeChallenge)

	return loginURL, nil
}

func (svc *oAuthService) Callback(code string, state string) error {
	svc.statesVerifiersMutex.Lock()
	codeVerifier, stateExists := svc.statesVerifiers[state]
	delete(svc.statesVerifiers, state)
	svc.statesVerifiersMutex.Unlock()

	if !stateExists {
		return errors.New("state not found")
	}
	token, err := svc.repo.Exchange(code, codeVerifier)
	if err != nil {
		return err
	}

	// Store token for user
	log.Printf("token: %+v\n", token)

	return nil
}

func (svc *oAuthService) generatePKCE() (verifier string, challenge string, err error) {
	cv := make([]byte, 32)
	_, err = rand.New(rand.NewSource(time.Now().UnixNano())).Read(cv)
	if err != nil {
		return "", "", err
	}

	verifier = base64.RawURLEncoding.EncodeToString(cv)

	cvhash := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(cvhash[:])

	return
}
