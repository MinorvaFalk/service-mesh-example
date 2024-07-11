package handler

import (
	"time"
	"worker-mesh/config"
	"worker-mesh/internal/model"
	"worker-mesh/pkg/messaging"

	"github.com/gofiber/fiber/v2"
)

const (
	NotificationContext = "notification"
)

type HTTPHandler interface {
	ParseRequestBody(c *fiber.Ctx) error
	Publish(c *fiber.Ctx) error
	PublishDeferred(c *fiber.Ctx) error
}

type httpHandler struct {
	producer messaging.Producer
}

func NewHttpHandler(producer messaging.Producer) HTTPHandler {
	return &httpHandler{producer}
}

func (h *httpHandler) ParseRequestBody(c *fiber.Ctx) error {
	var notification Notification
	if err := c.BodyParser(&notification); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(model.HTTPResponse{
				Code:    fiber.StatusBadRequest,
				Status:  "Bad Request",
				Message: "invalid request body",
			})
	}

	c.Locals(NotificationContext, notification)

	return c.Next()
}

func (h *httpHandler) Publish(c *fiber.Ctx) error {
	notification := c.Locals(NotificationContext).(Notification)

	if err := h.producer.Publish(config.ReadConfig().Nsq.Topic, notification.Byte()); err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(model.HTTPResponse{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "failed to publish message",
			})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *httpHandler) PublishDeferred(c *fiber.Ctx) error {
	notification := c.Locals(NotificationContext).(Notification)
	delay, err := c.ParamsInt("delay", 1000)
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(model.HTTPResponse{
				Code:    fiber.StatusBadRequest,
				Status:  "Bad Request",
				Message: "invalid delay parameter",
			})
	}

	if err := h.producer.DeferredPublish(
		config.ReadConfig().Nsq.Topic,
		time.Millisecond*time.Duration(delay),
		notification.Byte(),
	); err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(model.HTTPResponse{
				Code:    fiber.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "failed to publish deferred message",
			})
	}

	return c.SendStatus(fiber.StatusOK)
}
