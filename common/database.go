package common

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func init() {
	dsn := Configuration.DataBaseDsn
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	if err != nil {
		panic(err)
	}

	DB = db
}

func AutoMigrate(entity ...interface{}) {
	if err := DB.AutoMigrate(entity...); err != nil {
		panic(err)
	}
}
