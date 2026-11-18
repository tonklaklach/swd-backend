package responsestruct

type Success[T any] struct {
	Data T `json:"data"`
}
