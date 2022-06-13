//go:build windows
// +build windows

package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func driverOpen(con string) gorm.Dialector {
	return sqlite.Open(con)
}
