package catalog

import (
	"fmt"
	"github.com/moxeed/store/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CallBackModel struct {
	ValidationCallBack string
	CheckOutCallBack   string
}

type CreateProductModel struct {
	Title              string
	Price              uint
	IsPermanent        bool
	CategoryKey        string
	CategoryTitle      string
	ValidationCallBack string
	CheckOutCallBack   string
}

type ProductModel struct {
	ID            uint
	Title         string
	Price         uint
	IsPermanent   bool
	CategoryKey   string
	CategoryTitle string
}

func CreateProduct(model CreateProductModel) ProductModel {
	category := Category{Key: model.CategoryKey}
	result := common.DB.
		Where(&category).
		First(&category)

	if result.Error == gorm.ErrRecordNotFound {
		category.Title = model.CategoryTitle
		category.ValidationCallBack = model.ValidationCallBack
		category.CheckOutCallBack = model.CheckOutCallBack
	}

	product := Product{
		Title:       model.Title,
		Price:       model.Price,
		IsPermanent: model.IsPermanent,
		Category:    category,
	}

	common.DB.Create(&product)

	return product.ToModel()
}

func GetProduct(id uint) (Product, error) {
	product := Product{}
	dbResult := common.DB.Preload(clause.Associations).First(&product, id)

	if dbResult.Error == gorm.ErrRecordNotFound {
		return product, fmt.Errorf("محصول پیدا نشد")
	}

	return product, nil
}

func GetProductCallBacks(productIds []uint) map[uint]CallBackModel {
	var products []Product
	common.DB.Preload(clause.Associations).Find(&products, productIds)

	result := make(map[uint]CallBackModel)
	for _, product := range products {
		result[product.ID] = CallBackModel{
			ValidationCallBack: product.Category.ValidationCallBack,
			CheckOutCallBack:   product.Category.CheckOutCallBack,
		}
	}

	return result
}

func (p *Product) ToModel() ProductModel {
	return ProductModel{
		ID:            p.ID,
		Title:         p.Title,
		Price:         p.Price,
		IsPermanent:   p.IsPermanent,
		CategoryKey:   p.Category.Key,
		CategoryTitle: p.Category.Title,
	}
}
