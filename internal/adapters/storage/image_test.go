package storage

import (
	"context"
	"net/url"
	"testing"
)

func TestUploadImage_Success(t *testing.T) {
	var codeLen int = 3
	service := NewImageStorage("http://localhost:6969", int64(codeLen))
	publicURL, err := service.UploadImage(context.Background(), 0, []byte("X"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := url.Parse(publicURL); err != nil {
		t.Fatalf("expected no url parsing err , got %v", err)
	}
}

func TestUploadImage_Fail(t *testing.T) {
	var codeLen int = 3
	service := NewImageStorage("http://localhost:0001", int64(codeLen))
	_, err := service.UploadImage(context.Background(), 0, []byte("X"))
	if err == nil {
		t.Fatalf("expected error, got no errors")
	}
}

func TestGetImageURL_Success(t *testing.T) {
	var codeLen int = 3
	service := NewImageStorage("http://localhost:6969", int64(codeLen))
	publicURL := service.GetImageURL(0, "wKj")
	if _, err := url.Parse(publicURL); err != nil {
		t.Fatalf("expected no url parsing err , got %v", err)
	}
}

func TestGetImageURL_Fail(t *testing.T) {
	var codeLen int = 3
	service := NewImageStorage("http://localhost:6969", int64(codeLen))
	publicURL := service.GetImageURL(0, "___")
	if _, err := url.Parse(publicURL); err != nil {
		t.Fatalf("expected no url parsing err , got %v", err)
	}
}
