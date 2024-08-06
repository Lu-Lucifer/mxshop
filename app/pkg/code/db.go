package code

const (
	// ErrConnectDB - 500: Connect db error.
	ErrConnectDB int = iota + 100601
	// ErrConnectGRPC - 500: Connect grpc error.
	ErrConnectGRPC
)
