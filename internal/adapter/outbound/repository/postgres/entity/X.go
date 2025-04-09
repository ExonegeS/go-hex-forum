package entity

import (
	"time"

	"github.com/ExonegeS/go-hex-forum/internal/domain/model"
)

type X struct {
	ID        int64     `json:"id" db:"id"`
	Data      string    `json:"data" db:"data"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func XToEntity(x *model.X) *X {
	return &X{
		ID:        int64(x.ID),
		Data:      x.Data,
		CreatedAt: x.CreatedAt,
		UpdatedAt: x.UpdatedAt,
	}
}

func XToModel(x *X) *model.X {
	return &model.X{
		ID:        int(x.ID),
		Data:      x.Data,
		CreatedAt: x.CreatedAt,
		UpdatedAt: x.UpdatedAt,
	}
}
