package core

import (
	"context"
	"errors"

	"github.com/jamesparry2/Muzz/app/store"
)

var (
	ErrSwipeMissingRequest = errors.New("swipe request provided was invalid")
)

type SwipeRequest struct {
	UserID uint

	MatchedID uint
	IsDesired string
}

type SwipeResponse struct {
	HasMatched bool
	MatchedID  uint
}

func (c *Client) Swipe(ctx context.Context, request *SwipeRequest) (*SwipeResponse, error) {
	if request == nil {
		return nil, ErrSwipeMissingRequest
	}

	user := store.User{}
	searchConditions := map[string]interface{}{"ID": request.UserID}
	if err := c.store.FindUser(ctx, &user, searchConditions); err != nil {
		return nil, err
	}

	if err := c.store.UpsertSwipe(ctx, &store.Swipe{
		MatchedID: request.MatchedID,
		IsDesired: request.IsDesired,
		UserID:    user.ID,
	}); err != nil {
		return nil, err
	}

	// If the user has swipped no on this person, we can short circuit and
	// not be concerned if they'd matched or not
	if request.IsDesired != "YES" {
		return &SwipeResponse{HasMatched: false}, nil
	}

	hasMatched, err := c.store.HasMatched(ctx, user.ID, request.MatchedID)
	if err != nil {
		return nil, err
	}

	response := &SwipeResponse{HasMatched: hasMatched}
	if response.HasMatched {
		response.MatchedID = request.MatchedID
	}

	return response, nil
}
