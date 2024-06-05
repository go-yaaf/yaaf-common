// Token Utils tests

package test

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v5"

	"github.com/go-yaaf/yaaf-common/utils"
)

const (
	tokenApiSecret = "thisTokenSecretKeyMustBe32Length"
	tokenApiVector = "thisIsTokenInitializingVectorN32"
)

func TestApiKey(t *testing.T) {
	skipCI(t)
	appName := "my-application"

	if err := utils.SetSecret([]byte(tokenApiSecret), []byte(tokenApiVector)); err != nil {
		t.Fail()
	}

	tu := utils.TokenUtils()
	apiKey, err := tu.CreateApiKey(appName)
	if err != nil {
		t.Fail()
	}

	fmt.Println(apiKey)

	if name, err := tu.ParseApiKey(apiKey); err != nil {
		t.Fail()
	} else {
		fmt.Println(name)
	}
}

func TestJwtToken(t *testing.T) {
	skipCI(t)
	// Create the Registered Claims and map Token Data fields
	rc := &jwt.RegisteredClaims{
		ID:       "claimId",
		Issuer:   "accountId",
		Subject:  "subjectId",
		Audience: []string{},
	}

	if err := utils.SetSecret([]byte(tokenApiSecret), []byte(tokenApiVector)); err != nil {
		t.Fail()
	}

	tu := utils.TokenUtils()
	token, err := tu.CreateToken(rc)
	if err != nil {
		t.Fail()
	}

	fmt.Println(token)

	if claims, err := tu.ParseToken(token); err != nil {
		t.Fail()
	} else {
		fmt.Println(claims.ID)
	}
}
