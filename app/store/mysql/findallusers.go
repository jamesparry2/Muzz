package mysql

import (
	"context"
	"fmt"

	"github.com/jamesparry2/Muzz/app/store"
)

/*
SELECT *
FROM muzz.users as users
INNER JOIN muzz.locations as location
ON location.user_id = users.id
WHERE (SELECT ST_Distance_Sphere(point ('51.601151', '-3.3499'), point(location.lat, `long`)) * .000621371192) <= location.distance_from_me
AND users.id != 1 AND users.id NOT IN (2) AND users.gender = 'male' AND (users.age BETWEEN 25 AND 30);
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
		// Gorm doesn't seem to play well ST_Distance_Sphere, so formating using DB values with floats, if this was strings this should not be done
		query.InnerJoins("Location").Where("@query <= @distance_from_me", map[string]interface{}{
			"query":            fmt.Sprintf("ST_Distance_Sphere(point ('%f', '%f'), point(lat, `long`)) * .000621371192", user.Location.Lat, user.Location.Long),
			"distance_from_me": user.Location.DistanceFromMe,
		})
	}

	result := query.Find(&users)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}
