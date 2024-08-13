package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func MonitorMiddleware(a *fiber.App) {
	a.Get("/metrics", monitor.New(
		monitor.Config{
			Title: "Roast Github API Monitor",
		},
	))
}
