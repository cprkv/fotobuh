package server

import (
	"fmt"
	"fotobuh/lib"
	"fotobuh/lib/db"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/amber"
)

type App struct {
}

type CreateCategoryParams struct {
	Name string `json:"name" xml:"name" form:"name"`
}

func StartApp() {
	viewEngine := amber.New("./views", ".amber").Debug(true)
	// viewEngine := pug.New("./views", ".pug").Debug(true)

	app := fiber.New(fiber.Config{
		Views: viewEngine,
	})

	app.Use("static", filesystem.New(filesystem.Config{
		Root: http.Dir("./static"),
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{"Title": "ABOBA"})
	})

	app.Get("/admin/category/create", func(c *fiber.Ctx) error {
		return c.Render("admin/category-create", fiber.Map{"Title": "create category"})
	})

	app.Get("/admin/category/manage", func(c *fiber.Ctx) error {
		cat, err := db.Context.GetCategories()
		if err != nil {
			return err
		}

		return c.Render("admin/category-manage", fiber.Map{
			"Title":      "manage categories",
			"Categories": catArrToModel(cat),
		})
	})

	app.Post("/admin/category/create", func(c *fiber.Ctx) error {
		p := new(CreateCategoryParams)
		if err := c.BodyParser(p); err != nil {
			return err
		}

		id, err := db.Context.CreateCategory(p.Name)
		if err != nil {
			return c.Render("admin/category-create", fiber.Map{
				"Title": "create category",
				"Error": err.Error(),
			})
		}

		return c.Redirect(fmt.Sprintf("/admin/category/%d", id))
	})

	app.Get("/admin/category/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}

		cat, err := db.Context.GetCategory(uint(id))
		if err != nil {
			return err
		}

		return c.Render("admin/category", fiber.Map{
			"Title":    "category: " + cat.Name,
			"Category": catToModel(&cat),
		})
	})

	// app.Post("/api/category/:name", func(c *fiber.Ctx) error {
	// 	name := c.Params("name")
	// 	if len(name) == 0 {
	// 		return c.JSON(fiber.Map{"error": "category name not defined"})
	// 	}
	// 	return c.JSON(fiber.Map{"message": "Hello World"})
	// })

	log.Fatal(app.Listen(lib.Config.Http.Address))
}
