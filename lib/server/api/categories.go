package api

import (
	"fotobuh/lib/db"

	"github.com/gofiber/fiber/v2"
)

func GetCategories(c *fiber.Ctx) error {
	cats, err := db.Context.GetCategories()
	if err != nil {
		return internalServerError(c, err)
	}

	return c.JSON(catArrToModel(cats))
}

func GetCategory(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return badRequest(c, "param id is invalid or missing")
	}

	cat, err := db.Context.GetCategoryWithPictures(uint(id))
	if err != nil {
		return internalServerError(c, err)
	}

	return c.JSON(catToModel(&cat))
}
