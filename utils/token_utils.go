// JWT token utilities
//

package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key to encode API keys (must be 32 characters)
var tokenApiSecret []byte

// Initializing vector to encode API keys (must be 32 characters)
var tokenSigningKey []byte

// region Initialize secrets -------------------------------------------------------------------------------------------

// SetSecret set the secret key and initializing vector to encode/decode API keys
func SetSecret(secret, vector []byte) error {
	if len(secret) != 32 {
		return fmt.Errorf("secret must be 32 bytes length")
	}

	if len(vector) != 32 {
		return fmt.Errorf("vector must be 32 bytes length")
	}

	tokenApiSecret = secret
	tokenSigningKey = vector
	return nil
}

// endregion

// region Singleton Pattern --------------------------------------------------------------------------------------------

type TokenUtilsStruct struct {
}

var doOnceForTokenUtils sync.Once

var tokenUtilsSingleton *TokenUtilsStruct = nil

// TokenUtils is a factory method that acts as a static member
func TokenUtils() *TokenUtilsStruct {
	doOnceForTokenUtils.Do(func() {
		tokenUtilsSingleton = &TokenUtilsStruct{}
	})
	return tokenUtilsSingleton
}

// endregion

// region Access Token parsing helpers ---------------------------------------------------------------------------------

// CreateToken build JWT token from Token Data structure
func (t *TokenUtilsStruct) CreateToken(claims *jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tokenSigningKey)
}

// ParseToken rebuild Token Data structure from JWT token
func (t *TokenUtilsStruct) ParseToken(tokenString string) (*jwt.RegisteredClaims, error) {

	rc := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, rc, func(token *jwt.Token) (interface{}, error) {
		return tokenSigningKey, nil
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
func (t *TokenUtilsStruct) CreateApiKey(appName string) (string, error) {
	return t.encrypt(appName)
}

// ParseApiKey extract application name from API key
func (t *TokenUtilsStruct) ParseApiKey(apiKey string) (string, error) {
	return t.decrypt(apiKey)
}

// endregion

// region PRIVATE SECTION ----------------------------------------------------------------------------------------------

// encrypt string using AES and return base64
func (t *TokenUtilsStruct) encrypt(value string) (string, error) {

	block, err := aes.NewCipher(tokenApiSecret)
	if err != nil {
		return "", err
	}

	// Generate a new random IV
	cipherText := make([]byte, aes.BlockSize+len(value))
	iv := cipherText[:aes.BlockSize]
	if _, er := io.ReadFull(rand.Reader, iv); er != nil {
		return "", er
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(value))

	return hex.EncodeToString(cipherText), nil
}

// decrypt base64 string using AES
func (t *TokenUtilsStruct) decrypt(value string) (string, error) {
	cipherTextBytes, err := hex.DecodeString(value)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(tokenApiSecret)
	if err != nil {
		return "", err
	}

	if len(cipherTextBytes) < aes.BlockSize {
		return "", fmt.Errorf("cipher text too short")
	}

	iv := cipherTextBytes[:aes.BlockSize]
	cipherTextBytes = cipherTextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherTextBytes, cipherTextBytes)

	return string(cipherTextBytes), nil
}

// endregion
