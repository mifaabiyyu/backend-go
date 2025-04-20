package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	authentication "github.com/mifaabiyyu/backend-go/cmd/api/auth"
	"github.com/mifaabiyyu/backend-go/cmd/api/user"
	"github.com/mifaabiyyu/backend-go/internal/auth"
	"github.com/mifaabiyyu/backend-go/internal/env"
	"github.com/mifaabiyyu/backend-go/internal/mailer"
	"github.com/mifaabiyyu/backend-go/internal/ratelimiter"
	"github.com/mifaabiyyu/backend-go/internal/store"
	"github.com/mifaabiyyu/backend-go/internal/store/cache"
	"github.com/mifaabiyyu/backend-go/utils"

	"github.com/swaggo/swag/example/basic/docs"
	"go.uber.org/zap"
)

const version = "1.1.0"

type Application struct {
	ApiURL        string
	Config        Config
	CacheStorage  cache.Storage
	Logger        *zap.SugaredLogger
	Store         *store.Store
	RateLimiter   ratelimiter.Limiter
	Authenticator auth.Authenticator
	Mailer        mailer.Client
	middleware    AppAll
}

type Config struct {
	Addr        string
	Db          DbConfig
	Env         string
	ApiURL      string
	Mail        MailConfig
	FrontendURL string
	Auth        AuthConfig
	RedisCfg    RedisConfig
	RateLimiter ratelimiter.Config
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type RedisConfig struct {
	Addr    string
	Pw      string
	Db      int
	Enabled bool
}

type AuthConfig struct {
	Basic BasicConfig
	Token TokenConfig
}

type TokenConfig struct {
	Secret string
	Exp    time.Duration
	Iss    string
}

type BasicConfig struct {
	User string
	Pass string
}

type MailConfig struct {
	SendGrid  SendGridConfig
	MailTrap  MailTrapConfig
	FromEmail string
	Exp       time.Duration
}

type MailTrapConfig struct {
	ApiKey string
}

type SendGridConfig struct {
	ApiKey string
}

func (app *Application) InitMiddleware() {
	app.middleware = AppAll{
		AppWrapper: &utils.AppWrapper{
			Logger: app.Logger,
		},
		Application: app, // Pastikan `app` implementasikan method `GetUserFromCacheOrDb`
	}
}

func (app *Application) Mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5174")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(v1 chi.Router) {
		// Grouped routes for users
		app.mountUserRoutes(v1)

		// In the future:
		// app.mountProductRoutes(v1)
		// app.mountAuthRoutes(v1)
	})

	return r
}

func (app *Application) mountUserRoutes(r chi.Router) {
	userHandler := user.InitUserModule(app.Store.Queries)
	authHandler := authentication.InitAuthModule(app.Store, app.middleware.AppWrapper, app.Authenticator)

	r.Route("/users", func(r chi.Router) {
		r.With(app.middleware.AuthTokenMiddleware, app.middleware.RequirePermission("user:read")).Get("/", userHandler.ListUsers)
		r.With(app.middleware.AuthTokenMiddleware, app.middleware.RequirePermission("user:read")).Get("/{id}", userHandler.GetUser)

	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})
}

func (app *Application) Run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.Config.ApiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.Logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.Logger.Infow("server has started", "addr", app.Config.Addr, "env", app.Config.Env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.Logger.Infow("server has stopped", "addr", app.Config.Addr, "env", app.Config.Env)

	return nil
}
