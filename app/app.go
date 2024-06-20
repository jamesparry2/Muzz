package app

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/jamesparry2/Muzz/app/auth"
	"github.com/jamesparry2/Muzz/app/core"
	"github.com/jamesparry2/Muzz/app/handler"
	"github.com/jamesparry2/Muzz/app/store/mysql"
	"github.com/labstack/echo/v4"
)

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUsername string `env:"DB_USERNAME"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`

	AuthSecretKey string `env:"AUTH_SECRET_KEY"`

	ApiPort string `env:"API_PORT"`
}

func Run() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		os.Exit(1)
	}

	// Add custom env injection for 12 factor standard
	store, err := mysql.NewClientWithConection(&mysql.ClientOptions{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Username: cfg.DBUsername,
		Password: cfg.DBPassword,
		Database: cfg.DBName,
	})
	if err != nil {
		os.Exit(1)
	}

	auth := auth.NewAuth(&auth.AuthOptions{
		SecretKey: cfg.AuthSecretKey,
	})

	core := core.NewClient(&core.ClientOptions{
		Store: store,
		Auth:  auth,
	})

	handlers := handler.NewHandler(&handler.HandlerOption{
		Core: core,
	})

	server := echo.New()
	authGroup := server.Group("")
	authGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return handleAuthMiddleware(auth, next)
	})

	setupNonAuthRoutes(server, handlers)
	setupAuthRoutes(authGroup, handlers)

	server.Logger.Fatal(server.Start(cfg.ApiPort))
}

func setupNonAuthRoutes(server *echo.Echo, handlers *handler.Handler) {
	// Testing endpoint for creating a random user
	server.POST("user/create", handlers.CreateUser)
	server.POST("login", handlers.Login)
}

func setupAuthRoutes(server *echo.Group, handlers *handler.Handler) {
	// To add
	server.POST("user/:id/swipe", handlers.Swipe)
	server.GET("user/:id/discovery", handlers.Discovery)
	server.POST("user/:id/preference", handlers.Preference)

	server.POST("user/:id/location", func(c echo.Context) error { return nil })
}

// Tidy this up tomorrow
func handleAuthMiddleware(auth *auth.Auth, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return errors.New("")
		}

		modifiedString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := auth.VerifyToken(modifiedString)
		if err != nil {
			return err
		}

		userId, err := handler.GetUserIDPathParam(c)
		if err != nil {
			return err
		}

		sub, err := token.Claims.GetSubject()
		if err != nil {
			return err
		}

		if sub != fmt.Sprint(userId) {
			return errors.New("token not allowed to access this resource")
		}

		return next(c)
	}
}
