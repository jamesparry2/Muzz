package mysql_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jamesparry2/Muzz/app/store"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFindAllUsers(t *testing.T) {
	t.Run("should run a base select when no swipes, location or preferences are provided", func(t *testing.T) {
		client, mock, err := setupMockDB()

		assert.NoError(t, err, "setup should not error")

		rows := sqlmock.NewRows([]string{"email", "password"}).
			AddRow("email1", "password2").
			AddRow("email2", "password2")

		mock.ExpectQuery("SELECT \\* FROM `users` WHERE users\\.id != \\? AND `users`\\.`deleted_at` IS NULL").WillReturnRows(rows)

		users, err := client.FindAllUsers(context.Background(), &store.User{
			Model: gorm.Model{ID: 2},
		})

		assert.NoError(t, err, "no error should be returned")
		assert.Len(t, users, 2, "unexpected amount of users returned")
	})

	t.Run("should include filtered out IDs when swipes are provided", func(t *testing.T) {
		client, mock, err := setupMockDB()

		assert.NoError(t, err, "setup should not error")

		rows := sqlmock.NewRows([]string{"email", "password"}).
			AddRow("email1", "password2").
			AddRow("email2", "password2")

		mock.ExpectQuery("SELECT \\* FROM `users` WHERE users\\.id != \\? AND users\\.id NOT IN \\(\\?\\) AND `users`\\.`deleted_at` IS NULL").WillReturnRows(rows)

		users, err := client.FindAllUsers(context.Background(), &store.User{
			Model: gorm.Model{ID: 2},
			Swipes: []*store.Swipe{{
				MatchedID: 3,
			}},
		})

		assert.NoError(t, err, "no error should be returned")
		assert.Len(t, users, 2, "unexpected amount of users returned")
	})

	t.Run("should include preferences when provided on the user request", func(t *testing.T) {
		client, mock, err := setupMockDB()

		assert.NoError(t, err, "setup should not error")

		rows := sqlmock.NewRows([]string{"email", "password"}).
			AddRow("email1", "password2").
			AddRow("email2", "password2")

		mock.ExpectQuery("SELECT \\* FROM `users` WHERE users\\.id != \\? AND users\\.gender = \\? AND \\(users\\.age BETWEEN \\? AND \\?\\) AND `users`\\.`deleted_at` IS NULL").WillReturnRows(rows)

		users, err := client.FindAllUsers(context.Background(), &store.User{
			Model: gorm.Model{ID: 2},
			Preferences: &store.Preferences{
				Gender:     "female",
				MinimumAge: 24,
				MaximumAge: 28,
			},
		})

		assert.NoError(t, err, "no error should be returned")
		assert.Len(t, users, 2, "unexpected amount of users returned")
	})

	t.Run("should include locations when provided on the user reqeust", func(t *testing.T) {
		client, mock, err := setupMockDB()

		assert.NoError(t, err, "setup should not error")

		rows := sqlmock.NewRows([]string{"email", "password"}).
			AddRow("email1", "password2").
			AddRow("email2", "password2")

		mock.ExpectQuery("SELECT users\\.id, users\\.name, users\\.password, users\\.gender, users\\.age, ST_Distance_Sphere\\(point \\(\\?, \\?\\), point\\(lat, `long`\\)\\) \\* \\.000621371192 as distance_from_me FROM `users` INNER JOIN `locations` `Location` ON `users`\\.`id` = `Location`\\.`user_id` AND `Location`\\.`deleted_at` IS NULL WHERE users\\.id != \\? AND `users`\\.`deleted_at` IS NULL ORDER BY distance_from_me DESC").WillReturnRows(rows)

		users, err := client.FindAllUsers(context.Background(), &store.User{
			Model: gorm.Model{ID: 2},
			Location: &store.Location{
				Lat:  0,
				Long: 1,
			},
		})

		assert.NoError(t, err, "no error should be returned")
		assert.Len(t, users, 2, "unexpected amount of users returned")
	})
}
