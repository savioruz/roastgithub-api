package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savioruz/roastgithub-api/app/handlers"
)

func PublicRoutes(a *fiber.App) {
	route := a.Group("/api/v1/")

	// Get
	a.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger")
	})

	// Post
	route.Post("/roast/github", handlers.GetGithubRoast)
	route.Post("/roast/resume", handlers.GetResumeRoast)
}
