package mysql

import (
	"context"

	"github.com/jamesparry2/Muzz/app/store"
)

/*
SELECT users.id, users.name, users.password, users.gender, users.age, ST_Distance_Sphere(point (51.601151, -3.3499), point(lat, `long`)) * .000621371192 as distance_from_me
FROM `users`
INNER JOIN `locations` `Location`
ON `users`.`id` = `Location`.`user_id` AND `Location`.`deleted_at` IS NULL
WHERE users.id != 1
AND users.id NOT IN (2)
AND users.gender = 'male'
AND (users.age BETWEEN 25 AND 30)
AND `users`.`deleted_at` IS NULL
ORDER BY distance_from_me DESC
*/
func (c *Client) FindAllUsers(ctx context.Context, user *store.User) ([]store.User, error) {
	users := []store.User{}

	query := c.db.Where("users.id != ?", user.ID)
	if len(user.Swipes) != 0 {
		swippedIds := []uint{}
		for _, swipped := range user.Swipes {
			swippedIds = append(swippedIds, swipped.ID)
		}

		query.Where("users.id NOT IN (?)", swippedIds)
	}

	if user.Preferences != nil {
		if user.Preferences.Gender != "" {
			query.Where("users.gender = ?", user.Preferences.Gender)
		}

		if user.Preferences.MinimumAge >= 18 && user.Preferences.MaximumAge >= 19 {
			query.Where("users.age BETWEEN ? AND ?", user.Preferences.MinimumAge, user.Preferences.MaximumAge)
		}
	}

	if user.Location != nil {
		query.InnerJoins("Location").Select("users.id, users.name, users.password, users.gender, users.age, ST_Distance_Sphere(point (?, ?), point(lat, `long`)) * .000621371192 as distance_from_me", user.Location.Lat, user.Location.Long).Order("distance_from_me DESC")
	}

	// Just need to map distance in
	result := query.Find(&users)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}
