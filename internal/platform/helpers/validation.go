package helpers

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/lackmus/settlementgengo/pkg/model"
)

func ValidateSettlementName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("settlement name cannot be empty")
	}
	if len(name) > 50 {
		return errors.New("settlement name cannot exceed 50 characters")
	}
	return nil
}

func ValidateFactionName(faction string) error {
	if strings.TrimSpace(faction) == "" {
		return errors.New("faction name cannot be empty")
	}
	if len(faction) > 50 {
		return errors.New("faction name cannot exceed 50 characters")
	}
	return nil
}

func ValidateCoordinates(x, y int) error {
	if x < 0 || x > 1000 {
		return fmt.Errorf("x coordinate must be between 0 and 1000, got %d", x)
	}
	if y < 0 || y > 1000 {
		return fmt.Errorf("y coordinate must be between 0 and 1000, got %d", y)
	}
	return nil
}

func ValidatePopulation(population int) error {
	if population < 0 {
		return fmt.Errorf("population cannot be negative, got %d", population)
	}
	if population > 1000000 {
		return fmt.Errorf("population cannot exceed 1,000,000, got %d", population)
	}
	return nil
}

// ValidateNPCName checks if the NPC name is valid (not empty and not too long) and denys code injection attempts
func ValidateNotes(notes string) error {
	if len(notes) > 500 {
		return errors.New("notes cannot exceed 500 characters")
	}
	if containsCodeInjection(notes) {
		return errors.New("notes cannot contain code injection attempts")
	}
	return nil
}

func containsCodeInjection(input string) bool {
	normalized := strings.ToLower(input)
	patterns := []string{"<script>", "</script>", "javascript:", "onerror=", "onload="}
	for _, pattern := range patterns {
		if strings.Contains(normalized, pattern) {
			return true
		}
	}
	return false
}

func ValidateSettlement(settlement model.Settlement) error {
	if err := ValidateSettlementName(settlement.Name); err != nil {
		return fmt.Errorf("invalid settlement name: %w", err)
	}
	if err := ValidateFactionName(settlement.Faction); err != nil {
		return fmt.Errorf("invalid faction: %w", err)
	}
	if err := ValidateCoordinates(settlement.XCoord, settlement.YCoord); err != nil {
		return fmt.Errorf("invalid coordinates: %w", err)
	}
	if err := ValidatePopulation(settlement.Population); err != nil {
		return fmt.Errorf("invalid population: %w", err)
	}
	if err := ValidateNotes(settlement.Notes); err != nil {
		return fmt.Errorf("invalid notes: %w", err)
	}
	return nil
}

func IsNilOrEmpty[T any](t T) bool {
	if reflect.ValueOf(t).IsZero() {
		return true
	}

	if str, ok := any(t).(string); ok && str == "" {
		return true
	}

	return false
}
