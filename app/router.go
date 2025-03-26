package app

import (
	"github.com/baeorg/buddy/pkg/manager"
	"github.com/gofiber/fiber/v2"
)

func SetRoute(router *fiber.App) {
	router.Post("/user/im", manager.RegisterUserIM)
}
