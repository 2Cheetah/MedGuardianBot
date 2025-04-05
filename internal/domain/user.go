package domain

import "errors"

type User struct {
	FirstName string
	LastName  string
	ID        int64
	Username  string
}

func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username field can't be empty")
	}

	if u.ID == 0 {
		return errors.New("id can't be zero")
	}

	return nil
}
