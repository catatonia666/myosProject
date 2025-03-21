package models

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Nickname: "sample",
		Email:    "user@example.org",
		Password: "password",
	}
}
