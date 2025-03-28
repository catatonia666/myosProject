package apiserver

import (
	"dialogue/internal/models"
	"dialogue/internal/store/teststore"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Home(t *testing.T) {
	os.Chdir("../../..")
	newCash, _ := newCash("redis://:@localhost:6379/0")
	s := newServer(teststore.New(), newCash)

	testCases := []struct {
		name         string
		method       string
		expectedCode int
	}{
		{
			name:         "valid",
			method:       "GET",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid method POST",
			method:       "POST",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "invalid method DELETE",
			method:       "DELETE",
			expectedCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, "/home", nil)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestHandler_UserSignup(t *testing.T) {
	os.Chdir("../../..")
	store := teststore.New()
	s := newServer(store, nil)

	existingUser := &models.User{
		Nickname: "example",
		Email:    "emailinuse@example.org",
		Password: "password",
	}
	store.User().Create(existingUser)

	testCases := []struct {
		name         string
		payload      url.Values
		expectedCode int
	}{
		{
			name: "valid",
			payload: url.Values{
				"nickname": {"example"},
				"email":    {"user@example.org"},
				"password": {"password"},
			},
			expectedCode: http.StatusFound,
		},
		{
			name: "no nickname provided",
			payload: url.Values{
				"email":    {"user@example.org"},
				"password": {"password"},
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "no email provided",
			payload: url.Values{
				"nickname": {"example"},
				"password": {"password"},
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "no password provided",
			payload: url.Values{
				"nickname": {"example"},
				"email":    {"user@example.org"},
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "password is too short",
			payload: url.Values{
				"nickname": {"example"},
				"email":    {"user@example.org"},
				"password": {"short"},
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "duplicate email",
			payload: url.Values{
				"nickname": {"example"},
				"email":    {"emailinuse@example.org"},
				"password": {"password"},
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "incorrect email",
			payload: url.Values{
				"nickname": {"example"},
				"email":    {"email"},
				"password": {"password"},
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			body := strings.NewReader(tc.payload.Encode())

			req, _ := http.NewRequest(http.MethodPost, "/user/signup", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)

			fmt.Println(s.store.User().FindByEmail(""))
		})
	}
}
