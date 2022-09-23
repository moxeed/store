package catalog_model

type CallBackModel struct {
	ValidationCallBack string `json:"validationCallBack,omitempty"`
	CheckOutCallBack   string `json:"checkOutCallBack,omitempty"`
}

type CreateProductModel struct {
	ReferenceCode      uint
	Title              string `json:"title,omitempty"`
	Price              uint   `json:"price,omitempty"`
	IsPermanent        bool   `json:"isPermanent,omitempty"`
	CategoryKey        string `json:"categoryKey,omitempty"`
	CategoryTitle      string `json:"categoryTitle,omitempty"`
	ValidationCallBack string `json:"validationCallBack,omitempty"`
	CheckOutCallBack   string `json:"checkOutCallBack,omitempty"`
}

type ProductModel struct {
	ID            uint `json:"ID,omitempty"`
	ReferenceCode uint
	Title         string `json:"title,omitempty"`
	Price         uint   `json:"price,omitempty"`
	IsPermanent   bool   `json:"isPermanent,omitempty"`
	CategoryKey   string `json:"categoryKey,omitempty"`
	CategoryTitle string `json:"categoryTitle,omitempty"`
}
