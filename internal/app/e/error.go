package e

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// sort err

var (
	ErrUnexpected = errors.New("unexpected error")

	ErrCredentialsInvalid = errors.New("login or password is incorrect")
	ErrToken              = errors.New("token error")

	ErrCreatingUser         = errors.New("error during creating user, email is already exist")
	ErrGetUser              = errors.New("email or password incorrect")
	ErrDeletingUser         = errors.New("error during deleting user")
	ErrDeletingUserNotFound = errors.New("error during deleting user, user was not found")

	ErrCreatingLink       = errors.New("error link is already exist")
	ErrDeleteLink         = errors.New("error during deleting link")
	ErrDeleteLinkNotFound = errors.New("error link was not found")
	ErrGetLink            = errors.New("error can not find link")
	ErrGetAllLinks        = errors.New("error during getting all links")

	ErrUserUnauthorized    = errors.New("user unauthorized")
	ErrUnprocessableEntity = errors.New("data structure is not correct")

	ErrUserIDWasNotFound = errors.New("user was not found in create handler")
)

type CustomErrorWithCode struct {
	Code int
	Msg  error
}

func New(originalError, forUserError error, httpCode int) *CustomErrorWithCode {
	return &CustomErrorWithCode{
		Msg:  fmt.Errorf("original error msg: %v, for user error message: %w", originalError, forUserError),
		Code: httpCode,
	}
}

func (e *CustomErrorWithCode) Error() string {
	return e.Msg.Error()
}

func (e *CustomErrorWithCode) ErrorForUser() string {
	return UnwrapRecursive(e.Msg).Error()
}

func GetInfo() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}

	fn := runtime.FuncForPC(pc).Name()
	funcName := fn[strings.LastIndex(fn, ".")+1:]

	return fmt.Sprintf("%v: %%w", funcName)
}

func UnwrapRecursive(err error) error {
	unwrappedError := errors.Unwrap(err)
	if unwrappedError != nil {
		return UnwrapRecursive(unwrappedError)
	}
	return err
}
