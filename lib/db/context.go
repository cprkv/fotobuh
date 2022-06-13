package db

import (
	"fotobuh/lib"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseContext struct {
	*gorm.DB
}

var (
	Context = DatabaseContext{}
)

func (self *DatabaseContext) Init() error {
	nowFunc := func() time.Time {
		return time.Now().UTC()
	}
	logger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Millisecond * 100,
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	config := &gorm.Config{
		NowFunc: nowFunc,
		Logger:  logger,
	}

	dialector := driverOpen(lib.Config.Database.Connection)

	var err error
	self.DB, err = gorm.Open(dialector, config)
	if err != nil {
		return err
	}

	err = self.AutoMigrate(
		&Picture{},
		&Category{},
	)
	if err != nil {
		return err
	}

	log.Printf("database initialized: %v", lib.Config.Database.Connection)
	return nil
}
