package newrelic

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/maintenancewindow"
	"github.com/newrelic/newrelic-client-go/v2/pkg/workloads"
)

func expandMaintenanceWindowCreateInput(d *schema.ResourceData) maintenancewindow.MaintenanceWindowConfigCreateInput {
	return maintenancewindow.MaintenanceWindowConfigCreateInput{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		StartTime:          maintenancewindow.NaiveDateTime(d.Get("start_time").(string)),
		Duration:           workloads.Duration(d.Get("duration").(string)),
		Rrule:              d.Get("rrule").(string),
		Timezone:           d.Get("timezone").(string),
		AffectedEntityType: maintenancewindow.MaintenanceWindowAffectedEntityType(d.Get("affected_entity_type").(string)),
		AffectedEntities:   expandMaintenanceWindowAffectedEntities(d.Get("affected_entities").([]interface{})),
		Scope:              expandMaintenanceWindowScope(d.Get("scope").([]interface{})),
	}
}

func expandMaintenanceWindowScope(cfg []interface{}) maintenancewindow.MaintenanceWindowScopedReferenceInput {
	scope := cfg[0].(map[string]interface{})
	return maintenancewindow.MaintenanceWindowScopedReferenceInput{
		ID:   scope["id"].(string),
		Type: maintenancewindow.MaintenanceWindowEntityScope(scope["type"].(string)),
	}
}

func expandMaintenanceWindowUpdateInput(d *schema.ResourceData) maintenancewindow.MaintenanceWindowConfigUpdateInput {
	return maintenancewindow.MaintenanceWindowConfigUpdateInput{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		StartTime:        maintenancewindow.NaiveDateTime(d.Get("start_time").(string)),
		Duration:         workloads.Duration(d.Get("duration").(string)),
		Rrule:            d.Get("rrule").(string),
		Timezone:         d.Get("timezone").(string),
		AffectedEntities: expandMaintenanceWindowAffectedEntities(d.Get("affected_entities").([]interface{})),
	}
}

func expandMaintenanceWindowAffectedEntities(cfg []interface{}) []string {
	entities := make([]string, len(cfg))
	for i, e := range cfg {
		entities[i] = e.(string)
	}
	return entities
}

func flattenMaintenanceWindow(mw *maintenancewindow.MaintenanceWindowConfigDefinition, d *schema.ResourceData) error {
	if err := d.Set("name", mw.Name); err != nil {
		return err
	}
	if err := d.Set("description", mw.Description); err != nil {
		return err
	}
	if err := d.Set("scope", flattenMaintenanceWindowScope(mw.Scope)); err != nil {
		return err
	}
	if err := d.Set("start_time", string(mw.StartTime)); err != nil {
		return err
	}
	if err := d.Set("duration", string(mw.Duration)); err != nil {
		return err
	}
	if err := d.Set("rrule", mw.Rrule); err != nil {
		return err
	}
	if err := d.Set("timezone", mw.Timezone); err != nil {
		return err
	}
	if err := d.Set("affected_entity_type", string(mw.AffectedEntityType)); err != nil {
		return err
	}
	if err := d.Set("affected_entities", mw.AffectedEntities); err != nil {
		return err
	}

	return nil
}

func flattenMaintenanceWindowScope(scope maintenancewindow.MaintenanceWindowScopedReference) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":   scope.ID,
			"type": string(scope.Type),
		},
	}
}

// validate functions which perform custom validation, mirroring the server-side
// rules in maintenance-window-service's ValidationService/RRuleValidator so that
// invalid input is caught at `terraform plan` time instead of a round-trip API error.

func validateMaintenanceWindowTimeZone(val interface{}, key string) (warns []string, errs []error) {
	timezone := val.(string)
	if _, err := time.LoadLocation(timezone); err != nil {
		errs = append(errs, err)
	}

	return warns, errs
}

func validateMaintenanceWindowStartTime(val interface{}, key string) (warns []string, errs []error) {
	startTime := val.(string)
	if _, err := time.Parse("2006-01-02T15:04:05", startTime); err != nil {
		errs = append(errs, fmt.Errorf("%q must be an ISO 8601 date-time without a time zone offset (e.g. 2026-12-01T22:00:00): %w", key, err))
	}

	return warns, errs
}

func validateMaintenanceWindowDuration(val interface{}, key string) (warns []string, errs []error) {
	if _, err := parseISO8601Duration(val.(string)); err != nil {
		errs = append(errs, fmt.Errorf("%q: %w", key, err))
	}

	return warns, errs
}

var (
	maintenanceWindowRruleCharsetRegex      = regexp.MustCompile(`^[A-Za-z0-9=;,:+TZ-]*$`)
	maintenanceWindowWeeklyByDayRegex       = regexp.MustCompile(`^(MO|TU|WE|TH|FR|SA|SU)(,(MO|TU|WE|TH|FR|SA|SU))*$`)
	maintenanceWindowMonthlyByDayRegex      = regexp.MustCompile(`^-?[1-4](MO|TU|WE|TH|FR|SA|SU)$`)
	maintenanceWindowMonthlyByMonthDayRegex = regexp.MustCompile(`^-?([1-9]|[12][0-9]|3[01])(,-?([1-9]|[12][0-9]|3[01]))*$`)
)

// validateMaintenanceWindowRrule mirrors RRuleValidator.java: max 200 chars, a
// restricted charset, no COUNT, INTERVAL must be 1 if present, FREQ limited to
// DAILY/WEEKLY/MONTHLY, and frequency-specific BYDAY/BYMONTHDAY requirements.
func validateMaintenanceWindowRrule(val interface{}, key string) (warns []string, errs []error) {
	rrule := val.(string)
	if rrule == "" {
		return warns, errs
	}

	if len(rrule) > 200 {
		errs = append(errs, fmt.Errorf("%q must be 200 characters or fewer", key))
		return warns, errs
	}

	if !maintenanceWindowRruleCharsetRegex.MatchString(rrule) {
		errs = append(errs, fmt.Errorf("%q contains characters not permitted in an RRULE", key))
		return warns, errs
	}

	rules := map[string]string{}
	for _, part := range strings.Split(rrule, ";") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			errs = append(errs, fmt.Errorf("%q contains a malformed component %q — expected KEY=VALUE", key, part))
			return warns, errs
		}
		rules[strings.ToUpper(kv[0])] = kv[1]
	}

	if _, ok := rules["COUNT"]; ok {
		errs = append(errs, fmt.Errorf("%q must not include COUNT", key))
	}

	if interval, ok := rules["INTERVAL"]; ok && interval != "1" {
		errs = append(errs, fmt.Errorf("%q: INTERVAL must be 1 if present", key))
	}

	freq, ok := rules["FREQ"]
	if !ok {
		errs = append(errs, fmt.Errorf("%q must include FREQ", key))
		return warns, errs
	}

	switch freq {
	case "DAILY":
		for k := range rules {
			if strings.HasPrefix(k, "BY") {
				errs = append(errs, fmt.Errorf("%q: FREQ=DAILY does not permit any BY* rules (found %s)", key, k))
			}
		}
	case "WEEKLY":
		byDay, hasByDay := rules["BYDAY"]
		if !hasByDay {
			errs = append(errs, fmt.Errorf("%q: FREQ=WEEKLY requires BYDAY", key))
		} else if !maintenanceWindowWeeklyByDayRegex.MatchString(byDay) {
			errs = append(errs, fmt.Errorf("%q: BYDAY for FREQ=WEEKLY must be a comma-separated list of day codes with no numeric prefix (e.g. MO,WE,FR)", key))
		}
	case "MONTHLY":
		byDay, hasByDay := rules["BYDAY"]
		byMonthDay, hasByMonthDay := rules["BYMONTHDAY"]
		switch {
		case hasByDay && hasByMonthDay:
			errs = append(errs, fmt.Errorf("%q: FREQ=MONTHLY must specify exactly one of BYDAY or BYMONTHDAY, not both", key))
		case !hasByDay && !hasByMonthDay:
			errs = append(errs, fmt.Errorf("%q: FREQ=MONTHLY must specify exactly one of BYDAY or BYMONTHDAY", key))
		case hasByDay:
			if !maintenanceWindowMonthlyByDayRegex.MatchString(byDay) {
				errs = append(errs, fmt.Errorf("%q: BYDAY for FREQ=MONTHLY must include a numeric prefix of -1 or 1-4 (e.g. 2TU, -1FR)", key))
			}
		case hasByMonthDay:
			if !maintenanceWindowMonthlyByMonthDayRegex.MatchString(byMonthDay) {
				errs = append(errs, fmt.Errorf("%q: BYMONTHDAY for FREQ=MONTHLY must be -1 or 1-31", key))
			}
		}
	default:
		errs = append(errs, fmt.Errorf("%q: FREQ must be one of DAILY, WEEKLY, or MONTHLY (got %s)", key, freq))
	}

	return warns, errs
}

// validateMaintenanceWindowAttributes is wired up as the resource's CustomizeDiff,
// for validation that depends on more than one field at once.
func validateMaintenanceWindowAttributes(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	rrule := d.Get("rrule").(string)
	if rrule == "" {
		return nil
	}

	duration, err := parseISO8601Duration(d.Get("duration").(string))
	if err != nil {
		// Already caught by the field-level ValidateFunc on "duration"; nothing more to add here.
		return nil
	}

	if duration > 24*time.Hour {
		return fmt.Errorf("duration must not exceed 24 hours when rrule is set (got %s)", d.Get("duration").(string))
	}

	return nil
}

// parseISO8601Duration parses a subset of ISO 8601 durations (e.g. PT1H, P1D,
// P1DT2H30M) into a time.Duration. Y and M (calendar month) components are
// approximated as 365 and 30 days respectively, which is sufficient for the
// one thing this parser is used for: validating format and the 24-hour cap.
func parseISO8601Duration(s string) (time.Duration, error) {
	if s == "" || s[0] != 'P' {
		return 0, fmt.Errorf("duration must start with \"P\" as defined by ISO 8601 (e.g. PT1H, P1D)")
	}

	datePart := s[1:]
	timePart := ""
	if idx := strings.Index(datePart, "T"); idx != -1 {
		timePart = datePart[idx+1:]
		datePart = datePart[:idx]
	}

	if datePart == "" && timePart == "" {
		return 0, fmt.Errorf("duration must specify at least one component (e.g. PT1H, P1D)")
	}

	var total time.Duration
	matched := false

	dateUnits := []struct {
		suffix string
		unit   time.Duration
	}{
		{"Y", 365 * 24 * time.Hour},
		{"M", 30 * 24 * time.Hour},
		{"W", 7 * 24 * time.Hour},
		{"D", 24 * time.Hour},
	}

	remaining := datePart
	for _, u := range dateUnits {
		idx := strings.Index(remaining, u.suffix)
		if idx == -1 {
			continue
		}
		n, err := strconv.Atoi(remaining[:idx])
		if err != nil {
			return 0, fmt.Errorf("invalid ISO 8601 duration: %q", s)
		}
		total += time.Duration(n) * u.unit
		remaining = remaining[idx+1:]
		matched = true
	}
	if remaining != "" {
		return 0, fmt.Errorf("invalid ISO 8601 duration: %q", s)
	}

	timeUnits := []struct {
		suffix string
		unit   time.Duration
	}{
		{"H", time.Hour},
		{"M", time.Minute},
		{"S", time.Second},
	}

	remaining = timePart
	for _, u := range timeUnits {
		idx := strings.Index(remaining, u.suffix)
		if idx == -1 {
			continue
		}
		n, err := strconv.Atoi(remaining[:idx])
		if err != nil {
			return 0, fmt.Errorf("invalid ISO 8601 duration: %q", s)
		}
		total += time.Duration(n) * u.unit
		remaining = remaining[idx+1:]
		matched = true
	}
	if remaining != "" {
		return 0, fmt.Errorf("invalid ISO 8601 duration: %q", s)
	}

	if !matched {
		return 0, fmt.Errorf("duration must specify at least one component (e.g. PT1H, P1D)")
	}

	return total, nil
}
