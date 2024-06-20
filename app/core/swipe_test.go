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

func TestSwipe(t *testing.T) {
	t.Run("should return an error when now swipe request is provided", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{})

		response, err := client.Swipe(context.Background(), nil)

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, core.ErrSwipeMissingRequest, "unexpected error was returned")
	})

	t.Run("should return an error when the user provided couldn't be found", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return store.ErrFindUserMissingUserDetails
				},
			},
		})

		response, err := client.Swipe(context.Background(), &core.SwipeRequest{})

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, store.ErrFindUserMissingUserDetails, "unexpected error was returned")
	})

	t.Run("should return an error when failing to update/insert a swipe record", func(t *testing.T) {
		upsertErr := errors.New("failed to save upsert")
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return nil
				},
				MockUpsertSwipe: func(ctx context.Context, swipe *store.Swipe) error {
					return upsertErr
				},
			},
		})

		response, err := client.Swipe(context.Background(), &core.SwipeRequest{})

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, upsertErr, "unexpected error was returned")
	})

	t.Run("should return a success with a non match when the swipe was non-desirable", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return nil
				},
				MockUpsertSwipe: func(ctx context.Context, swipe *store.Swipe) error {
					return nil
				},
			},
		})

		response, err := client.Swipe(context.Background(), &core.SwipeRequest{
			IsDesired: "NO",
		})

		assert.Nil(t, err, "no err should be returned in success flow")
		assert.Equal(t, false, response.HasMatched, "unexpected matched status was returned")
	})

	t.Run("should return an error when trying to evaluate if the users have matched", func(t *testing.T) {
		hasMatchedErr := errors.New("failed to match")
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return nil
				},
				MockUpsertSwipe: func(ctx context.Context, swipe *store.Swipe) error {
					return nil
				},
				MockHasMatched: func(ctx context.Context, userID, matchedID uint) (bool, error) {
					return false, hasMatchedErr
				},
			},
		})

		response, err := client.Swipe(context.Background(), &core.SwipeRequest{
			IsDesired: "YES",
		})

		assert.Nil(t, response, "no response should be returned in error flow")
		assert.ErrorIs(t, err, hasMatchedErr, "unexpected error was returned")
	})

	t.Run("should return a success and a matched ID if both users have swipped desirably for each other", func(t *testing.T) {
		client := core.NewClient(&core.ClientOptions{
			Store: &store_mock.MockStore{
				MockFindUser: func(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
					return nil
				},
				MockUpsertSwipe: func(ctx context.Context, swipe *store.Swipe) error {
					return nil
				},
				MockHasMatched: func(ctx context.Context, userID, matchedID uint) (bool, error) {
					return true, nil
				},
			},
		})

		response, err := client.Swipe(context.Background(), &core.SwipeRequest{
			IsDesired: "YES",
			MatchedID: 2,
		})

		assert.Nil(t, err, "no err should be returned in success flow")
		assert.Equal(t, true, response.HasMatched, "unexpected matched status was returned")
		assert.Equal(t, uint(2), response.MatchedID, "unexpected match id was returned")
	})
}
