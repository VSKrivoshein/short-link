package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *User
		isValid bool
	}{
		{
			name:    "valid",
			u:       func() *User { return TestUser(t) },
			isValid: true,
		},
		{
			name: "empty email",
			u: func() *User {
				u := TestUser(t)
				u.Email = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid email",
			u: func() *User {
				u := TestUser(t)
				u.Email = "invalidEmail"
				return u
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *User {
				u := TestUser(t)
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "short password",
			u: func() *User {
				u := TestUser(t)
				u.Email = "1234"
				return u
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	u := TestUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}
