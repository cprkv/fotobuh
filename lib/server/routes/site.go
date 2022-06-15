package routes

import (
	"fotobuh/lib/db"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func SiteMain(c *fiber.Ctx) error {
	cats, err := db.Context.GetCategories()
	if err != nil {
		return err
	}

	return c.Render("index", fiber.Map{
		"Title":      "latest",
		"Categories": catArrToModel(cats),
		"Current":    "/",
	})
}

func SiteCategoryGet(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	cats, err := db.Context.GetCategories()
	if err != nil {
		return err
	}

	cat, err := db.Context.GetCategoryWithPictures(uint(id))
	if err != nil {
		return err
	}

	return c.Render("category", fiber.Map{
		"Title":      cat.Name,
		"Category":   catToModel(&cat),
		"Categories": catArrToModel(cats),
		"Current":    strconv.FormatInt(int64(id), 10),
	})
}
