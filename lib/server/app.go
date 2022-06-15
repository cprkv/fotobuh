package server

import (
	"fotobuh/lib"
	"fotobuh/lib/server/api"
	"fotobuh/lib/server/routes"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/amber"
	"github.com/jbrodriguez/mlog"
)

func StartApp() {
	fiberConfig := fiber.Config{
		Views:     amber.New("./views", ".amber").Debug(true),
		BodyLimit: 256 * 1024 * 1024,
	}
	app := fiber.New(fiberConfig)

	fiberLogConfig := logger.Config{Format: "[${ip}]:${port} ${status} - ${method} ${path}\n"}
	app.Use(logger.New(fiberLogConfig))

	staticFSConfig := filesystem.Config{Root: http.Dir("./static")}
	app.Use("static", filesystem.New(staticFSConfig))

	imagesFSConfig := filesystem.Config{Root: http.Dir(lib.Config.Pictures.Storage)}
	app.Use("images", filesystem.New(imagesFSConfig))

	app.Use(cors.New())

	app.Get("/", routes.SiteMain)
	app.Get("/category/:id", routes.SiteCategoryGet)
	app.Get("/admin", routes.AdminMain)
	app.Get("/admin/category/manage", routes.AdminCategoryManage)
	app.Post("/admin/category/create", routes.AdminCategoryCreate)
	app.Get("/admin/category/:id", routes.AdminCategoryGet)
	app.Post("/admin/category/:id/delete", routes.AdminCategoryDelete)
	app.Post("/admin/category/:id/upload", routes.AdminPicturesUpload)

	app.Get("/api/category", api.GetCategories)
	app.Get("/api/category/:id", api.GetCategory)

	err := app.Listen(lib.Config.Http.Address)
	if err != nil {
		mlog.Fatal(err)
	}
}

func init() {
	mlog.StartEx(mlog.LevelTrace, "fotobuh.log", 5*1024*1024, 5)
}
