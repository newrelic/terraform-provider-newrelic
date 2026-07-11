//go:build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

// These tests pin the tag-reading helpers used by resourceNewRelicMonitorDowntimeRead.
// They are the helpers whose missing-tag behavior used to panic the plugin
// (issue #3117): getStringEntityTag returned out[0] without checking length,
// and setMonitorDowntimeStartTime/EndTime would emit 1970-01-01 strings when
// their backing tag was absent. The contract verified here is:
//   - findEntityTagValue reports presence/absence correctly
//   - getStringEntityTag returns "" for absent tags (no panic)
//   - the time helpers return "" for absent/unparseable tags (no stale state)
//   - the scalar helpers (mode/timezone/accountID) tolerate missing tags

func tagSet(pairs ...string) []entities.EntityTag {
	if len(pairs)%2 != 0 {
		panic("tagSet requires key/value pairs")
	}
	out := make([]entities.EntityTag, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, entities.EntityTag{Key: pairs[i], Values: []string{pairs[i+1]}})
	}
	return out
}

func TestFindEntityTagValue(t *testing.T) {
	tags := tagSet("type", "DAILY", "timezone", "UTC")

	if v, ok := findEntityTagValue(tags, "type"); !ok || v != "DAILY" {
		t.Fatalf("present tag: got (%q, %v), want (DAILY, true)", v, ok)
	}
	if v, ok := findEntityTagValue(tags, "missing"); ok || v != "" {
		t.Fatalf("absent tag: got (%q, %v), want (\"\", false)", v, ok)
	}
	if v, ok := findEntityTagValue(nil, "type"); ok || v != "" {
		t.Fatalf("nil slice: got (%q, %v), want (\"\", false)", v, ok)
	}

	// A tag with no values should be reported as absent so callers don't get an
	// empty string they have to special-case.
	empty := []entities.EntityTag{{Key: "type", Values: nil}}
	if v, ok := findEntityTagValue(empty, "type"); ok || v != "" {
		t.Fatalf("empty values: got (%q, %v), want (\"\", false)", v, ok)
	}
}

func TestGetStringEntityTag_AbsentDoesNotPanic(t *testing.T) {
	// The pre-fix version returned out[0] unconditionally and panicked with
	// "index out of range" when the named tag wasn't present. This must stay
	// safe — every downstream setter (mode/timezone/accountID) depends on it.
	got := getStringEntityTag(tagSet("type", "DAILY"), "no-such-tag")
	if got != "" {
		t.Fatalf("absent tag: got %q, want \"\"", got)
	}
}

func TestSetMonitorDowntimeTimeHelpers_MissingTagReturnsEmpty(t *testing.T) {
	// Without these guards, a missing startTime/endTime tag parses as int64 0
	// and renders as "1970-01-01T00:00:00", which would silently overwrite
	// state with a bogus timestamp.
	if got := setMonitorDowntimeStartTime(tagSet("type", "DAILY")); got != "" {
		t.Fatalf("missing startTime: got %q, want \"\"", got)
	}
	if got := setMonitorDowntimeEndTime(tagSet("type", "DAILY")); got != "" {
		t.Fatalf("missing endTime: got %q, want \"\"", got)
	}
}

func TestSetMonitorDowntimeTimeHelpers_UnparseableReturnsEmpty(t *testing.T) {
	tags := tagSet("startTime", "not-a-number", "endTime", "also-not", "timezone", "UTC")
	if got := setMonitorDowntimeStartTime(tags); got != "" {
		t.Fatalf("unparseable startTime: got %q, want \"\"", got)
	}
	if got := setMonitorDowntimeEndTime(tags); got != "" {
		t.Fatalf("unparseable endTime: got %q, want \"\"", got)
	}
}

func TestSetMonitorDowntimeTimeHelpers_BadTimezoneFallsBackToUTC(t *testing.T) {
	// A garbage timezone tag must not abort the read: fall back to UTC so the
	// happy path still produces a usable timestamp.
	tags := tagSet("startTime", "1700000000000", "timezone", "Not/AZone")
	got := setMonitorDowntimeStartTime(tags)
	if got == "" {
		t.Fatalf("bad timezone: got empty string, want UTC-rendered timestamp")
	}
}

func TestSetMonitorDowntimeScalarHelpers_MissingTags(t *testing.T) {
	// All scalar helpers route through getStringEntityTag and must inherit its
	// safe-empty behavior.
	empty := []entities.EntityTag{}
	if got := setMonitorDowntimeMode(empty); got != "" {
		t.Fatalf("missing type: got %q, want \"\"", got)
	}
	if got := setMonitorDowntimeTimezone(empty); got != "" {
		t.Fatalf("missing timezone: got %q, want \"\"", got)
	}
	if got := setMonitorDowntimeAccountID(empty); got != 0 {
		t.Fatalf("missing accountId: got %d, want 0", got)
	}
}

func TestSetMonitorDowntimeAccountID_NonIntegerReturnsZero(t *testing.T) {
	got := setMonitorDowntimeAccountID(tagSet("accountId", "not-a-number"))
	if got != 0 {
		t.Fatalf("non-integer accountId: got %d, want 0", got)
	}
}
