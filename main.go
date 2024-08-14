package main

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
	"os"
	_ "roastgithub-api/docs" // load API Docs files (Swagger)
	"roastgithub-api/pkg/middleware"
	"roastgithub-api/pkg/routes"
	"roastgithub-api/pkg/utils"
)

// @title Roast GitHub API
// @version 0.1
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email jakueenak@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	app := fiber.New()

	middleware.FiberMiddleware(app)
	middleware.LimiterMiddleware(app)
	middleware.MonitorMiddleware(app)

	routes.PublicRoutes(app)
	routes.SwaggerRoute(app)

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
