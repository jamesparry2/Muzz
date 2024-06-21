package core

import (
	"context"
	"errors"

	"github.com/jamesparry2/Muzz/app/store"
)

var (
	ErrLocationMissingRequest = errors.New("location request provided was invalid")
)

type LocationRequest struct {
	Lat            float64
	Long           float64
	DistanceFromMe int

	UserID uint
}

func (c *Client) Location(ctx context.Context, request *LocationRequest) error {
	if request == nil {
		return ErrLocationMissingRequest
	}

	user := store.User{}
	searchConditions := map[string]interface{}{"ID": request.UserID}
	if err := c.store.FindUser(ctx, &user, searchConditions); err != nil {
		return err
	}

	user.Location = &store.Location{
		Lat:  request.Lat,
		Long: request.Long,
	}

	if err := c.store.UpsertUser(ctx, &user); err != nil {
		return err
	}

	return nil
}
