package api

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func FiberListenAndServe() error {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}
			err = ctx.Status(code).JSON(fiber.Map{"message": message})
			if err != nil {
				return ctx.Status(500).SendString("Internal Server Error")
			}
			return nil
		},
	})
	app.Use(cors.New(cors.Config{AllowOrigins: "http://127.0.0.1:3000"}))
	app.Use(recover.New())
	prometheus := fiberprometheus.New("fiber")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)
	SetupRoutes(app)
	app.Static("/", "html")
	err := app.Listen(":5000")
	if err != nil {
		return err
	}
	return nil
}
