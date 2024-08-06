package request

type GoodsFilter struct {
	PriceMin    int32  `json:"priceMin,omitempty"`
	PriceMax    int32  `json:"priceMax,omitempty"`
	IsHot       bool   `json:"isHot,omitempty"`
	IsNew       bool   `json:"isNew,omitempty"`
	IsTab       bool   `json:"isTab,omitempty"`
	TopCategory int32  `json:"topCategory,omitempty"`
	Pages       int32  `json:"pages,omitempty"`
	PagePerNums int32  `json:"pagePerNums,omitempty"`
	KeyWords    string `json:"keyWords,omitempty"`
	Brand       int32  `json:"brand,omitempty"`
}
