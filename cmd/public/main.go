package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"komodo-address-api/internal/config"

	"github.com/rdevitto86/komodo-forge-sdk-go/api/handlers/health"

	srv "github.com/rdevitto86/komodo-forge-sdk-go/api/server"
	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secretsmanager"
	"github.com/rdevitto86/komodo-forge-sdk-go/crypto/jwt"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// init runs once per execution environment (cold start on Lambda, once on Fargate/local).
// Order matters: logger first, then SM (loads JWT_* keys), then JWT init.
func init() {
	logger.Init(
		os.Getenv(config.APP_NAME),
		os.Getenv(config.LOG_LEVEL),
		os.Getenv(config.ENV),
	)

	smCfg := awsSM.Config{
		Region:   os.Getenv(config.AWS_REGION),
		Endpoint: os.Getenv(config.AWS_ENDPOINT),
		Prefix:   os.Getenv(config.AWS_SECRET_PREFIX),
		Batch:    os.Getenv(config.AWS_SECRET_BATCH),
		Keys: []string{
			config.JWT_PUBLIC_KEY,
			config.JWT_PRIVATE_KEY,
			config.JWT_AUDIENCE,
			config.JWT_ISSUER,
			config.JWT_KID,
			config.ADDRESS_PROVIDER_API_KEY,
			config.MAX_CONTENT_LENGTH,
			config.RATE_LIMIT_RPS,
			config.RATE_LIMIT_BURST,
		},
	}
	sm, err := awsSM.New(context.Background(), smCfg)
	if err != nil {
		logger.Fatal("failed to initialize secrets manager", err)
		os.Exit(1)
	}

	secrets, err := sm.GetSecrets(smCfg.Keys, smCfg.Prefix, smCfg.Batch)
	if err != nil {
		logger.Fatal("failed to fetch secrets", err)
		os.Exit(1)
	}
	for k, v := range secrets {
		os.Setenv(k, v)
	}

	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	logger.Info("address-api: bootstrap complete")
}

func main() {
	// stack := []func(http.Handler) http.Handler{
	// 	mw.RequestIDMiddleware,
	// 	mw.TelemetryMiddleware,
	// 	mw.RateLimiterMiddleware,
	// 	mw.CORSMiddleware,
	// 	mw.SecurityHeadersMiddleware,
	// 	mw.AuthMiddleware,
	// 	mw.NormalizationMiddleware,
	// 	mw.RuleValidationMiddleware,
	// 	mw.SanitizationMiddleware,
	// }

	// providerClient := provider.NewClient(os.Getenv("ADDRESS_PROVIDER_API_KEY"))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health.HealthHandler)
	// mux.Handle("POST /addresses/validate", mw.Chain(handlers.Validate(providerClient), stack...))
	// mux.Handle("POST /addresses/normalize", mw.Chain(handlers.Normalize(providerClient), stack...))
	// mux.Handle("POST /addresses/geocode", mw.Chain(handlers.Geocode(providerClient), stack...))

	server := &http.Server{
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	srv.Run(server, os.Getenv(config.PORT), 30*time.Second)
}
