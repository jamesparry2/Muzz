package core

import (
	"context"
	"errors"

	"github.com/jamesparry2/Muzz/app/store"
)

var (
	ErrPreferenceMissingRequest = errors.New("preference request provided was invalid")
)

type PreferenceRequest struct {
	Gender string
	MaxAge int
	MinAge int

	UserID uint
}

func (c *Client) Preference(ctx context.Context, request *PreferenceRequest) error {
	if request == nil {
		return ErrPreferenceMissingRequest
	}

	user := store.User{}
	searchConditions := map[string]interface{}{"ID": request.UserID}
	if err := c.store.FindUser(ctx, &user, searchConditions); err != nil {
		return err
	}

	user.Preferences = &store.Preferences{
		Gender:     request.Gender,
		MinimumAge: request.MinAge,
		MaximumAge: request.MaxAge,
	}

	if err := c.store.UpsertUser(ctx, &user); err != nil {
		return err
	}

	return nil
}
