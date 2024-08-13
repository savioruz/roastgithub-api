package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"regexp"
	"strings"
	"time"
)

var (
	uri = []*regexp.Regexp{
		regexp.MustCompile(`^/monitor`),
		regexp.MustCompile(`^/swagger`),
	}
	ua = []string{
		"Mozilla",
		"Chrome",
		"Safari",
		"Swagger",
	}
)

func LimiterMiddleware(a *fiber.App) {
	a.Use(limiter.New(
		limiter.Config{
			Next: func(c *fiber.Ctx) bool {
				return whitelistPath(c) || whitelistRequest(c)
			},
			Max:               30,
			Expiration:        60 * time.Second,
			LimiterMiddleware: limiter.SlidingWindow{},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Too many requests",
				})
			},
		}))
}

// whitelistPath func for checking the next middleware uri
func whitelistPath(c *fiber.Ctx) bool {
	originalURL := strings.ToLower(c.OriginalURL())

	for _, pattern := range uri {
		if pattern.MatchString(originalURL) {
			return true
		}
	}
	return false
}

// whitelistRequest func for checking the next middleware request
func whitelistRequest(c *fiber.Ctx) bool {
	userAgent := c.Get("User-Agent")
	origin := c.Get("Origin")

	for _, pattern := range ua {
		if strings.Contains(userAgent, pattern) {
			if origin == "" {
				return true
			}
		}
	}

	return false
}
