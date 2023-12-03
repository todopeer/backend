package model

import (
	"database/sql"

	"github.com/todopeer/backend/orm"
)

func ConvertFromRegistrationInput(registrationInput *UserRegistrationInput) *orm.User {
	if registrationInput == nil {
		return nil
	}

	u := &orm.User{
		Name:  &registrationInput.Name,
		Email: registrationInput.Email,
	}

	if registrationInput.Username != nil {
		u.Username = &sql.NullString{
			String: *registrationInput.Username,
			Valid:  true,
		}
	} else {
		u.Username = &sql.NullString{
			Valid: false,
		}
	}

	u.SetPassword(registrationInput.Password)
	return u
}
