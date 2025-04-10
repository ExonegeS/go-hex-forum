package unit

import (
	"context"
	"fmt"
	"testing"

	"github.com/ExonegeS/go-hex-forum/internal/adapter/outbound/API/rick_morty"
)

// -----------------------------------------------------------
// Success
func TestGetCharacter_5_Fail(t *testing.T) {
	client := rick_morty.Rick_MortyClient{}

	id := 1
	characterData, err := client.GetCharacterByID(context.Background(), id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if characterData == nil {
		t.Fatalf("expected character data, got nil")
	}
	if characterData.ID != id {
		t.Fatalf("expected character ID %d, got %d", id, characterData.ID)
	}
	fmt.Println("Character data:", characterData)
}

// -----------------------------------------------------------

// -----------------------------------------------------------
// Edge Case
func TestGetCharacter_1_Edge(t *testing.T) {
	client := rick_morty.Rick_MortyClient{}

	id := 1
	characterData, err := client.GetCharacterByID(context.Background(), id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if characterData == nil {
		t.Fatalf("expected character data, got nil")
	}
	if characterData.ID != id {
		t.Fatalf("expected character ID %d, got %d", id, characterData.ID)
	}
	fmt.Println("Character data:", characterData)
}

func TestGetCharacter_826_Edge(t *testing.T) {
	client := rick_morty.Rick_MortyClient{}

	id := 826
	characterData, err := client.GetCharacterByID(context.Background(), id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if characterData == nil {
		t.Fatalf("expected character data, got nil")
	}
	if characterData.ID != id {
		t.Fatalf("expected character ID %d, got %d", id, characterData.ID)
	}
	fmt.Println("Character data:", characterData)
}

// -----------------------------------------------------------

// -----------------------------------------------------------
// Failure
func TestGetCharacter_0_Fail(t *testing.T) {
	client := rick_morty.Rick_MortyClient{}

	characterData, err := client.GetCharacterByID(context.Background(), 0)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if characterData != nil {
		t.Fatalf("expected nil character data, got %v", characterData)
	}
	fmt.Println("Error:", err)
	fmt.Println("Character data:", characterData)

	if err.Error() != "failed to get character: 404 Not Found" {
		t.Fatalf("expected 404 error, got %v", err)
	}
}

func TestGetCharacter_1000000_Fail(t *testing.T) {
	client := rick_morty.Rick_MortyClient{}

	characterData, err := client.GetCharacterByID(context.Background(), 1000000)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if characterData != nil {
		t.Fatalf("expected nil character data, got %v", characterData)
	}
	fmt.Println("Error:", err)
	fmt.Println("Character data:", characterData)

	if err.Error() != "failed to get character: 404 Not Found" {
		t.Fatalf("expected 404 error, got %v", err)
	}
}

// -----------------------------------------------------------
