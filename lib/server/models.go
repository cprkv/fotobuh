package server

import (
	"fotobuh/lib/db"
	"strconv"
)

func picToModel(pic *db.Picture) map[string]interface{} {
	r := make(map[string]interface{})
	r["name"] = pic.Name
	r["fileName"] = pic.FileName
	r["createdAt"] = pic.CreatedAt.Format("2006.01.02 15:04")
	return r
}

func catToModel(cat *db.Category) map[string]interface{} {
	r := make(map[string]interface{})
	r["id"] = strconv.FormatUint(uint64(cat.ID), 10)
	r["name"] = cat.Name
	r["createdAt"] = cat.CreatedAt.Format("2006.01.02 15:04")

	pictures := make([]map[string]interface{}, len(cat.Pictures))
	for i := 0; i < len(cat.Pictures); i++ {
		pictures[i] = picToModel(cat.Pictures[i])
	}
	r["pictures"] = pictures

	return r
}

func catArrToModel(cat []db.Category) []map[string]interface{} {
	r := make([]map[string]interface{}, len(cat))
	for i := 0; i < len(cat); i++ {
		r[i] = catToModel(&cat[i])
	}
	return r
}
