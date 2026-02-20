package middleware

import (
	"basekarya-backend/pkg/response"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type RateLimiterMiddleware struct {
}

func NewRateLimiterMiddleware() *RateLimiterMiddleware {
	return &RateLimiterMiddleware{}
}

func (m *RateLimiterMiddleware) Init() echo.MiddlewareFunc {
	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(5.0 / 60.0),
				Burst:     5,
				ExpiresIn: 3 * time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return response.NewResponses[any](context,
				http.StatusTooManyRequests, "Too many login attempts. Please try again in 1 minute.", nil, nil, nil)
		},
	})
}
