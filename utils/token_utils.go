// Package utils provides a collection of utility functions, including helpers for handling JWT tokens and API keys.
// This file contains utilities for creating, parsing, and managing JWT tokens and encrypted API keys.
// It uses the "github.com/golang-jwt/jwt/v5" library for JWT functionality and AES for API key encryption.
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

// tokenApiSecret is the secret key used for AES encryption of API keys. It must be 32 bytes long.
var tokenApiSecret []byte

// tokenSigningKey is the key used for signing and verifying JWT tokens.
var tokenSigningKey []byte

// region Initialize secrets -------------------------------------------------------------------------------------------

// SetSecret initializes the secret keys used for token and API key operations.
// This function must be called before using any of the token or API key functions.
//
// Parameters:
//
//	secret: A 32-byte secret key for AES encryption.
//	signingKey: A byte slice used for signing JWT tokens.
//
// Returns:
//
//	An error if the secret or vector lengths are not valid.
func SetSecret(secret, signingKey []byte) error {
	if len(secret) != 32 {
		return fmt.Errorf("secret key must be 32 bytes long")
	}
	tokenApiSecret = secret
	tokenSigningKey = signingKey
	return nil
}

// endregion

// region Singleton Pattern --------------------------------------------------------------------------------------------

// TokenUtilsStruct provides methods for token and API key manipulation.
// It is used as a singleton to provide a centralized and consistent way to handle tokens.
type TokenUtilsStruct struct{}

var (
	doOnceForTokenUtils sync.Once
	tokenUtilsSingleton *TokenUtilsStruct
)

// TokenUtils returns a singleton instance of TokenUtilsStruct.
// This ensures that all token operations are handled through a single, consistent interface.
func TokenUtils() *TokenUtilsStruct {
	doOnceForTokenUtils.Do(func() {
		tokenUtilsSingleton = &TokenUtilsStruct{}
	})
	return tokenUtilsSingleton
}

// endregion

// region Access Token parsing helpers ---------------------------------------------------------------------------------

// CreateToken generates a new JWT token with the provided claims.
//
// Example:
//
//	claims := &jwt.RegisteredClaims{
//	    ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
//	    Issuer:    "my-app",
//	}
//	token, err := TokenUtils().CreateToken(claims)
//
// Parameters:
//
//	claims: The JWT claims to be included in the token.
//
// Returns:
//
//	A signed JWT token string.
//	An error if token signing fails.
func (t *TokenUtilsStruct) CreateToken(claims *jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tokenSigningKey)
}

// ParseToken parses a JWT token string and returns the claims.
// It also verifies the token's signature.
//
// Example:
//
//	claims, err := TokenUtils().ParseToken(tokenString)
//	if err != nil {
//	    // Handle error
//	}
//
// Parameters:
//
//	tokenString: The JWT token string to parse.
//
// Returns:
//
//	The parsed JWT claims.
//	An error if parsing or validation fails.
func (t *TokenUtilsStruct) ParseToken(tokenString string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tokenSigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}

// endregion

// region API Key parsing helpers --------------------------------------------------------------------------------------

// CreateApiKey generates an encrypted API key from an application name.
// The encryption is done using AES.
//
// Parameters:
//
//	appName: The application name to encrypt.
//
// Returns:
//
//	An encrypted, hex-encoded API key.
//	An error if encryption fails.
func (t *TokenUtilsStruct) CreateApiKey(appName string) (string, error) {
	return t.encrypt(appName)
}

// ParseApiKey decrypts an API key to retrieve the original application name.
//
// Parameters:
//
//	apiKey: The encrypted, hex-encoded API key.
//
// Returns:
//
//	The decrypted application name.
//	An error if decryption fails.
func (t *TokenUtilsStruct) ParseApiKey(apiKey string) (string, error) {
	return t.decrypt(apiKey)
}

// endregion

// region PRIVATE SECTION ----------------------------------------------------------------------------------------------

// encrypt encrypts a string using AES-CFB and returns it as a hex-encoded string.
func (t *TokenUtilsStruct) encrypt(value string) (string, error) {
	if len(tokenApiSecret) == 0 {
		return "", fmt.Errorf("token API secret is not set")
	}

	block, err := aes.NewCipher(tokenApiSecret)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(value))
	iv := cipherText[:aes.BlockSize]
	if _, er := io.ReadFull(rand.Reader, iv); er != nil {
		return "", er
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(value))

	return hex.EncodeToString(cipherText), nil
}

// decrypt decrypts a hex-encoded, AES-CFB encrypted string.
func (t *TokenUtilsStruct) decrypt(value string) (string, error) {
	if len(tokenApiSecret) == 0 {
		return "", fmt.Errorf("token API secret is not set")
	}

	cipherText, err := hex.DecodeString(value)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(tokenApiSecret)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

// endregion
