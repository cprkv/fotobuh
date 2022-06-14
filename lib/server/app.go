package server

import (
	"bufio"
	"fmt"
	"fotobuh/lib"
	"fotobuh/lib/db"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/amber"
	"github.com/google/uuid"
	"github.com/jbrodriguez/mlog"
)

type App struct {
}

type CreateCategoryParams struct {
	Name string `json:"name" xml:"name" form:"name"`
}

func StartApp() {
	viewEngine := amber.New("./views", ".amber").Debug(true)

	app := fiber.New(fiber.Config{
		Views:     viewEngine,
		BodyLimit: 256 * 1024 * 1024,
	})
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Use("static", filesystem.New(filesystem.Config{
		Root: http.Dir("./static"),
	}))

	app.Use("images", filesystem.New(filesystem.Config{
		Root: http.Dir(lib.Config.Pictures.Storage),
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		cats, err := db.Context.GetCategories()
		if err != nil {
			return err
		}

		return c.Render("index", fiber.Map{
			"Title":      "latest",
			"Categories": catArrToModel(cats),
			"Current":    "/",
		})
	})

	app.Get("/category/:id", func(c *fiber.Ctx) error {
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
	})

	app.Get("/admin", func(c *fiber.Ctx) error {
		cats, err := db.Context.GetCategories()
		if err != nil {
			return err
		}

		return c.Render("admin/admin", fiber.Map{
			"Title":      "admin page",
			"Categories": catArrToModel(cats),
		})
	})

	app.Get("/admin/category/manage", func(c *fiber.Ctx) error {
		cats, err := db.Context.GetCategories()
		if err != nil {
			return err
		}

		return c.Render("admin/category-manage", fiber.Map{
			"Title":      "manage categories",
			"Categories": catArrToModel(cats),
		})
	})

	app.Post("/admin/category/create", func(c *fiber.Ctx) error {
		p := new(CreateCategoryParams)
		if err := c.BodyParser(p); err != nil {
			return err
		}

		id, err := db.Context.CreateCategory(p.Name)
		if err != nil {
			// todo: error?
			return c.Redirect("/admin/category/manage")
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

		cats, err := db.Context.GetCategories()
		if err != nil {
			return err
		}

		return c.Render("admin/category", fiber.Map{
			"Title":      "category: " + cat.Name,
			"Category":   catToModel(&cat),
			"Categories": catArrToModel(cats),
		})
	})

	app.Post("/admin/category/:id/delete", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}

		cat, err := db.Context.GetCategory(uint(id))
		if err != nil {
			return err
		}

		err = db.Context.DeleteCategory(&cat)
		if err != nil {
			return err
		}

		return c.Redirect("/admin/category/manage")
	})

	app.Post("/admin/category/:id/upload", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}

		cat, err := db.Context.GetCategory(uint(id))
		if err != nil {
			return err
		}

		form, err := c.MultipartForm()
		if err != nil {
			return err
		}

		for formFieldName, fileHeaders := range form.File {
			mlog.Trace("file %v", formFieldName)
			for _, header := range fileHeaders {
				mlog.Trace("  header Filename: %v", header.Filename)
				mlog.Trace("  header Header: %v", header.Header)
				mlog.Trace("  header Size: %v", header.Size)
				err = uploadFile(header, &cat)
				if err != nil {
					return err
				}
			}
		}

		return c.Redirect(fmt.Sprintf("/admin/category/%d", id))
	})

	// app.Post("/api/category/:name", func(c *fiber.Ctx) error {
	// 	name := c.Params("name")
	// 	if len(name) == 0 {
	// 		return c.JSON(fiber.Map{"error": "category name not defined"})
	// 	}
	// 	return c.JSON(fiber.Map{"message": "Hello World"})
	// })

	err := app.Listen(lib.Config.Http.Address)
	if err != nil {
		mlog.Fatal(err)
	}
}

func uploadFile(header *multipart.FileHeader, category *db.Category) error {
	extension := filepath.Ext(header.Filename)
	name := strings.TrimSuffix(header.Filename, extension)
	fileName := uuid.New().String() + extension
	resultPath := filepath.Join(lib.Config.Pictures.Storage, fileName)
	picture := &db.Picture{}

	handle, err := header.Open()
	if err != nil {
		return err
	}
	defer handle.Close()

	outFile, err := os.Create(resultPath)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(handle)
	bytesWritten, err := reader.WriteTo(outFile)
	outFile.Close()

	if err != nil {
		err = fmt.Errorf("error writing file %v: %v", outFile, err)
		goto exitError
	}

	if bytesWritten != header.Size {
		err = fmt.Errorf("error writing file %v: written %v bytes, expected %v bytes",
			outFile, bytesWritten, header.Size)
		goto exitError
	}

	picture.Name = name
	picture.FileName = fileName
	picture.CreatedAt = time.Now()
	picture.Categories = []*db.Category{category}

	err = db.Context.CreatePicture(picture)
	if err != nil {
		goto exitError
	}

	mlog.Trace("created picture: %v", picture.ID)
	return nil

exitError:
	os.Remove(resultPath)
	return err
}

func init() {
	mlog.StartEx(mlog.LevelTrace, "fotobuh.log", 5*1024*1024, 5)
}
