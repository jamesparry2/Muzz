package core

import (
	"context"
	"errors"

	"github.com/jamesparry2/Muzz/app/store"
)

var (
	ErrDiscoveryRequestNotValid = errors.New("request provided is invalid")
)

type DiscoveryRequest struct {
	UserID uint
}

type DiscoveryResponse struct {
	Users []store.User
}

func (c *Client) Discovery(ctx context.Context, request *DiscoveryRequest) (*DiscoveryResponse, error) {
	if request == nil {
		return nil, ErrDiscoveryRequestNotValid
	}

	user := store.User{}
	searchConditions := map[string]interface{}{"ID": request.UserID}
	if err := c.store.FindUser(ctx, &user, searchConditions); err != nil {
		return nil, err
	}

	users, err := c.store.FindAllUsers(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &DiscoveryResponse{
		Users: users,
	}, nil
}
