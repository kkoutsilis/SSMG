package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

type Participant struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	WishList []string `json:"wishlist,omitempty"`
}

func LoadParticipants(path string) ([]Participant, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read participants file: %w", err)
	}

	var participants []Participant
	if err := json.Unmarshal(file, &participants); err != nil {
		return nil, fmt.Errorf("parse participants file: %w", err)
	}
	if len(participants) == 0 {
		return nil, fmt.Errorf("participants file is empty")
	}
	return participants, nil
}
