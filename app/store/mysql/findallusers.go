package mysql

import (
	"context"

	"github.com/jamesparry2/Muzz/app/store"
)

func (c *Client) FindAllUsers(ctx context.Context, user *store.User) ([]store.User, error) {
	users := []store.User{}

	// Filter out the following:
	// - own user
	// - any user that this user has already swipped on (no need to consider responded swipes)

	result := c.db.Find(&users)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}
