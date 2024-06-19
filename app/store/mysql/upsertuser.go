package mysql

import (
	"context"

	"github.com/jamesparry2/Muzz/app/store"
)

func (c *Client) UpsertUser(ctx context.Context, request *store.User) error {
	if request == nil {
		return store.ErrUpsertUserMissingUser
	}

	if resp := c.db.Save(request); resp.Error != nil || resp.RowsAffected == 0 {
		return store.ErrUpsertUserDBError
	}

	return nil
}
