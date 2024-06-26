package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// basicFunctions are the set of initial functions provided to every template.
//
//nolint:unused // templating functions.
var basicFunctions = template.FuncMap{
	"json": func(v interface{}) string {
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(v)
		// Remove the trailing new line added by the encoder
		return strings.TrimSpace(buf.String())
	},
	"split":    strings.Split,
	"join":     strings.Join,
	"title":    cases.Title(language.English).String,
	"lower":    cases.Lower(language.English).String,
	"upper":    cases.Upper(language.English).String,
	"pad":      padWithSpace,
	"padlen":   padToLength,
	"padmax":   padToMaxLength,
	"truncate": truncateWithLength,
	"tf":       stringTrueFalse,
	"yn":       stringYesNo,
	"t":        stringTab,
	"time":     timeFormat,
	"date":     dateFormat,
}

// padToLength adds whitespace to pad to the supplied length.
//
//nolint:unused // is used by basicFunctions.
func padToMaxLength(source string) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", 0), source)
}

// padToLength adds whitespace to pad to the supplied length.
//
//nolint:unused // is used by basicFunctions.
func padToLength(source string, prefix int) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", prefix), source)
}

// padWithSpace adds whitespace to the input if the input is non-empty.
//
//nolint:unused // is used by basicFunctions.
func padWithSpace(source string, prefix, suffix int) string {
	if source == "" {
		return source
	}

	return strings.Repeat(" ", prefix) + source + strings.Repeat(" ", suffix)
}

// timeFormat returns time in RFC3339 format.
//
//nolint:unused // is used by basicFunctions.
func timeFormat(source interface{}) string {
	switch s := source.(type) {
	case time.Time:
		return s.Format(time.RFC3339)
	case timestamppb.Timestamp:
		return s.AsTime().Format(time.RFC3339)
	case *timestamppb.Timestamp:
		return s.AsTime().Format(time.RFC3339)
	default:
		return fmt.Sprintf("%s", source)
	}
}

// dateFormat returns date in YYYY-MM-DD format.
//
//nolint:unused // is used by basicFunctions.
func dateFormat(source interface{}) string {
	switch s := source.(type) {
	case time.Time:
		return s.Format("2006-01-02")
	case timestamppb.Timestamp:
		return s.AsTime().Format("2006-01-02")
	case *timestamppb.Timestamp:
		return s.AsTime().Format("2006-01-02")
	default:
		return fmt.Sprintf("%q", source)
	}
}

// stringTrueFalse returns "true" or "false" for boolean input.
//
//nolint:unused // is used by basicFunctions.
func stringTrueFalse(source bool) string {
	if source {
		return "true"
	}

	return "false"
}

// stringYesNo returns "yes" or "no" for boolean input.
//
//nolint:unused // is used by basicFunctions.
func stringYesNo(source bool) string {
	if source {
		return "yes"
	}

	return "no"
}

// stringTab returns a tab character.
//
//nolint:unused // is used by basicFunctions.
func stringTab() string {
	return "\t"
}

// truncateWithLength truncates the source string up to the length provided by the input.
//
//nolint:unused // is used by basicFunctions.
func truncateWithLength(source string, length int) string {
	if len(source) < length {
		return source
	}

	return source[:length]
}
