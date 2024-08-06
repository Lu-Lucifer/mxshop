package code

//go:generate codeg -type=int -doc -output ./error_code_generated.md
const (
	// ErrShopCartItemNotFound - 404: ShopCart item not found.
	ErrShopCartItemNotFound int = iota + 100701

	// ErrSubmitOrder - 400: Submit order error.
	ErrSubmitOrder

	// ErrOrderNotFound - 404: No Goods selected.
	ErrNoGoodsSelect
)
