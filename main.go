package main

import (
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mifaabiyyu/backend-go/api"
	"github.com/mifaabiyyu/backend-go/internal/auth"
	"github.com/mifaabiyyu/backend-go/internal/db"
	"github.com/mifaabiyyu/backend-go/internal/env"
	"github.com/mifaabiyyu/backend-go/internal/mailer"
	"github.com/mifaabiyyu/backend-go/internal/ratelimiter"
	"github.com/mifaabiyyu/backend-go/internal/store"
	"github.com/mifaabiyyu/backend-go/internal/store/cache"
	"go.uber.org/zap"
)

func main() {
	cfg := api.Config{
		Addr: env.GetString("PORT", ":3000"),
		Db: api.DbConfig{
			Addr:         env.GetString("DB", "postgresql://root:password@localhost:5432/backend?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		RedisCfg: api.RedisConfig{
			Addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			Pw:      env.GetString("REDIS_PW", ""),
			Db:      env.GetInt("REDIS_DB", 0),
			Enabled: env.GetBool("REDIS_ENABLED", false),
		},
		Env: env.GetString("ENV", "development"),
		Mail: api.MailConfig{
			Exp:       time.Hour * 24 * 3, // 3 days
			FromEmail: env.GetString("FROM_EMAIL", ""),
			SendGrid: api.SendGridConfig{
				ApiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			MailTrap: api.MailTrapConfig{
				ApiKey: env.GetString("MAILTRAP_API_KEY", "3b94f4137017efc8e12be18da20096c3"),
			},
		},
		Auth: api.AuthConfig{
			Basic: api.BasicConfig{
				User: env.GetString("AUTH_BASIC_USER", "admin"),
				Pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			Token: api.TokenConfig{
				Secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				Exp:    time.Hour * 24 * 3, // 3 days
				Iss:    "gophersocial",
			},
		},
		RateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	dbCon, err := db.New(
		cfg.Db.Addr,
		cfg.Db.MaxOpenConns,
		cfg.Db.MaxIdleConns,
		cfg.Db.MaxIdleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	// Cache
	var rdb *redis.Client
	if cfg.RedisCfg.Enabled {
		rdb = cache.NewRedisClient(cfg.RedisCfg.Addr, cfg.RedisCfg.Pw, cfg.RedisCfg.Db)
		logger.Info("redis cache connection established")

		defer rdb.Close()
	}

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.RateLimiter.RequestsPerTimeFrame,
		cfg.RateLimiter.TimeFrame,
	)

	// Mailer
	// mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	mailtrap, err := mailer.NewMailTrapClient(cfg.Mail.MailTrap.ApiKey, cfg.Mail.FromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	// Authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.Auth.Token.Secret,
		cfg.Auth.Token.Iss,
		cfg.Auth.Token.Iss,
	)

	defer dbCon.Close()
	logger.Info("database connection pool established")

	store := store.NewStore(dbCon)
	cacheStorage := cache.NewRedisStorage(rdb)

	app := api.Application{
		Config:        cfg,
		Store:         store,
		CacheStorage:  cacheStorage,
		Logger:        logger,
		Mailer:        mailtrap,
		Authenticator: jwtAuthenticator,
		RateLimiter:   rateLimiter,
	}
	app.InitMiddleware()
	mux := app.Mount()

	log.Fatal(app.Run(mux))
}
