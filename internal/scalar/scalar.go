// Package scalar holds the scalar rules the JSON, XML and YAML packages share.
//
// It exists because the three formats used to answer "is this text a number?"
// three different ways, and the answers disagreed. strconv.ParseFloat accepts
// "Inf", "NaN" and hexadecimal floats, none of which JSON can write, so an XML
// element holding the word "Inf" became a number on the way in and an invalid
// JSON document on the way out. One rule, used everywhere, is the fix.
//
// The package is internal on purpose: the Node model stays a plain string-based
// data model, so nothing here is part of the module's public API.
package scalar

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
)

// IsJSONNumber reports whether s is a number literal JSON accepts verbatim.
//
// The leading-character test is what separates a number from the other bare
// JSON literals: json.Valid also accepts "true" and "null".
func IsJSONNumber(s string) bool {
	if s == "" {
		return false
	}
	if c := s[0]; c != '-' && (c < '0' || c > '9') {
		return false
	}
	return json.Valid([]byte(s))
}

// Number normalizes a numeric literal into text every format can write, and
// reports whether it is a number at all.
//
// A literal JSON already accepts is returned untouched, so 9007199254740993
// and 1e10 keep the digits they were written with instead of being routed
// through float64. Forms only YAML accepts (0x1F, 0o17, 1_000, +5) are
// converted to plain decimal. Infinities and NaN are not numbers here: no
// format in this module can write them, so they travel as text.
func Number(s string) (string, bool) {
	if IsJSONNumber(s) {
		return s, true
	}
	if i, err := strconv.ParseInt(s, 0, 64); err == nil {
		return strconv.FormatInt(i, 10), true
	}
	if u, err := strconv.ParseUint(s, 0, 64); err == nil {
		return strconv.FormatUint(u, 10), true
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil && !math.IsInf(f, 0) && !math.IsNaN(f) {
		return strconv.FormatFloat(f, 'f', -1, 64), true
	}
	return s, false
}

// IsInteger reports whether a literal already known to be a number is integral.
// The caller uses it to pick between YAML's !!int and !!float tags.
func IsInteger(s string) bool {
	if _, err := strconv.ParseInt(s, 10, 64); err == nil {
		return true
	}
	_, err := strconv.ParseUint(s, 10, 64)
	return err == nil
}

// Bool parses the boolean spellings the three formats use between them, and
// reports whether the text was one of them. YAML 1.1's yes/no/on/off are
// included because documents in the wild still carry them.
func Bool(s string) (bool, bool) {
	if b, err := strconv.ParseBool(s); err == nil {
		return b, true
	}
	switch strings.ToLower(s) {
	case "yes", "on":
		return true, true
	case "no", "off":
		return false, true
	}
	return false, false
}
