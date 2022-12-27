// Copyright 2022. Motty Cohen
//
// JWT token utilities
//

package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/golang-jwt/jwt/v4"
)

// Secret key to encode API keys (must be 32 characters)
var tokenApiSecret string

// Initializing vector to encode API keys (must be 16 characters)
var tokenApiVector string

// Secret key to sign JWT token (must be 32 characters)
var tokenSigningKey string

// region Initialize secrets -------------------------------------------------------------------------------------------

// SetAPISecret set the secret key and initializing vector to encode/decode API keys
func SetAPISecret(secret, vector string) error {
	if len(secret) < 32 {
		return fmt.Errorf("secret must be 32 characters length")
	}

	if len(vector) < 16 {
		return fmt.Errorf("vector must be 16 characters length")
	}

	tokenApiSecret = secret
	tokenApiVector = vector[0:16]
	return nil
}

// SetJWTSecret set the private key to sign JWT token
func SetJWTSecret(secret string) error {
	if len(secret) < 32 {
		return fmt.Errorf("secret must be 32 characters length")
	} else {
		tokenSigningKey = secret
		return nil
	}
}

// endregion

// region Singleton Pattern --------------------------------------------------------------------------------------------

type tokenUtils struct {
}

var doOnceForTokenUtils sync.Once

var tokenUtilsSingleton *tokenUtils = nil

// TokenUtils is a factory method that acts as a static member
func TokenUtils() *tokenUtils {
	doOnceForTokenUtils.Do(func() {
		tokenUtilsSingleton = &tokenUtils{}
	})
	return tokenUtilsSingleton
}

// endregion

// region Access Token parsing helpers ---------------------------------------------------------------------------------

// CreateToken build JWT token from Token Data structure
func (t *tokenUtils) CreateToken(claims *jwt.RegisteredClaims) (string, error) {
	signingKey := []byte(tokenSigningKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

// ParseToken rebuild Token Data structure from JWT token
func (t *tokenUtils) ParseToken(tokenString string) (*jwt.RegisteredClaims, error) {

	rc := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, rc, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSigningKey), nil
	})

	if err != nil {
		return nil, err
	} else {
		return rc, nil
	}
}

// endregion

// region API Key parsing helpers --------------------------------------------------------------------------------------

// CreateApiKey generate API Key from application name
func (t *tokenUtils) CreateApiKey(appName string) (string, error) {
	return encrypt(appName)
}

// ParseApiKey extract application name from API key
func (t *tokenUtils) ParseApiKey(apiKey string) (string, error) {
	return decrypt(apiKey)
}

// endregion

// region PRIVATE SECTION ----------------------------------------------------------------------------------------------

// encrypt string using AES and return base64
func encrypt(value string) (string, error) {

	if err := validateSecret(); err != nil {
		return "", err
	}

	text := []byte(value)
	block, err := aes.NewCipher([]byte(tokenApiSecret))
	if err != nil {
		return "", err
	}

	bytes := []byte(tokenApiVector)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(text))
	cfb.XORKeyStream(cipherText, text)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// decrypt base64 string using AES
func decrypt(value string) (string, error) {
	if err := validateSecret(); err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(value)

	block, err := aes.NewCipher([]byte(tokenApiSecret))
	if err != nil {
		return "", err
	}

	bytes := []byte(tokenApiVector)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

// validate secret and vector
func validateSecret() error {

	if len(tokenApiSecret) < 32 {
		return fmt.Errorf("secret must be 32 characters length")
	}

	if len(tokenApiVector) < 16 {
		return fmt.Errorf("vector must be 16 characters length")
	}

	return nil
}

// endregion
