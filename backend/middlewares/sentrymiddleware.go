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

				if err := c.Status(500).JSON(fiber.Map{"error": "internal server error"}); err != nil {
					// We use a blank assignment here to satisfy the linter's requirement for checking the error,
					// as there is no meaningful recovery path if sending the error response itself fails.
					_ = err
				}
			}
		}()

		return c.Next()
	}
}
