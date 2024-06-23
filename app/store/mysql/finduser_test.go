package mysql_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jamesparry2/Muzz/app/store"
	"github.com/stretchr/testify/assert"
)

func TestFindUser(t *testing.T) {
	t.Run("should return an error when no user to map to is passed in", func(t *testing.T) {
		client, _, err := setupMockDB()

		assert.NoError(t, err, "setup should not error")

		dbErr := client.FindUser(context.Background(), nil, map[string]interface{}{})

		assert.ErrorIs(t, dbErr, store.ErrFindUserMissingUserDetails, "unexpected error was returned")
	})

	t.Run("should return a standard when user not found error", func(t *testing.T) {
		client, mock, err := setupMockDB()

		assert.NoError(t, err, "setup should not error")

		mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`\\.`deleted_at` IS NULL").WillReturnError(store.ErrFindUserNotFound)

		dbErr := client.FindUser(context.Background(), &store.User{}, map[string]interface{}{})

		assert.ErrorIs(t, dbErr, store.ErrFindUserNotFound, "unexpected error was returned")
	})

	t.Run("should return a standard when user not found error", func(t *testing.T) {
		client, mock, err := setupMockDB()

		assert.NoError(t, err, "setup should not error")

		unExpectedErr := errors.New("boo, didn't expect this error did you")
		mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`\\.`deleted_at` IS NULL").WillReturnError(unExpectedErr)

		dbErr := client.FindUser(context.Background(), &store.User{}, map[string]interface{}{})

		assert.ErrorIs(t, dbErr, unExpectedErr, "unexpected error was returned")
	})

	t.Run("should return a mapped user from a success query", func(t *testing.T) {
		client, mock, err := setupMockDB()

		assert.NoError(t, err, "setup should not error")

		userRow := mock.NewRows([]string{"email", "password"}).
			AddRow("email1", "password2")

		mock.ExpectQuery("SELECT \\* FROM `users` WHERE `email` = \\? AND `users`\\.`deleted_at` IS NULL").WillReturnRows(userRow)

		dbErr := client.FindUser(context.Background(), &store.User{}, map[string]interface{}{"email": "email1"})

		assert.NoError(t, dbErr, "unexpected error was returned")
	})
}
