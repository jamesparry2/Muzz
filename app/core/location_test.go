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

func TestLocation(t *testing.T) {
	t.Run("should return an error when the request provided is ommited", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{})

		err := client.Location(context.Background(), nil)

		assert.ErrorIs(t, err, core.ErrLocationMissingRequest, "unexpected error was returned")
	})

	t.Run("should return an error when the user provided couldn't be found", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return store.ErrFindUserMissingUserDetails
				},
			},
		})

		err := client.Location(context.Background(), &core.LocationRequest{Lat: 1.0, Long: 2.0})

		assert.ErrorIs(t, err, store.ErrFindUserMissingUserDetails, "unexpected error was returned")
	})

	t.Run("should return an error when trying to update user for the location", func(t *testing.T) {
		errUpsertLocation := errors.New("failed to upsert")

		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return nil
				},
				MockUpsertUser: func(ctx context.Context, user *store.User) error {
					return errUpsertLocation
				},
			},
		})

		err := client.Location(context.Background(), &core.LocationRequest{UserID: 1})

		assert.ErrorIs(t, err, errUpsertLocation, "unexpected error was returned")
	})

	t.Run("should return a success when no errors occured", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return nil
				},
				MockUpsertUser: func(ctx context.Context, user *store.User) error {
					return nil
				},
			},
		})

		err := client.Location(context.Background(), &core.LocationRequest{UserID: 1})

		assert.Nil(t, err, "no error should have been returned")
	})
}
