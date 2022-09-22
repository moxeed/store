package catalog_model

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
