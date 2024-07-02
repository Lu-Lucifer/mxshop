package code

//go:generate codegen -type=int
const (
	// ErrGoodsNotFound - 404: Goods not found.
	ErrGoodsNotFound int = iota + 100501
	// ErrCategoryNotFound - 404: Category not found.
	ErrCategoryNotFound
	// ErrEsUnmarshal - 500: Es unmarshal error.
	ErrEsUnmarshal
)
