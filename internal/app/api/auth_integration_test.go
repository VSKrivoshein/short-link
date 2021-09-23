package api_test

import (
	"fmt"
	"github.com/VSKrivoshein/short-link/internal/app/api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// SignUP
func TestSignUp(t *testing.T) {
	tests := []struct {
		name             string
		email            string
		password         string
		wantCode         int
		isRequiredDelete bool
	}{
		{
			name:             "correct credentials",
			email:            "test@mail.ru",
			password:         "qwerty",
			wantCode:         http.StatusOK,
			isRequiredDelete: true,
		},
		{
			name:             "invalid email",
			email:            "test.ru",
			password:         "qwerty",
			wantCode:         http.StatusUnprocessableEntity,
			isRequiredDelete: false,
		},
		{
			name:             "invalid password (less than 6 char)",
			email:            "test@mail.ru",
			password:         "qwert",
			wantCode:         http.StatusUnprocessableEntity,
			isRequiredDelete: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := SignUp(t, test.email, test.password)
			defer RespClose(t, resp)

			assert.Equal(t, resp.StatusCode, test.wantCode, RespToString(t, resp))

			if test.isRequiredDelete {
				respDel := DeleteUser(t, test.email, test.password)
				defer RespClose(t, respDel)
				if respDel.StatusCode != http.StatusOK {
					t.Fatalf(
						"TestSignUp, %v: Fatal DeleteUser: %v",
						test.name,
						RespToString(t, respDel))
				}
			}
		})
	}
}

// Check that you can't create user with exist email
func TestSignUpExistedUser(t *testing.T) {
	var (
		email    = "test@email.com"
		password = "qwerty"
	)

	// Perform test
	resp := SignUp(t, email, password)
	defer RespClose(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Creating first user with unique email")
	respSecond := SignUp(t, email, password)
	defer RespClose(t, respSecond)
	assert.Equal(t, respSecond.StatusCode, http.StatusConflict, "Creating second user with non unique email")

	// Clean DB
	respDel := DeleteUser(t, email, password)
	defer RespClose(t, respDel)
	assert.Equal(t, respDel.StatusCode, http.StatusOK)
}

// SignIn
func TestSignIn(t *testing.T) {
	// Prepare for table test
	var (
		email    = "test@email.com"
		password = "qwerty"
	)
	respSetUp := SignUp(t, email, password)
	defer RespClose(t, respSetUp)
	assert.Equal(t, respSetUp.StatusCode, http.StatusOK, "Creating first user with unique email")

	// rm user
	defer func() {
		respDel := DeleteUser(t, email, password)
		defer RespClose(t, respDel)
		assert.Equal(t, respDel.StatusCode, http.StatusOK)
	}()

	// Test of SignIn
	tests := []struct {
		name      string
		email     string
		password  string
		wantCode  int
		isSuccess bool
	}{
		{
			name:     "Correct credentials",
			email:    email,
			password: password,
			wantCode: http.StatusOK,
		},
		{
			name:     "Incorrect Email",
			email:    "testWrong@mail.com",
			password: password,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "Incorrect password",
			email:    email,
			password: "shlapaPassword",
			wantCode: http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := SignIn(t, test.email, test.password)
			defer RespClose(t, resp)
			assert.Equal(t, resp.StatusCode, test.wantCode)

			if test.isSuccess {
				assert.NotEqual(t, resp.Header.Get(api.ShortenerCookieName), "")
			}
		})
	}
}

func TestSignOut(t *testing.T) {
	// Create user for test
	var (
		email    = "test@email.com"
		password = "qwerty"
	)
	respSetUp := SignUp(t, email, password)
	defer RespClose(t, respSetUp)
	assert.Equal(t, respSetUp.StatusCode, http.StatusOK, "Creating first user with unique email")

	// rm user
	defer func() {
		respDel := DeleteUser(t, email, password)
		defer RespClose(t, respDel)
		assert.Equal(t, respDel.StatusCode, http.StatusOK)
	}()

	// Test
	url := fmt.Sprintf("%v/auth/sign-out", TestSrv.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Fatatl http.NewRequest: %v", err)
	}

	req.Header.Set(api.ShortenerCookieName, "mockJWT")

	resp, err := TestSrv.Client().Do(req)
	defer RespClose(t, resp)
	if err != nil {
		t.Fatalf("Fatal TestSrv.Client().Do(): %v", err)
	}

	jwt := resp.Header.Get(api.ShortenerCookieName)

	assert.Equal(t, jwt, "")
}

func TestDeleteUser(t *testing.T) {
	// Create user for test
	var (
		email    = "test@email.com"
		password = "qwerty"
	)
	respSetUp := SignUp(t, email, password)
	defer RespClose(t, respSetUp)
	assert.Equal(t, respSetUp.StatusCode, http.StatusOK, "Creating first user with unique email")

	// rm user
	defer func() {
		// Required only in case of fatal one of the test
		respDel := DeleteUser(t, email, password)
		defer RespClose(t, respDel)
	}()

	tests := []struct {
		name     string
		email    string
		password string
		wantCode int
	}{
		{
			name:     "Incorrect email",
			email:    "wrong@mail.com",
			password: password,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "Incorrect password",
			email:    email,
			password: "incorrectPassword",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "Delete user",
			email:    email,
			password: password,
			wantCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := DeleteUser(t, test.email, test.password)
			assert.Equal(t, resp.StatusCode, test.wantCode)
			RespClose(t, resp)
		})
	}
}
