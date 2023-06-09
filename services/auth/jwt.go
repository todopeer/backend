package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/todopeer/backend/orm"
)

var (
	jwtKey = []byte("to-be-replaced")
)

// SetJWTKey would be called during setup
func SetJWTKey(key string) {
	jwtKey = []byte(key)
}

const (
	jwtExpireDuration = 180 * 24 * time.Hour // Tokens will expire after half a year
	jwtIssuer         = "todopeer.com"
)

func tokenToClaim(tokenStr string) (*jwt.StandardClaims, error) {
	r := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, r, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claim, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claim: %s", token.Claims.Valid().Error())
	}
	return claim, nil
}

func getToken(claim *jwt.StandardClaims) (string, error) {

	now := time.Now()
	claim.IssuedAt = now.Unix()
	claim.ExpiresAt = now.Add(jwtExpireDuration).Unix()
	claim.Issuer = "diarier"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(jwtKey)
}

func GetTokenFromUser(user *orm.User) (string, error) {
	// Create the Claims
	claims := &jwt.StandardClaims{
		Subject: fmt.Sprintf("%d", user.ID),
		Id:      fmt.Sprintf("%d", user.SessionID),
	}

	return getToken(claims)
}

func claimToUserInfo(claim *jwt.StandardClaims) (userid int64, sessionid int32, err error) {
	userid, err = strconv.ParseInt(claim.Subject, 10, 64)

	if userid == 0 || err != nil {
		return 0, 0, errors.New("invalid userid")
	}

	session64, err := strconv.ParseInt(claim.Id, 10, 32)
	if err != nil {
		return 0, 0, errors.New("invalid sessionid")
	}

	return userid, int32(session64), nil
}
