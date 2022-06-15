package api

import "github.com/gofiber/fiber/v2"

func internalServerError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).
		SendString(errToModel(err))
}

func badRequest(c *fiber.Ctx, what string) error {
	return c.Status(fiber.StatusBadRequest).
		SendString(what)
}
