package generate

import (
	"testing"

	"github.com/google/uuid"
)

func TestFixedUUID(t *testing.T) {
	// Set a fixed UUID
	SetFixedUUID("123e4567-e89b-12d3-a456-426614174000")
	defer ClearFixedUUID()

	// Call the UUID function
	uuid := UUID()

	// Check if the returned UUID matches the fixed UUID
	if uuid != "123e4567-e89b-12d3-a456-426614174000" {
		t.Errorf("Expected UUID to be '123e4567-e89b-12d3-a456-426614174000', got '%s'", uuid)
	}
}

func TestRandomUUID(t *testing.T) {
	// Clear any fixed UUID
	ClearFixedUUID()

	// Call the UUID function
	id := UUID()

	// Check if the returned UUID is not empty
	if id == "" {
		t.Error("Expected a non-empty UUID")
	}

	// Check if the returned UUID is a valid UUID format
	if _, err := uuid.Parse(id); err != nil {
		t.Errorf("Expected a valid UUID, got '%s'", id)
	}
}
