package cli

import (
	"strings"
	"testing"
)

func TestPickCustomerAddress_NoAddresses(t *testing.T) {
	_, err := pickCustomerAddress(nil, "")
	if err == nil || !strings.Contains(err.Error(), "no customer addresses") {
		t.Fatalf("expected no addresses error, got: %v", err)
	}
}

func TestPickCustomerAddress_Single(t *testing.T) {
	addr := map[string]any{"id": "1"}
	got, err := pickCustomerAddress([]map[string]any{addr}, "")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got["id"] != "1" {
		t.Fatalf("unexpected address: %#v", got)
	}
}

func TestPickCustomerAddress_ByID(t *testing.T) {
	a1 := map[string]any{"id": "1"}
	a2 := map[string]any{"id": "2"}
	got, err := pickCustomerAddress([]map[string]any{a1, a2}, "2")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got["id"] != "2" {
		t.Fatalf("unexpected address: %#v", got)
	}
}

func TestPickCustomerAddress_ByID_NotFound(t *testing.T) {
	a1 := map[string]any{"id": "1"}
	a2 := map[string]any{"id": "2"}
	_, err := pickCustomerAddress([]map[string]any{a1, a2}, "3")
	if err == nil || !strings.Contains(err.Error(), "not found") || !strings.Contains(err.Error(), "available: 1,2") {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestPickCustomerAddress_Selected(t *testing.T) {
	a1 := map[string]any{"id": "1"}
	a2 := map[string]any{"id": "2", "is_selected": true}
	got, err := pickCustomerAddress([]map[string]any{a1, a2}, "")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got["id"] != "2" {
		t.Fatalf("unexpected address: %#v", got)
	}
}

func TestPickCustomerAddress_Default(t *testing.T) {
	a1 := map[string]any{"id": "1", "is_default": true}
	a2 := map[string]any{"id": "2"}
	got, err := pickCustomerAddress([]map[string]any{a2, a1}, "")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got["id"] != "1" {
		t.Fatalf("unexpected address: %#v", got)
	}
}

func TestPickCustomerAddress_Multiple_NoHeuristic(t *testing.T) {
	a1 := map[string]any{"id": "1"}
	a2 := map[string]any{"id": "2"}
	_, err := pickCustomerAddress([]map[string]any{a1, a2}, "")
	if err == nil || !strings.Contains(err.Error(), "--address-id") {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestIsTruthy(t *testing.T) {
	cases := []struct {
		v    any
		want bool
	}{
		{nil, false},
		{false, false},
		{true, true},
		{"true", true},
		{"TRUE", true},
		{" 1 ", true},
		{"yes", true},
		{"no", false},
		{"0", false},
		{float64(0), false},
		{float64(2), true},
		{0, false},
		{3, true},
	}
	for i, c := range cases {
		if got := isTruthy(c.v); got != c.want {
			t.Fatalf("case %d: isTruthy(%v)=%v want %v", i, c.v, got, c.want)
		}
	}
}
