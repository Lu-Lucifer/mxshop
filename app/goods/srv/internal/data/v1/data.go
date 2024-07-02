package v1

import "gorm.io/gorm"

// 依赖注入-工厂模式，建立一个data层的store接口的工厂
type DataFactory interface {
	Goods() GoodsStore
	Category() CategoryStore
	Brands() BrandsStore
	Banners() BannerStore
	CategoryBrands() GoodsCategoryBrandStore
	Begin() *gorm.DB
}
