package author

type Service interface {
	SingIn(user *User) error
	SingUp(user *User) error
	CheckAuthAndRefresh(user *User) error
	DeleteUser(user *User) error
}
