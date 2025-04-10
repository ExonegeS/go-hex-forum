package rick_morty

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ExonegeS/go-hex-forum/internal/domain/model"
)

type Rick_MortyClient struct {
}

func (c *Rick_MortyClient) GetCharacterByID(ctx context.Context, id int) (*model.CharacterData, error) {
	url := fmt.Sprintf("https://rickandmortyapi.com/api/character/%d", id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get character: %s", resp.Status)
	}

	var character model.CharacterData
	if err := json.NewDecoder(resp.Body).Decode(&character); err != nil {
		return nil, err
	}

	return &character, nil
}
