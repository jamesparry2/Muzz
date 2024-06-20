package mysql

import (
	"context"

	"github.com/jamesparry2/Muzz/app/store"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *Client) FindUser(ctx context.Context, user *store.User, conditions map[string]interface{}) error {
	if user == nil {
		return store.ErrFindUserMissingUserDetails
	}

	if resp := c.db.Preload(clause.Associations).Find(&user, conditions); resp.Error != nil {
		switch err := resp.Error; err {
		case gorm.ErrRecordNotFound:
			return store.ErrFindUserNotFound
		default:
			return resp.Error
		}
	}

	return nil
}
