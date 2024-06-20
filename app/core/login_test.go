package core_test

import (
	"context"
	"errors"
	"testing"

	auth_mock "github.com/jamesparry2/Muzz/app/auth/mock"
	"github.com/jamesparry2/Muzz/app/core"
	"github.com/jamesparry2/Muzz/app/store"
	store_mock "github.com/jamesparry2/Muzz/app/store/mock"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	t.Run("should return an error when no request is provided", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{})

		response, err := client.Login(context.Background(), nil)

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, core.ErrLoginMissingRequest, "unexpected error was returned")
	})

	t.Run("should return an error when the user provided couldn't be found", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return store.ErrFindUserMissingUserDetails
				},
			},
		})

		response, err := client.Login(context.Background(), &core.LoginRequest{})

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, core.ErrLoginInvalidCreds, "unexpected error was returned")
	})

	t.Run("should return an error when the encryption failed", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return nil
				},
			},
			Auth: &auth_mock.AuthMock{
				MockEncryptPassword: func(pwd string) (string, error) {
					return "", errors.New("failed to encrypt password")
				},
			},
		})

		response, err := client.Login(context.Background(), &core.LoginRequest{})

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, core.ErrLoginStandardError, "unexpected error was returned")
	})

	t.Run("should return an error when the password provided does not match", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return nil
				},
			},
			Auth: &auth_mock.AuthMock{
				MockEncryptPassword: func(pwd string) (string, error) {
					return "password", nil
				},
			},
		})

		response, err := client.Login(context.Background(), &core.LoginRequest{
			Password: "different_password",
		})

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, core.ErrLoginInvalidCreds, "unexpected error was returned")
	})

	t.Run("should return an error when the token creation has failed", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					user.Password = "password"
					return nil
				},
			},
			Auth: &auth_mock.AuthMock{
				MockEncryptPassword: func(pwd string) (string, error) {
					return "password", nil
				},
				MockCreateToken: func(subject, aud string) (string, error) {
					return "", errors.New("failed to gen token")
				},
			},
		})

		response, err := client.Login(context.Background(), &core.LoginRequest{
			Password: "password",
		})

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, core.ErrLoginStandardError, "unexpected error was returned")
	})

	t.Run("should return a valid token that can be used for authentication", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					user.Password = "password"
					return nil
				},
			},
			Auth: &auth_mock.AuthMock{
				MockEncryptPassword: func(pwd string) (string, error) {
					return "password", nil
				},
				MockCreateToken: func(subject, aud string) (string, error) {
					return "awesomeauthtoken", nil
				},
			},
		})

		response, err := client.Login(context.Background(), &core.LoginRequest{
			Password: "password",
		})

		assert.Nil(t, err, "no err should be returned in a success flow")
		assert.Equal(t, "awesomeauthtoken", response.Token, "unexpected token was returned")
	})
}
