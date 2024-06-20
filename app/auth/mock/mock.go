package mock

import "github.com/golang-jwt/jwt/v5"

type AuthMock struct {
	MockCreateToken     func(subject, aud string) (string, error)
	MockVerifyToken     func(tokenStr string) (*jwt.Token, error)
	MockEncryptPassword func(pwd string) (string, error)
}

func (a *AuthMock) CreateToken(subject, aud string) (string, error) {
	return a.MockCreateToken(subject, aud)
}

func (a *AuthMock) VerifyToken(tokenStr string) (*jwt.Token, error) {
	return a.MockVerifyToken(tokenStr)
}

func (a *AuthMock) EncryptPassword(pwd string) (string, error) {
	return a.MockEncryptPassword(pwd)
}
