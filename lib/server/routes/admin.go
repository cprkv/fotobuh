package routes

import (
	"bufio"
	"fmt"
	"fotobuh/lib"
	"fotobuh/lib/db"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jbrodriguez/mlog"
)

type CreateCategoryParams struct {
	Name string `json:"name" xml:"name" form:"name"`
}

func AdminMain(c *fiber.Ctx) error {
	cats, err := db.Context.GetCategories()
	if err != nil {
		return err
	}

	return c.Render("admin/admin", fiber.Map{
		"Title":      "admin page",
		"Categories": catArrToModel(cats),
	})
}

func AdminCategoryManage(c *fiber.Ctx) error {
	cats, err := db.Context.GetCategories()
	if err != nil {
		return err
	}

	return c.Render("admin/category-manage", fiber.Map{
		"Title":      "manage categories",
		"Categories": catArrToModel(cats),
	})
}

func AdminCategoryCreate(c *fiber.Ctx) error {
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
}

func AdminCategoryGet(c *fiber.Ctx) error {
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
}

func AdminCategoryDelete(c *fiber.Ctx) error {
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
}

func AdminPicturesUpload(c *fiber.Ctx) error {
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
