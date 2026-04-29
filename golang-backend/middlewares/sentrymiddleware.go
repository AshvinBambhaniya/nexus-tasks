package middlewares

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
)

func SentryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Capture panic in Sentry
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				sentry.CaptureException(err)
				sentry.Flush(2 * time.Second)

				c.Status(500).JSON(fiber.Map{"error": "internal server error"})
			}
		}()

		return c.Next()
	}
}
