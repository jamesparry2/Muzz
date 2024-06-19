package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrSecretKeyNotProvided = errors.New("missing secret key")
	ErrInvalidToken         = errors.New("invalid token provided")
)

type AuthIface interface {
	CreateToken(subject, aud string) (string, error)
	VerifyToken(tokenStr string) (*jwt.Token, error)
	EncryptPassword(pwd string) (string, error)
}

type AuthOptions struct {
	SecretKey string
}

type Auth struct {
	secretKey string
}

func NewAuth(opts *AuthOptions) *Auth {
	return &Auth{
		secretKey: opts.SecretKey,
	}
}

func (a *Auth) CreateToken(subject, aud string) (string, error) {
	if a.secretKey == "" {
		return "", ErrSecretKeyNotProvided
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": subject,
		"aud": aud,
		"iss": "muzz",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	token, err := claims.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *Auth) VerifyToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token, nil
}

func (a Auth) EncryptPassword(pwd string) (string, error) {
	return pwd, nil
}
