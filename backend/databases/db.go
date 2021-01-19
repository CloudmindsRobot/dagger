package databases

import (
	"dagger/backend/models"
	"dagger/backend/runtime"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
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
	dbMode, _ := runtime.Cfg.GetValue("db", "mode")
	if dbMode == "sqlite" {
		DB, err = gorm.Open(sqlite.Open("db/dagger.db"), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   loglevel,
		})
	} else if dbMode == "mysql" {
		addr, _ := runtime.Cfg.GetValue("db", "address")
		DB, err = gorm.Open(mysql.Open(addr), &gorm.Config{})
	} else {
		log.Panicf("no support database")
	}
	if err != nil {
		log.Panicf("db connect error %v", err)
	}
	if DB.Error != nil {
		log.Panicf("database error %v", DB.Error)
	}
}

func MigrateDB(db *gorm.DB) {
	if !db.Migrator().HasTable(&models.User{}) {
		db.Migrator().CreateTable(&models.User{})
		username, _ := runtime.Cfg.GetValue("users", "admin_username")
		password, _ := runtime.Cfg.GetValue("users", "admin_passwod")
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Panicf("init admin user error %v", err)
		}
		db.Model(&models.User{}).Create(&models.User{
			Username:    username,
			Password:    string(hash),
			IsActive:    true,
			IsSuperuser: true,
			CreateAt:    time.Now(),
		})
	}
	if !db.Migrator().HasTable(&models.LogHistory{}) {
		db.Migrator().CreateTable(&models.LogHistory{})
	}
	if !db.Migrator().HasTable(&models.LogRule{}) {
		db.Migrator().CreateTable(&models.LogRule{})
	}
	if !db.Migrator().HasTable(&models.LogSnapshot{}) {
		db.Migrator().CreateTable(&models.LogSnapshot{})
	}
	if !db.Migrator().HasTable(&models.LogUserGroup{}) {
		db.Migrator().CreateTable(&models.LogUserGroup{})
	}
	if !db.Migrator().HasTable(&models.LogUser{}) {
		db.Migrator().CreateTable(&models.LogUser{})
	}
	if !db.Migrator().HasTable(&models.LogEvent{}) {
		db.Migrator().CreateTable(&models.LogEvent{})
	}
	if !db.Migrator().HasTable(&models.LogLabel{}) {
		db.Migrator().CreateTable(&models.LogLabel{})
	}
	if !db.Migrator().HasTable(&models.LogGroup{}) {
		db.Migrator().CreateTable(&models.LogGroup{})
	}
	if !db.Migrator().HasTable(&models.LogEventDetail{}) {
		db.Migrator().CreateTable(&models.LogEventDetail{})
	}
}
