package manager

import (
	"log/slog"
	"net/http"

	"github.com/baeorg/buddy/pkg/helper"
	"github.com/baeorg/buddy/pkg/storage"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserIM(c *fiber.Ctx) error {
	var req User
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Bad Request:" + err.Error())
	}

	slog.Info("RegisterUserIM", "req", req)
	// Validate the request
	if err := helper.Validate.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).SendString("Bad Request: " + err.Error())
	}

	// Register the user
	err = storage.UpdatePermission(req.ID, req.Token)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error:" + err.Error())
	}

	return c.SendStatus(http.StatusCreated)
}

type User struct {
	ID    uint64 `json:"id" validate:"required"`
	Token string `json:"token" validate:"required"`
}
