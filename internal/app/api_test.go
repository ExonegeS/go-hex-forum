package apiserver

import "testing"

func TestRunApp(t *testing.T) {
	server := NewAPIServer(nil, nil, nil)
	err := server.Run()
	if err == nil {
		t.Fatal("Expected error, got no errors")
	}
}
