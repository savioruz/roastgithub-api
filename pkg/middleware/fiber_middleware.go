package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func FiberMiddleware(a *fiber.App) {
	a.Use(logger.New())
	a.Use(recover.New())
	a.Use(cors.New(cors.Config{
		AllowOrigins: "https://*.savioruz.me, https://roastgithub-api.vercel.app",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST",
	}))
}
