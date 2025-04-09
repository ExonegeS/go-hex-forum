package application

import (
	"context"

	"github.com/ExonegeS/go-hex-forum/internal/domain/model"
	"github.com/ExonegeS/go-hex-forum/internal/domain/port"
)

type XService struct {
	Xrepo port.XRepository
}

func NewXService(Xrepo port.XRepository) *XService {
	return &XService{
		Xrepo,
	}
}

func (s *XService) CreateX(ctx context.Context, data string) (int, error) {
	x := model.X{
		Data: data,
	}
	id, err := s.Xrepo.Save(ctx, x)
	if err != nil {
		return -1, err
	}
	return id, nil
}
