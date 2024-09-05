package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	_ "github.com/savioruz/roastgithub-api/docs" // load API Docs files (Swagger)
	"github.com/savioruz/roastgithub-api/pkg/middleware"
	"github.com/savioruz/roastgithub-api/pkg/routes"
	"net/http"
)

// initializeHandler is a function to initialize the handler
func initializeHandler() http.HandlerFunc {
	app := fiber.New()

	middleware.FiberMiddleware(app)
	middleware.LimiterMiddleware(app)
	middleware.MonitorMiddleware(app)

	routes.PublicRoutes(app)
	routes.SwaggerRoute(app)
	routes.NotFoundRoute(app)

	return adaptor.FiberApp(app)
}

// Handler is a function to handle the request from fiber app
// @title Roast GitHub API
// @version 0.1
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email jakueenak@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func Handler(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = r.URL.String()

	initializeHandler().ServeHTTP(w, r)
}
