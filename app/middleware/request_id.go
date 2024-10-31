package middleware

import (
	"context"
	"user_api/utils/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

func RequestID(c *fiber.Ctx) error {
	rid := c.Get("X-Request-Id")
	if rid == "" {
		rid = uuid.New().String()
	}

	ctxValue := context.WithValue(c.UserContext(), RequestIDKey, rid)
	ctx := logger.WithLogger(ctxValue, logger.GetLogger(ctxValue).WithFields(logger.Fields{
		"request_id": rid,
	}))

	c.SetUserContext(ctx)

	return c.Next()
}
