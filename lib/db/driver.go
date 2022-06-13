//go:build !windows
// +build !windows

package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func driverOpen(con string) gorm.Dialector {
	return sqlite.Open(con)
}
