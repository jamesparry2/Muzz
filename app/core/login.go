package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/jamesparry2/Muzz/app/store"
)

var (
	ErrLoginMissingRequest = errors.New("login was passed an invalid request")
	ErrLoginInvalidCreds   = errors.New("provided username/password was incorrect")
	// Standard error to ensure no critial errors leak out for potential attacks
	ErrLoginStandardError = errors.New("login is currently not doing the login")
)

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token string
}

func (c *Client) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	if request == nil {
		return nil, ErrLoginMissingRequest
	}

	user := store.User{}
	searchConditions := map[string]interface{}{"email": request.Email}
	if err := c.store.FindUser(ctx, &user, searchConditions); err != nil {
		return nil, ErrLoginInvalidCreds
	}

	encryptedRequestPassword, err := c.auth.EncryptPassword(request.Password)
	if err != nil {
		return nil, ErrLoginStandardError
	}

	if user.Password != encryptedRequestPassword {
		return nil, ErrLoginInvalidCreds
	}

	token, err := c.auth.CreateToken(fmt.Sprint(user.ID), "user")
	if err != nil {
		return nil, ErrLoginStandardError
	}

	return &LoginResponse{Token: token}, nil
}
