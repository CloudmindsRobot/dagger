package databases

import (
	"dagger/backend/runtime"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func init() {
	debug, _ := runtime.Cfg.Bool("global", "debug")
	loglevel := logger.Default.LogMode(logger.Error)
	if debug {
		loglevel = logger.Default.LogMode(logger.Info)
	}

	var err error

	addr, _ := runtime.Cfg.GetValue("db", "address")
	DB, err = gorm.Open(mysql.Open(addr), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   loglevel,
	})

	if err != nil {
		log.Panicf("db connect error %v", err)
	}
	if DB.Error != nil {
		log.Panicf("database error %v", DB.Error)
	}
}
