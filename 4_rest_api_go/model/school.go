package model

type SchoolRequest struct {
	ID           int64  `json:"id"`
	Name         string `json:"name" validate:"required,notblank"`
	Abbreviation string `json:"abbreviation" validate:"required,notblank"`
	Address      string `json:"address" validate:"required,notblank"`
}
