package catalog

import (
	"fmt"
	"github.com/moxeed/store/catalog/catalog_model"
	"github.com/moxeed/store/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateProduct(model catalog_model.CreateProductModel) catalog_model.ProductModel {
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
		ReferenceCode: model.ReferenceCode,
		Title:         model.Title,
		Price:         model.Price,
		IsPermanent:   model.IsPermanent,
		Category:      category,
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

func GetProductInfo(productIds []uint) map[uint]catalog_model.ProductModel {
	var products []Product
	common.DB.Preload(clause.Associations).Find(&products, productIds)

	result := make(map[uint]catalog_model.ProductModel)
	for _, product := range products {
		result[product.ID] = product.ToModel()
	}

	return result
}

func GetProductCallBacks(productIds []uint) map[uint]catalog_model.CallBackModel {
	var products []Product
	common.DB.Preload(clause.Associations).Find(&products, productIds)

	result := make(map[uint]catalog_model.CallBackModel)
	for _, product := range products {
		result[product.ID] = catalog_model.CallBackModel{
			ValidationCallBack: product.Category.ValidationCallBack,
			CheckOutCallBack:   product.Category.CheckOutCallBack,
		}
	}

	return result
}

func (p *Product) ToModel() catalog_model.ProductModel {
	return catalog_model.ProductModel{
		ID:            p.ID,
		ReferenceCode: p.ReferenceCode,
		Title:         p.Title,
		Price:         p.Price,
		IsPermanent:   p.IsPermanent,
		CategoryKey:   p.Category.Key,
		CategoryTitle: p.Category.Title,
	}
}
