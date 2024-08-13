package routes

import (
	"github.com/gofiber/fiber/v2"
	"roastgithub-api/app/handlers"
)

func PublicRoutes(a *fiber.App) {
	route := a.Group("/api/")

	// Get
	a.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"uri":  c.OriginalURL(),
			"path": c.Path(),
		})
	})

	// Post
	route.Post("/roast", handlers.GetRoast)
}
