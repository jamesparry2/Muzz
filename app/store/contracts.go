package store

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrFindUserMissingUserDetails = errors.New("missing user details")
	ErrFindUserNotFound           = errors.New("unable to find record")

	ErrNewClientMissingRequiredOptions = errors.New("missing critical client optiosn: host, port, username, password or database")

	ErrUpsertUserMissingUser = errors.New("missing user details to upsert")
	ErrUpsertUserDBError     = errors.New("failed to save the new user to the DB")

	ErrUpsertSwipeMissingSwipe = errors.New("missing swipe details to upsert")
	ErrUpsertSwipeDBError      = errors.New("failed to save swipe to the DB")

	ErrHasMatchedMissingIDs = errors.New("missing ids to verify if matched")

	ErrUpsertPreferencesMissing = errors.New("missing preference details to upsert")
	ErrUpsertPreferenceDBError  = errors.New("failed to save preference to the DB")
)

type User struct {
	gorm.Model
	Email          string
	Password       string
	Name           string
	Gender         string
	Age            int
	DistanceFromMe float64 `gorm:"distance_from_me"`

	Location    *Location
	Preferences *Preferences

	Swipes []*Swipe
}

type Location struct {
	gorm.Model
	Lat  float64
	Long float64

	UserID uint
}

type Preferences struct {
	gorm.Model
	MinimumAge int
	MaximumAge int
	Gender     string

	UserID uint
}

type Swipe struct {
	gorm.Model
	MatchedID uint
	IsDesired string

	UserID uint
}

type StoreIface interface {
	UpsertUser(ctx context.Context, request *User) error
	FindUser(ctx context.Context, user *User, conditions map[string]interface{}) error
	UpsertSwipe(ctx context.Context, swipe *Swipe) error
	HasMatched(ctx context.Context, userID, matchedID uint) (bool, error)
	FindAllUsers(ctx context.Context, user *User) ([]User, error)
}
