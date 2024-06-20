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

func TestCreateUser(t *testing.T) {
	t.Run("should return an error when the request provided is ommited", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{})

		response, err := client.CreateUser(context.Background(), nil)

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, core.ErrCreateUserMissingRequest, "unexpected error was returned")
	})

	t.Run("should return an error when the user being created already has the email provided in use", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return store.ErrFindUserMissingUserDetails
				},
			},
		})

		response, err := client.CreateUser(context.Background(), &core.CreateUserRequest{
			Email: "jamesparry2@gmail.com",
		})

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, store.ErrFindUserMissingUserDetails, "unexpected error was returned")
	})

	t.Run("should return an error when the ecryption process for the password has failed", func(t *testing.T) {
		pwdError := errors.New("unexpected error with password saving")

		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return nil
				},
			},
			Auth: &auth_mock.AuthMock{
				MockEncryptPassword: func(pwd string) (string, error) {
					return "", pwdError
				},
			},
		})

		response, err := client.CreateUser(context.Background(), &core.CreateUserRequest{
			Email: "jamesparry2@gmail.com",
		})

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, pwdError, "unexpected error was returned")
	})

	t.Run("should return an error when the user is failed to be saved", func(t *testing.T) {
		upsertErr := errors.New("failed to save the user")

		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return nil
				},
				MockUpsertUser: func(ctx context.Context, request *store.User) error {
					return upsertErr
				},
			},
			Auth: &auth_mock.AuthMock{
				MockEncryptPassword: func(pwd string) (string, error) {
					return "", nil
				},
			},
		})

		response, err := client.CreateUser(context.Background(), &core.CreateUserRequest{
			Email: "jamesparry2@gmail.com",
		})

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, upsertErr, "unexpected error was returned")
	})

	t.Run("should return a valid response when no errors are present", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return nil
				},
				MockUpsertUser: func(ctx context.Context, request *store.User) error {
					request.Age = 28
					request.Name = "James"

					return nil
				},
			},
			Auth: &auth_mock.AuthMock{
				MockEncryptPassword: func(pwd string) (string, error) {
					return "", nil
				},
			},
		})

		response, err := client.CreateUser(context.Background(), &core.CreateUserRequest{
			Email: "jamesparry2@gmail.com",
		})

		assert.Nil(t, err, "no err should be returned in success flow")
		assert.Equal(t, 28, response.Age, "unexpected age was returned")
		assert.Equal(t, "James", response.Name, "unexpected age was returned")
	})
}
