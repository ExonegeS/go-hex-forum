package port

import (
	"context"

	"github.com/ExonegeS/go-hex-forum/internal/domain/model"
)

type RickMortyAPI interface {
	GetCharacterByID(ctx context.Context, id int) (model.CharacterData, error)
}
