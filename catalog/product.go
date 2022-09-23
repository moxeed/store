package catalog

import (
	"github.com/moxeed/store/common"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Key                string
	Title              string
	ValidationCallBack string
	CheckOutCallBack   string
}

type Product struct {
	gorm.Model
	ReferenceCode uint
	Title         string
	Price         uint
	IsPermanent   bool
	CategoryID    uint
	Category      Category
}

func init() {
	common.AutoMigrate(
		&Product{},
		&Category{})
}
