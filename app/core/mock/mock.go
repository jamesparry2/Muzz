package mock

import (
	"context"

	"github.com/jamesparry2/Muzz/app/core"
)

type MockCore struct {
	MockCreateUser func(ctx context.Context, request *core.CreateUserRequest) (*core.CreateUserResponse, error)
	MockLogin      func(ctx context.Context, request *core.LoginRequest) (*core.LoginResponse, error)
	MockSwipe      func(ctx context.Context, request *core.SwipeRequest) (*core.SwipeResponse, error)
	MockDiscovery  func(ctx context.Context, request *core.DiscoveryRequest) (*core.DiscoveryResponse, error)
	MockPreference func(ctx context.Context, request *core.PreferenceRequest) error
}

func (mc *MockCore) CreateUser(ctx context.Context, request *core.CreateUserRequest) (*core.CreateUserResponse, error) {
	return mc.MockCreateUser(ctx, request)
}

func (mc *MockCore) Login(ctx context.Context, request *core.LoginRequest) (*core.LoginResponse, error) {
	return mc.MockLogin(ctx, request)
}

func (mc *MockCore) Swipe(ctx context.Context, request *core.SwipeRequest) (*core.SwipeResponse, error) {
	return mc.MockSwipe(ctx, request)
}

func (mc *MockCore) Discovery(ctx context.Context, request *core.DiscoveryRequest) (*core.DiscoveryResponse, error) {
	return mc.MockDiscovery(ctx, request)
}

func (mc *MockCore) Preference(ctx context.Context, request *core.PreferenceRequest) error {
	return mc.MockPreference(ctx, request)
}
