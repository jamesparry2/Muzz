package mock

import (
	"context"

	"github.com/jamesparry2/Muzz/app/store"
)

type MockStore struct {
	MockUpsertUser  func(ctx context.Context, request *store.User) error
	MockFindUser    func(ctx context.Context, user *store.User, conditions map[string]interface{}) error
	MockUpsertSwipe func(ctx context.Context, swipe *store.Swipe) error
	MockHasMatched  func(ctx context.Context, userID, matchedID uint) (bool, error)
	MockFindAllUser func(ctx context.Context, user *store.User) ([]store.User, error)
}

func (ms *MockStore) UpsertUser(ctx context.Context, request *store.User) error {
	return ms.MockUpsertUser(ctx, request)
}

func (ms *MockStore) FindUser(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
	return ms.MockFindUser(ctx, user, conditions)
}

func (ms *MockStore) UpsertSwipe(ctx context.Context, swipe *store.Swipe) error {
	return ms.MockUpsertSwipe(ctx, swipe)
}

func (ms *MockStore) HasMatched(ctx context.Context, userId, matchedId uint) (bool, error) {
	return ms.MockHasMatched(ctx, userId, matchedId)
}

func (ms *MockStore) FindAllUsers(ctx context.Context, user *store.User) ([]store.User, error) {
	return ms.MockFindAllUser(ctx, user)
}
