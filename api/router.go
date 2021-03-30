package api

import (
	"log"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/Mshivam2409/AnC-Courses/graph"
	"github.com/Mshivam2409/AnC-Courses/graph/generated"
	"github.com/Mshivam2409/AnC-Courses/services"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// SetupRoutes ....
func SetupRoutes(app *fiber.App) {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	cache, err := services.NewCache(viper.GetString("redis.host"), viper.GetString("redis.pwd"), 24*time.Hour)
	if err != nil {
		log.Printf("cannot create APQ redis cache: %v", err)
	}
	srv.Use(extension.AutomaticPersistedQuery{Cache: cache})
	gqlHandler := srv.Handler()
	app.All("/graphql", func(c *fiber.Ctx) error {
		a := c.Get("Authorizaton")
		log.Print(a)
		c.Locals("Authorizaton", a)
		gqlHandler(c.Context())
		return nil
	})
	restAPI := app.Group("/secure")
	restAPI.Post("/file/create", CreateFile)
	app.Get("/file/:fid", DownloadFile)
	app.All("/hydrate", Hydrate)
	app.Post("/register", Register)
}
