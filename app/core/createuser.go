package core

import (
	"context"
	"errors"
	"math/rand"

	"github.com/jamesparry2/Muzz/app/store"
)

var (
	ErrCreateUserMissingRequest = errors.New("missing inbound request unable to create user")
	ErrCreateUserUsernameInUser = errors.New("email provided is already in use")
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

type CreateUserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

func (c *Client) CreateUser(ctx context.Context, request *CreateUserRequest) (*CreateUserResponse, error) {
	if request == nil {
		return nil, ErrCreateUserMissingRequest
	}

	user := store.User{}
	searchConditions := map[string]interface{}{"email": request.Email}
	// We know if the record isn't found, thats not a terminal error case as in this flow we may not be registed
	if err := c.store.FindUser(ctx, &user, searchConditions); (err != nil && !errors.Is(err, store.ErrFindUserNotFound)) || user.ID > 0 {
		if user.ID > 0 {
			return nil, ErrCreateUserUsernameInUser
		}
		return nil, err
	}

	encryptedPassword, err := c.auth.EncryptPassword(request.Password)
	if err != nil {
		return nil, err
	}

	gender := map[int]string{
		0: "male",
		1: "female",
	}

	userToSave := &store.User{
		Email:    request.Email,
		Password: encryptedPassword,
		Name:     request.Name,
		Gender:   gender[rand.Intn(2)],
		Age:      rand.Intn(100-18) + 18,
		Location: &store.Location{
			Lat:  float64(rand.Intn(91)),
			Long: float64(rand.Intn(91)),
		},
	}
	if err := c.store.UpsertUser(ctx, userToSave); err != nil {
		return nil, err
	}

	return mapStoreUserToResponse(userToSave), nil
}

func mapStoreUserToResponse(user *store.User) *CreateUserResponse {
	return &CreateUserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		Gender:   user.Gender,
		Age:      user.Age,
		Name:     user.Name,
	}
}
