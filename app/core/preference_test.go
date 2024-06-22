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

func TestPreference(t *testing.T) {
	t.Run("should return an error when the request provided is ommited", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{})

		err := client.Preference(context.Background(), nil)

		assert.ErrorIs(t, err, core.ErrPreferenceMissingRequest, "unexpected error was returned")
	})

	t.Run("should return an error when the user provided couldn't be found", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return store.ErrFindUserMissingUserDetails
				},
			},
		})

		err := client.Preference(context.Background(), &core.PreferenceRequest{})

		assert.ErrorIs(t, err, store.ErrFindUserMissingUserDetails, "unexpected error was returned")
	})

	t.Run("should return an error when trying to update user for the preference", func(t *testing.T) {
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

		err := client.Preference(context.Background(), &core.PreferenceRequest{})

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

		err := client.Preference(context.Background(), &core.PreferenceRequest{})

		assert.Nil(t, err, "no error should have been returned")
	})
}
