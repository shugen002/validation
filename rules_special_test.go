package validation

import (
	"testing"

	"github.com/google/uuid"
)

func TestUuidRule(t *testing.T) {
	rule := &UuidRule{}
	validUUID := uuid.New().String()
	if !rule.Passes("id", validUUID) {
		t.Errorf("Expected valid UUID to pass")
	}
	if rule.Passes("id", "not-a-uuid") {
		t.Errorf("Expected invalid UUID to fail")
	}
}

func TestUlidRule(t *testing.T) {
	rule := &UlidRule{}
	validULID := "01ARZ3NDEKTSV4RRFFQ69G5FAV" // 26 chars, valid Crockford base32
	if !rule.Passes("ulid", validULID) {
		t.Errorf("Expected valid ULID to pass")
	}
	if rule.Passes("ulid", "not-a-ulid") {
		t.Errorf("Expected invalid ULID to fail")
	}
}

func TestNullableRule(t *testing.T) {
	rule := &NullableRule{}
	if !rule.Passes("any", nil) {
		t.Errorf("NullableRule should always pass (nil)")
	}
	if !rule.Passes("any", "something") {
		t.Errorf("NullableRule should always pass (non-nil)")
	}
}
