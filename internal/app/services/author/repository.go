package author

type Repository interface {
	GetUser(user *User) error
	CreateUser(user *User) error
	DeleteUser(user *User) error
}
