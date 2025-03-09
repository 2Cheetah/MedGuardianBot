package domain

type UserRepository interface {
	CreateUser(user *User) error
}
