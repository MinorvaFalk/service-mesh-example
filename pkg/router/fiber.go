package router

import (
	"worker-mesh/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func NewFiber(handler handler.HTTPHandler) *fiber.App {
	app := fiber.New()

	publish := app.Group("/publish", handler.ParseRequestBody)
	publish.Post("/", handler.Publish)
	publish.Post("/defer/:delay", handler.PublishDeferred)

	return app
}
