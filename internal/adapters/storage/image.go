package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type ImageStorage struct {
	baseURL    string
	client     *http.Client
	codeLength int64
}

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewImageStorage(baseURL string, codeLength int64) *ImageStorage {
	return &ImageStorage{
		baseURL:    baseURL,
		client:     &http.Client{Timeout: 30 * time.Second},
		codeLength: codeLength,
	}
}

func (s *ImageStorage) generateCode() string {
	b := make([]byte, s.codeLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (s *ImageStorage) UploadImage(ctx context.Context, userID int64, data []byte) (string, error) {
	code := s.generateCode()
	bucketPath := fmt.Sprintf("user-%d", userID)
	objectPath := fmt.Sprintf("%s/%s", bucketPath, code)
	base := s.baseURL

	put := func(path string, body io.Reader) (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, "PUT", base+"/"+path, body)
		if err != nil {
			return nil, err
		}
		return s.client.Do(req)
	}

	resp, err := put(objectPath, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		bResp, bErr := put(bucketPath, nil)
		if bErr != nil {
			return "", fmt.Errorf("failed to create bucket %q: %w", bucketPath, bErr)
		}
		bResp.Body.Close()
		if bResp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("bucket creation error: %d", bResp.StatusCode)
		}

		resp, err = put(objectPath, bytes.NewReader(data))
		if err != nil {
			return "", fmt.Errorf("retry upload failed: %w", err)
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("storage error: %d", resp.StatusCode)
	}

	return s.GetImageURL(userID, code), nil
}

func (s *ImageStorage) GetImageURL(userID int64, code string) string {
	return fmt.Sprintf("%s/user-%d/%s", s.baseURL, userID, code)
}
