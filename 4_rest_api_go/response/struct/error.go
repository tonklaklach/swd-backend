package responsestruct

type Error struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}
