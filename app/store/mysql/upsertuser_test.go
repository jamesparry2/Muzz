package mysql_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jamesparry2/Muzz/app/store"
	"github.com/stretchr/testify/assert"
)

func TestUpsertUser(t *testing.T) {
	t.Run("should return an error when no user is passed in to be saved", func(t *testing.T) {
		db, _, _ := setupMockDB()

		assert.ErrorIs(t, db.UpsertUser(context.Background(), nil), store.ErrUpsertUserMissingUser)
	})

	t.Run("should return an error when the user failed to save", func(t *testing.T) {
		db, mock, _ := setupMockDB()

		mock.ExpectExec("").WillReturnError(errors.New("termianl DB error"))

		assert.ErrorIs(t, db.UpsertUser(context.Background(), &store.User{}), store.ErrUpsertUserDBError)
	})
}
