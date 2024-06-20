package core

import (
	"context"

	"github.com/jamesparry2/Muzz/app/auth"
	"github.com/jamesparry2/Muzz/app/store"
)

type CoreIface interface {
	CreateUser(ctx context.Context, request *CreateUserRequest) (*CreateUserResponse, error)
	Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error)
	Swipe(ctx context.Context, request *SwipeRequest) (*SwipeResponse, error)
	Discovery(ctx context.Context, request *DiscoveryRequest) (*DiscoveryResponse, error)
	Preference(ctx context.Context, request *PreferenceRequest) error
}

type ClientOptions struct {
	Store store.StoreIface
	Auth  auth.AuthIface
}

type Client struct {
	store store.StoreIface
	auth  auth.AuthIface
}

func NewClient(opts *ClientOptions) *Client {
	return &Client{
		store: opts.Store,
		auth:  opts.Auth,
	}
}
