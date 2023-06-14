package resolver

import (
	"github.com/todopeer/backend/orm"
	"github.com/todopeer/backend/services/auth"
)

func generateToken(user *orm.User) (string, error) {
	return auth.GetTokenFromUser(user)
}
