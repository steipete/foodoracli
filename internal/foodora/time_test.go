package foodora

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFlexibleTime_StringZero(t *testing.T) {
	t.Parallel()
	var ft FlexibleTime
	if ft.String() != "" {
		t.Fatalf("expected empty")
	}
}

func TestFlexibleTime_Unmarshal_String(t *testing.T) {
	t.Parallel()
	var v struct {
		T FlexibleTime `json:"t"`
	}
	if err := json.Unmarshal([]byte(`{"t":"2025-12-20T00:00:00Z"}`), &v); err != nil {
		t.Fatalf("err: %v", err)
	}
	if v.T.IsZero() {
		t.Fatalf("expected time")
	}
	if got := v.T.UTC().Format(time.RFC3339); got != "2025-12-20T00:00:00Z" {
		t.Fatalf("got %q", got)
	}

	if err := json.Unmarshal([]byte(`{"t":"2025-12-20 01:02:03"}`), &v); err != nil {
		t.Fatalf("err: %v", err)
	}
	if v.T.IsZero() {
		t.Fatalf("expected time")
	}
}

func TestFlexibleTime_Unmarshal_NumberSecondsAndMillis(t *testing.T) {
	t.Parallel()
	var v struct {
		T FlexibleTime `json:"t"`
	}

	if err := json.Unmarshal([]byte(`{"t":1734652800}`), &v); err != nil {
		t.Fatalf("err: %v", err)
	}
	if v.T.Unix() != 1734652800 {
		t.Fatalf("got %d", v.T.Unix())
	}

	if err := json.Unmarshal([]byte(`{"t":1734652800123}`), &v); err != nil {
		t.Fatalf("err: %v", err)
	}
	if v.T.UnixMilli() != 1734652800123 {
		t.Fatalf("got %d", v.T.UnixMilli())
	}
}

func TestParseAPITimeString_Error(t *testing.T) {
	t.Parallel()
	if _, err := parseAPITimeString("not-a-time"); err == nil {
		t.Fatalf("expected error")
	}
}
