package mysql

import (
	"context"

	"github.com/jamesparry2/Muzz/app/store"
	"gorm.io/gorm"
)

func (c *Client) HasMatched(ctx context.Context, userID, matchedID uint) (bool, error) {
	if userID == 0 || matchedID == 0 {
		return false, store.ErrHasMatchedMissingIDs
	}

	conditions := map[string]interface{}{
		"user_id":    matchedID,
		"matched_id": userID,
		"is_desired": "YES",
	}

	swipe := store.Swipe{}
	if resp := c.db.Find(&swipe, conditions); resp.Error != nil {
		if resp.Error != nil {
			switch err := resp.Error; err {
			// If we didn't find a record, we can assume there isn't a match, however
			// any other error returned is a valid error
			case gorm.ErrRecordNotFound:
				return false, nil
			default:
				return false, resp.Error
			}
		}

		return false, nil
	}

	// Since we have had a row affected or no error returned, we can extropolate that there was a match
	return true, nil
}
