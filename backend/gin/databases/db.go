package databases

import (
	"dagger/backend/gin/models"
	"dagger/backend/gin/runtime"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func init() {
	loglevel := logger.Default.LogMode(logger.Error)
	if runtime.Debug {
		loglevel = logger.Default.LogMode(logger.Info)
	}

	var err error
	// dbMode, _ := runtime.Cfg.GetValue("db", "DB_MODE")
	// if dbMode == "sqlite" {
	// 	path, _ := runtime.Cfg.GetValue("sqlite", "SQLITE_PATH")
	// 	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{
	// 		DisableForeignKeyConstraintWhenMigrating: true,
	// 		Logger:                                   loglevel,
	// 	})
	// } else if dbMode == "mysql" {
	// 	addr, _ := runtime.Cfg.GetValue("mysql", "MYSQL_ADDR")
	// 	DB, err = gorm.Open(mysql.Open(addr), &gorm.Config{})
	// } else {
	// 	log.Panicf("no support database")
	// }
	DB, err = gorm.Open(sqlite.Open("db/dagger.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   loglevel,
	})
	if err != nil {
		log.Panicf("sqlite connect error %v", err)
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
	if !db.Migrator().HasTable(&models.LogSnapshot{}) {
		db.Migrator().CreateTable(&models.LogSnapshot{})
	}
}
