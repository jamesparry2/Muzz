package mysql_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jamesparry2/Muzz/app/store"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestHasMatched(t *testing.T) {
	t.Run("should return an error when ids are provided to check if matched", func(t *testing.T) {
		client, _, _ := setupMockDB()

		isMatched, err := client.HasMatched(context.Background(), 0, 0)

		assert.ErrorIs(t, err, store.ErrHasMatchedMissingIDs, "unexpected error was returned")
		assert.False(t, isMatched, "should be false")
	})

	t.Run("should return false but no error when no DB record is found, as it assume the person hasn't swipped on them", func(t *testing.T) {
		client, mock, _ := setupMockDB()

		mock.ExpectQuery("SELECT \\* FROM `swipes` WHERE \\(`is_desired` = \\? AND `matched_id` = \\? AND `user_id` = \\?\\) AND `swipes`\\.`deleted_at` IS NULL").WillReturnError(gorm.ErrRecordNotFound)

		isMatched, err := client.HasMatched(context.Background(), 1, 2)

		assert.NoError(t, err, "unexpected error was returned")
		assert.False(t, isMatched, "should be false")
	})

	t.Run("should return false with an error when a terminal DB error is returned", func(t *testing.T) {
		client, mock, _ := setupMockDB()

		terminalError := errors.New("unexpected error")
		mock.ExpectQuery("SELECT \\* FROM `swipes` WHERE \\(`is_desired` = \\? AND `matched_id` = \\? AND `user_id` = \\?\\) AND `swipes`\\.`deleted_at` IS NULL").WillReturnError(terminalError)

		isMatched, err := client.HasMatched(context.Background(), 1, 2)

		assert.ErrorIs(t, err, terminalError, "unexpected error was returned")
		assert.False(t, isMatched, "should be false")
	})

	t.Run("should return true when a match is found", func(t *testing.T) {
		client, mock, _ := setupMockDB()

		rows := mock.NewRows([]string{"user_id", "matched_id"}).AddRow(1, 1)

		terminalError := errors.New("unexpected error")
		mock.ExpectQuery("SELECT \\* FROM `swipes` WHERE \\(`is_desired` = \\? AND `matched_id` = \\? AND `user_id` = \\?\\) AND `swipes`\\.`deleted_at` IS NULL").WillReturnRows(rows)

		isMatched, err := client.HasMatched(context.Background(), 1, 2)

		assert.NoError(t, err, terminalError, "unexpected error was returned")
		assert.True(t, isMatched, "should be false")
	})
}
