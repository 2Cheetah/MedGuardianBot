package domain

type UserRepository interface {
	CreateUser(user *User) error
	GetUser(id uint64) (*User, error)
}
