package mock

import "github.com/jamesparry2/Muzz/app/core"

type MockCore struct {
	MockCreateUser func(request *core.CreateUserRequest) error
}

func (mc *MockCore) CreateUser(r *core.CreateUserRequest) error {
	return mc.MockCreateUser(r)
}
