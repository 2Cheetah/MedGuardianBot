package domain

type UserRepository interface {
	CreateUser(user *User) error
	GetUser(id int64) (*User, error)
}
