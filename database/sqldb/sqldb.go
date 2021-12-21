package sqldb

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLiteGetter() func(string) *gorm.DB {

	cache := make(map[string]*gorm.DB)

	return func(path string) *gorm.DB {

		if db, ok := cache[path]; ok {
			return db
		}

		db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		// Migrate the schema
		db.AutoMigrate(&User{})
		db.AutoMigrate(&Product{})
		db.AutoMigrate(&OrderItem{})
		db.AutoMigrate(&Order{})

		cache[path] = db

		return db
	}

}
