package mysql

import (
	"context"

	"github.com/jamesparry2/Muzz/app/store"
)

func (c *Client) UpsertSwipe(ctx context.Context, swipe *store.Swipe) error {
	if swipe == nil {
		return store.ErrUpsertSwipeMissingSwipe
	}

	if resp := c.db.Save(swipe); resp.Error != nil || resp.RowsAffected == 0 {
		return store.ErrUpsertSwipeDBError
	}

	return nil
}
