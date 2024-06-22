package app

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/jamesparry2/Muzz/app/auth"
	"github.com/jamesparry2/Muzz/app/core"
	"github.com/jamesparry2/Muzz/app/handler"
	"github.com/jamesparry2/Muzz/app/store/mysql"
	"github.com/jamesparry2/Muzz/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag"
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

	setupSwaggerDetails(docs.SwaggerInfo)

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
	server.GET("swagger/*", echoSwagger.WrapHandler)
}

func setupAuthRoutes(server *echo.Group, handlers *handler.Handler) {
	server.POST("user/:id/swipe", handlers.Swipe)
	server.GET("user/:id/discovery", handlers.Discovery)
	server.POST("user/:id/preference", handlers.Preference)
	server.POST("user/:id/location", handlers.Location)
}

func setupSwaggerDetails(spec *swag.Spec) {
	spec.Title = "Muzz Test API"
	spec.Version = "1.0"
	spec.Description = "This is a mini api that allows basic user creation, discovery and matching"
	spec.Host = "localhost:5001"
	spec.BasePath = "/"
	spec.Schemes = []string{"http"}
}

func handleAuthMiddleware(auth *auth.Auth, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, handler.NewAPIError(http.StatusUnauthorized, "auth", "no token provided"))
		}

		modifiedString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := auth.VerifyToken(modifiedString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, handler.NewAPIError(http.StatusUnauthorized, "auth", err.Error()))
		}

		userId, err := handler.GetUserIDPathParam(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, handler.NewAPIError(http.StatusUnauthorized, "auth", err.Error()))
		}

		sub, err := token.Claims.GetSubject()
		if err != nil {
			return c.JSON(http.StatusUnauthorized, handler.NewAPIError(http.StatusUnauthorized, "auth", err.Error()))
		}

		if sub != fmt.Sprint(userId) {
			return c.JSON(http.StatusUnauthorized, handler.NewAPIError(http.StatusUnauthorized, "auth", "not allowed to access the requested resource"))
		}

		return next(c)
	}
}
