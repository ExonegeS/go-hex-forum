package model

import (
	"encoding/json"
	"strings"
)

type CharacterData struct {
	ID      int          `json:"id"`
	Name    string       `json:"name"`
	Status  LivingStatus `json:"status"`
	Species string       `json:"species"`
	Type    string       `json:"type"`
	Gender  string       `json:"gender"`
	Image   string       `json:"image"`
}

type LivingStatus int

const (
	Alive LivingStatus = iota
	Dead
	Missing
	Unknown
)

var (
	livingStatusMap = map[LivingStatus]string{
		0: "Alive",
		1: "Dead",
		2: "Missing",
		3: "Unknown",
	}

	livingStatusMapReverse = map[string]LivingStatus{
		"alive":   Alive,
		"dead":    Dead,
		"missing": Missing,
		"unknown": Unknown,
	}
)

func (s LivingStatus) String() string {
	return livingStatusMap[s]
}

func ParseLivingStatus(status string) LivingStatus {
	if s, ok := livingStatusMapReverse[strings.ToLower(status)]; ok {
		return s
	}
	return Unknown
}

func (s *LivingStatus) UnmarshalJSON(data []byte) error {
	var statusStr string
	if err := json.Unmarshal(data, &statusStr); err != nil {
		return err
	}
	*s = ParseLivingStatus(statusStr)
	return nil
}

func (s LivingStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
