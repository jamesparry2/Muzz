package mysql_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jamesparry2/Muzz/app/store"
	"github.com/stretchr/testify/assert"
)

func TestUpsertSwipe(t *testing.T) {
	t.Run("should return an error when no swipe is passed in to be saved", func(t *testing.T) {
		db, _, _ := setupMockDB()

		assert.ErrorIs(t, db.UpsertSwipe(context.Background(), nil), store.ErrUpsertSwipeMissingSwipe)
	})

	t.Run("should return an error when the swipe failed to save", func(t *testing.T) {
		db, mock, _ := setupMockDB()

		mock.ExpectExec("").WillReturnError(errors.New("termianl DB error"))

		assert.ErrorIs(t, db.UpsertSwipe(context.Background(), &store.Swipe{}), store.ErrUpsertSwipeDBError)
	})
}
