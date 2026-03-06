package helpers

import (
	"strings"
	"testing"

	"github.com/lackmus/settlementgengo/pkg/model"
)

func TestValidateSettlementName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid", input: "Hollow Creek", wantErr: false},
		{name: "empty", input: "", wantErr: true},
		{name: "whitespace", input: "   ", wantErr: true},
		{name: "too long", input: strings.Repeat("a", 51), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSettlementName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateSettlementName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateNotes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid", input: "Quiet village by the river.", wantErr: false},
		{name: "script tag", input: "<script>alert(1)</script>", wantErr: true},
		{name: "javascript scheme", input: "Click javascript:alert(1)", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNotes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateNotes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSettlement(t *testing.T) {
	valid := model.Settlement{
		Name:       "Amber Ridge",
		Faction:    "Woodland Alliance",
		XCoord:     100,
		YCoord:     200,
		Population: 350,
		Notes:      "Trade-focused settlement.",
	}

	if err := ValidateSettlement(valid); err != nil {
		t.Fatalf("ValidateSettlement(valid) unexpected error: %v", err)
	}

	invalidCoords := valid
	invalidCoords.XCoord = -1
	if err := ValidateSettlement(invalidCoords); err == nil {
		t.Fatal("ValidateSettlement(invalidCoords) expected error, got nil")
	}
}
