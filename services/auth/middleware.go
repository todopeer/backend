package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/todopeer/backend/orm"
)

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{"user"}

// AuthMiddleware middleware function to check for authenticated user
func AuthMiddleware(userORM *orm.UserORM) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Token from the Authorization header
			tokenStr := r.Header.Get("Authorization")

			// Token should start with "Bearer "
			if len(tokenStr) > 7 && strings.ToUpper(tokenStr[0:7]) == "BEARER " {
				tokenStr = tokenStr[7:]
			} else {
				log.Println("tokenString not provided: ", tokenStr)
				// treat it as no auth
				next.ServeHTTP(w, r)
				return
			}

			// If there is an error, the token must have been expired or malformed

			var user *orm.User
			defer func() {
				if user == nil {
					// auth provided but no such user -- means invalid
					next.ServeHTTP(w, r)
				} else {
					ctx := context.WithValue(r.Context(), userCtxKey, user)
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}()

			// Parse the token
			claim, err := tokenToClaim(tokenStr)
			if err == nil {
				err = claim.Valid()
			}
			if err != nil {
				// TODO: render HTTP error
				return
			}
			userid, sessionid, err := claimToUserInfo(claim)
			if err != nil {
				return
			}
			user, err = userORM.GetUserByID(userid)
			if err != nil || user == nil {
				return
			}

			if int32(sessionid) < user.SessionID {
				user = nil
				return
			}

		})
	}
}

// UserFromContext finds the user from the context. Returns nil if a user isn't present.
func UserFromContext(ctx context.Context) *orm.User {
	raw, _ := ctx.Value(userCtxKey).(*orm.User)
	return raw
}
