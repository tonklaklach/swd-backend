package repository

import (
	"errors"

	"github.com/mattn/go-sqlite3"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }
func mapError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrNotFound
	}
	var sqlite sqlite3.Error
	if errors.As(err, &sqlite) && sqlite.Code == sqlite3.ErrConstraint {
		if sqlite.ExtendedCode == sqlite3.ErrConstraintForeignKey {
			return domain.ErrInvalidReference
		}
		return domain.ErrConflict
	}
	return err
}
