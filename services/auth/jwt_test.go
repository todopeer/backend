package auth

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/todopeer/backend/orm"
)

func TestGetTokenFromUser(t *testing.T) {
	// Setting up a user
	user := &orm.User{
		ID:        1,
		SessionID: 1,
	}

	// Call the function to be tested
	token, err := GetTokenFromUser(user)

	// Assert no error
	assert.NoError(t, err)

	// Assert the token is not empty
	assert.NotEmpty(t, token)
}

func TestTokenToClaim(t *testing.T) {
	// Create a claim for a user
	claim := &jwt.StandardClaims{
		Subject:   "1",
		Id:        "1",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(jwtExpireDuration).Unix(),
		Issuer:    "diarier",
	}

	// Get the token from the claim
	tokenStr, err := getToken(claim)
	assert.NoError(t, err)

	// Parse the token back to a claim
	claimTest, err := tokenToClaim(tokenStr)

	// Assert no error
	assert.NoError(t, err)

	// Assert the claimTest is equal to original claim
	assert.Equal(t, claim, claimTest)
}

func TestClaimToUserInfo(t *testing.T) {
	// Create a claim for a user
	claim := &jwt.StandardClaims{
		Subject: "1",
		Id:      "1",
	}

	// Call the function to be tested
	userid, sessionid, err := claimToUserInfo(claim)

	// Assert no error
	assert.NoError(t, err)

	// Assert the returned userid and sessionid are as expected
	assert.Equal(t, int64(1), userid)
	assert.Equal(t, int32(1), sessionid)
}
