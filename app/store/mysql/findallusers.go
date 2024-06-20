package mysql

import (
	"context"

	"github.com/jamesparry2/Muzz/app/store"
)

func (c *Client) FindAllUsers(ctx context.Context, user *store.User) ([]store.User, error) {
	users := []store.User{}

	query := c.db.Where("id != ?", user.ID)
	if len(user.Swipes) != 0 {
		swippedIds := []uint{}
		for _, swipped := range user.Swipes {
			swippedIds = append(swippedIds, swipped.ID)
		}

		query.Where("id NOT IN (?)", swippedIds)
	}

	if user.Preferences != nil {
		if user.Preferences.Gender != "" {
			query.Where("gender = ?", user.Preferences.Gender)
		}

		if user.Preferences.MinimumAge >= 18 && user.Preferences.MaximumAge >= 19 {
			query.Where("age BETWEEN ? AND ?", user.Preferences.MinimumAge, user.Preferences.MaximumAge)
		}
	}

	result := query.Find(&users)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}
