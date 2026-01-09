package foodora

import (
	"encoding/json"
	"fmt"
	"time"
)

// FlexibleTime decodes common API time formats:
// - RFC3339 strings
// - unix timestamps (seconds or milliseconds)
// - null
type FlexibleTime struct {
	time.Time
}

func (t FlexibleTime) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func (t *FlexibleTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		t.Time = time.Time{}
		return nil
	}

	if len(b) > 0 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if s == "" {
			t.Time = time.Time{}
			return nil
		}
		if parsed, err := parseAPITimeString(s); err == nil {
			t.Time = parsed
			return nil
		} else {
			return fmt.Errorf("parse time %q: %w", s, err)
		}
	}

	var n float64
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	if n == 0 {
		t.Time = time.Time{}
		return nil
	}
	iv := int64(n)
	// Heuristic: ms timestamps are already > 1e12 in 2001+.
	if iv > 1_000_000_000_000 {
		t.Time = time.UnixMilli(iv).UTC()
		return nil
	}
	t.Time = time.Unix(iv, 0).UTC()
	return nil
}

func parseAPITimeString(s string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	var lastErr error
	for _, layout := range layouts {
		tt, err := time.Parse(layout, s)
		if err == nil {
			return tt, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

const ReorderTimeLayout = "2006-01-02T15:04:05-0700"

func FormatReorderTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(ReorderTimeLayout)
}
