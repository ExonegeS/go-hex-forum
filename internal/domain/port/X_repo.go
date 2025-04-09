package port

import (
	"context"

	"github.com/ExonegeS/go-hex-forum/internal/domain/model"
)

type XRepository interface {
	Save(context.Context, model.X) (int, error)
}
