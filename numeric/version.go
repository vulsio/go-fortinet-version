// Package numeric parses and compares purely numeric Fortinet versions
// (e.g. 7.4.3, or a release train 7.2) — the scheme almost every Fortinet
// product uses. A milestone letter is rejected at parse time; use the
// nonnumeric package for FortiSASE-style versions (25.2.a).
package numeric

import (
	"fmt"

	core "github.com/vulsio/go-fortinet-version/internal/core"
)

// Version is a parsed numeric Fortinet version. Always construct it with
// NewVersion; the zero value is invalid.
type Version struct {
	v core.Version
}

// NewVersion parses a numeric Fortinet version. Each "."-separated component
// must be unsigned digits; a milestone letter, a signed/overflowing number, or
// an empty component (a stray dot) is an error.
func NewVersion(ver string) (Version, error) {
	v, err := core.Parse(ver, false)
	if err != nil {
		return Version{}, fmt.Errorf("parse numeric version %q: %w", ver, err)
	}
	return Version{v: v}, nil
}

// Compare returns -1, 0, or +1 for v < o, v == o, v > o, with trailing zeros a
// no-op (7.2 == 7.2.0). Numeric versions are totally ordered; the only error it
// returns is for a zero-value operand (built without NewVersion).
func (v Version) Compare(o Version) (int, error) {
	return v.v.Compare(o.v)
}

// String reformats the parsed version (canonical, dot-separated).
func (v Version) String() string {
	return v.v.String()
}
