package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jamesparry2/Muzz/app/core"
	"github.com/jamesparry2/Muzz/app/store"
	store_mock "github.com/jamesparry2/Muzz/app/store/mock"
	"github.com/stretchr/testify/assert"
)

func TestDiscovery(t *testing.T) {
	t.Run("should return an error when the request provided is ommited", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{})

		response, err := client.Discovery(context.Background(), nil)

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, core.ErrDiscoveryRequestNotValid, "unexpected error was returned")
	})

	t.Run("should return an error when the user provided couldn't be found", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return store.ErrFindUserMissingUserDetails
				},
			},
		})

		response, err := client.Discovery(context.Background(), &core.DiscoveryRequest{UserID: 1})

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, store.ErrFindUserMissingUserDetails, "unexpected error was returned")
	})

	t.Run("should return an error when trying to find users to display for discovery", func(t *testing.T) {
		errNotUsers := errors.New("no users to be discovered")

		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return nil
				},
				MockFindAllUser: func(ctx context.Context, user *store.User) ([]store.User, error) {
					return []store.User{}, errNotUsers
				},
			},
		})

		response, err := client.Discovery(context.Background(), &core.DiscoveryRequest{UserID: 1})

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, errNotUsers, "unexpected error was returned")
	})

	t.Run("should return a success when no errors occured", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return nil
				},
				MockFindAllUser: func(ctx context.Context, user *store.User) ([]store.User, error) {
					return []store.User{{}, {}}, nil
				},
			},
		})

		response, err := client.Discovery(context.Background(), &core.DiscoveryRequest{UserID: 1})

		assert.Nil(t, err, "no err should be returned in success flow")
		assert.Len(t, response.Users, 2, "unexpected amount of users returned")
	})
}
