// Package nonnumeric parses and compares non-numeric Fortinet versions — the
// FortiSASE milestone-letter scheme, where a component may be a single
// lowercase letter (a < b < c) as well as a number (e.g. 25.2.a, 25.1.a.2).
// Use the numeric package for the purely numeric scheme every other product
// uses.
package nonnumeric

import (
	"fmt"

	core "github.com/vulsio/go-fortinet-version/internal/core"
)

// ErrIncomparable is returned by Compare when the two versions have no defined
// order: a numeric component meets a milestone letter at the same position
// (e.g. a build "1.2.1" against a milestone "1.2.a"). (numeric versions are
// totally ordered, so only this package needs the sentinel.)
var ErrIncomparable = core.ErrIncomparable

// Version is a parsed non-numeric (milestone-letter) Fortinet version. Always
// construct it with NewVersion; the zero value is invalid.
type Version struct {
	v core.Version
}

// NewVersion parses a non-numeric Fortinet version. Each "."-separated component
// must be unsigned digits or a single lowercase milestone letter; a
// multi-character / non-[a-z] token (e.g. "alpha", "a10"), a signed/overflowing
// number, or an empty component (a stray dot) is an error.
func NewVersion(ver string) (Version, error) {
	v, err := core.Parse(ver, true)
	if err != nil {
		return Version{}, fmt.Errorf("parse nonnumeric version %q: %w", ver, err)
	}
	return Version{v: v}, nil
}

// Compare returns -1, 0, or +1 for v < o, v == o, v > o; milestone letters order
// a < b < c, trailing zeros are a no-op (7.2 == 7.2.0), and a bare train
// precedes its builds (25.2 < 25.2.a). It returns ErrIncomparable when a numeric
// and a letter component meet at the same position.
func (v Version) Compare(o Version) (int, error) {
	return v.v.Compare(o.v)
}

// String renders the parsed version as a dot-separated string, reflecting the
// components verbatim. Trailing zeros are not normalized, so 7.2 and 7.2.0
// render differently even though they compare equal.
func (v Version) String() string {
	return v.v.String()
}
