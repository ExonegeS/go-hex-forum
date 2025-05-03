package rickmorty

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"go-hex-forum/internal/core/domain"
)

type UserDataProvider struct {
	baseURL         string
	client          *http.Client
	availableIDs    []int
	inUse           map[int]time.Time
	mu              sync.Mutex
	randPool        *rand.Rand
	charactersCount int
}

func NewUserDataProvider(baseURL string, charactersCount int) *UserDataProvider {
	// Initialize with shuffled IDs 1-charactersCount
	ids := make([]int, charactersCount)
	for i := range ids {
		ids[i] = i + 1
	}

	p := &UserDataProvider{
		baseURL:         baseURL,
		client:          &http.Client{Timeout: 30 * time.Second},
		availableIDs:    ids,
		inUse:           make(map[int]time.Time),
		randPool:        rand.New(rand.NewSource(time.Now().UnixNano())),
		charactersCount: charactersCount,
	}

	// Initial shuffle
	p.randPool.Shuffle(len(p.availableIDs), func(i, j int) {
		p.availableIDs[i], p.availableIDs[j] = p.availableIDs[j], p.availableIDs[i]
	})

	// Start background cleanup
	go p.cleanupExpiredIDs()
	return p
}

func (p *UserDataProvider) GetUserData(ttl time.Duration) (domain.UserData, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Try to get from available IDs first
	if len(p.availableIDs) > 0 {
		id := p.availableIDs[0]
		p.availableIDs = p.availableIDs[1:]
		p.inUse[id] = time.Now().Add(ttl)
		return p.fetchUser(context.Background(), id)
	}

	// Fallback to random available ID
	for attempts := 0; attempts < 5; attempts++ {
		id := p.randPool.Intn(p.charactersCount) + 1
		if expiry, exists := p.inUse[id]; !exists || time.Now().After(expiry) {
			p.inUse[id] = time.Now().Add(ttl)
			return p.fetchUser(context.Background(), id)
		}
	}

	// If all else fails, return random without checking
	id := p.randPool.Intn(p.charactersCount) + 1
	return p.fetchUser(context.Background(), id)
}

func (p *UserDataProvider) fetchUser(ctx context.Context, id int) (domain.UserData, error) {
	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/character/%d", p.baseURL, id), nil)
	if err != nil {
		return domain.UserData{}, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return domain.UserData{}, err
	}
	defer resp.Body.Close()

	var data struct {
		Id    int64  `json:"id"`
		Name  string `json:"name"`
		Image string `json:"image"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return domain.UserData{}, err
	}

	return domain.UserData{
		ID:        data.Id,
		Name:      data.Name,
		AvatarURL: data.Image,
	}, nil
}

func (p *UserDataProvider) cleanupExpiredIDs() {
	for range time.Tick(5 * time.Minute) {
		p.mu.Lock()
		now := time.Now()
		for id, expiry := range p.inUse {
			if now.After(expiry) {
				delete(p.inUse, id)
				p.availableIDs = append(p.availableIDs, id)
			}
		}
		p.mu.Unlock()
	}
}
