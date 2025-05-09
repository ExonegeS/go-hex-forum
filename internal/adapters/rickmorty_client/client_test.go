package rickmorty

import (
	"fmt"
	"testing"
	"time"
)

func TestGetUserData_Success(t *testing.T) {
	client := NewUserDataProvider("https://rickandmortyapi.com/api", 5)
	if client == nil {
		t.Fatalf("UserDataProvider returned nil")
	}
	userdata, err := client.GetUserData(time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	fmt.Println(userdata)
}

func TestGetUserData_Success_Unique(t *testing.T) {
	usersLen := 10
	client := NewUserDataProvider("https://rickandmortyapi.com/api", usersLen)
	if client == nil {
		t.Fatalf("UserDataProvider returned nil")
	}
	users := make(map[string]int, usersLen)
	for i := 0; i < usersLen; i++ {
		userdata, err := client.GetUserData(time.Hour)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, ok := users[userdata.AvatarURL]; ok {
			t.Fatalf("not unique avatar, avatar: %v map: %v", userdata.AvatarURL, users)
		}
		users[userdata.AvatarURL]++
		fmt.Println(userdata)
	}
}
