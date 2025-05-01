package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"time"
)

func generatePostID(title, content string) string {
	h := sha256.New()
	h.Write([]byte(title))
	h.Write([]byte(content))
	h.Write([]byte(time.Now().String()))
	sum := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(sum[:16])
}
