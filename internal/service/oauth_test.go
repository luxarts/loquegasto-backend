package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePKCE(t *testing.T) {
	// Given
	svc := oAuthService{}

	// When
	verifier, challenge, err := svc.generatePKCE()

	// Then
	assert.Equal(t, 43, len(verifier))
	assert.NotEqual(t, "", challenge)
	assert.Nil(t, err)
}
