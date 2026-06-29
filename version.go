// Package version parses and compares Fortinet product version strings.
//
// A Fortinet version is a "."-separated list of components, each either an
// unsigned integer or a single lowercase milestone letter. Almost every product
// is purely numeric (e.g. 7.4.3, or a release train 7.2); FortiSASE labels
// releases with a milestone-letter scheme that carries an alphabetic component
// (e.g. 25.2.a, 25.1.a.2), which plain semver cannot represent.
//
// The two schemes share one comparison algorithm: a numeric component and a
// milestone letter meeting at the same position have no defined cross-scheme
// order, so Compare reports ErrIncomparable rather than guessing. Callers that
// model a purely numeric product can use IsNumeric to reject a lettered version
// before comparing.
package version

import (
	"cmp"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrIncomparable is returned by Compare when the two versions have no defined
// order: a numeric component meets a milestone letter at the same position
// (e.g. a build "1.2.1" against a milestone "1.2.a"), which spans Fortinet's
// numeric and milestone-letter schemes.
var ErrIncomparable = errors.New("incomparable fortinet versions")

type kind int

const (
	kindNumber kind = iota
	kindLetter
)

// component is one "."-separated piece of a version: an unsigned number or a
// single lowercase milestone letter.
type component struct {
	kind   kind
	number uint64
	letter byte
}

// Version is a parsed Fortinet version.
type Version struct {
	components []component
}

// NewVersion parses a Fortinet version string. Each "."-separated component must
// be unsigned digits (so a signed "+0"/"-1" or an overflowing number is
// rejected) or a single lowercase letter; a multi-character / non-[a-z] token
// (e.g. "alpha", "a10", "x"-style placeholders longer than one char) or an
// empty component (a stray/leading/trailing dot) is an error.
func NewVersion(ver string) (Version, error) {
	if ver == "" {
		return Version{}, fmt.Errorf("empty fortinet version")
	}
	ss := strings.Split(ver, ".")
	components := make([]component, 0, len(ss))
	for _, s := range ss {
		if n, err := strconv.ParseUint(s, 10, 64); err == nil {
			components = append(components, component{kind: kindNumber, number: n})
			continue
		}
		if len(s) == 1 && s[0] >= 'a' && s[0] <= 'z' {
			components = append(components, component{kind: kindLetter, letter: s[0]})
			continue
		}
		return Version{}, fmt.Errorf("unexpected fortinet version component %q in %q", s, ver)
	}
	return Version{components: components}, nil
}

// IsNumeric reports whether every component is numeric (the version carries no
// milestone letter). A purely numeric product can use this to refuse a lettered
// version (e.g. treating it as a non-match) before calling Compare.
func (v Version) IsNumeric() bool {
	for _, c := range v.components {
		if c.kind != kindNumber {
			return false
		}
	}
	return true
}

// Compare returns -1, 0, or +1 for v < o, v == o, v > o. Numeric components
// compare numerically and milestone letters lexically (a < b < c). When one
// version is a prefix of the other the remaining tail decides: trailing numeric
// zeros are a no-op (7.2 == 7.2.0), while a non-zero or lettered tail makes the
// longer version greater, so a bare train precedes its builds (25.2 < 25.2.a).
// It returns ErrIncomparable when a numeric and a letter component meet at the
// same position.
func (v Version) Compare(o Version) (int, error) {
	a, b := v.components, o.components
	for i := 0; i < len(a) || i < len(b); i++ {
		switch {
		case i >= len(a):
			return -tailSign(b[i:]), nil
		case i >= len(b):
			return tailSign(a[i:]), nil
		}
		ca, cb := a[i], b[i]
		switch {
		case ca.kind == kindNumber && cb.kind == kindNumber:
			if ca.number != cb.number {
				return cmp.Compare(ca.number, cb.number), nil
			}
		case ca.kind == kindLetter && cb.kind == kindLetter:
			if ca.letter != cb.letter {
				return cmp.Compare(ca.letter, cb.letter), nil
			}
		default:
			return 0, ErrIncomparable
		}
	}
	return 0, nil
}

// tailSign reports whether the trailing components of the longer version make it
// greater (1) or leave the two equal (0, a trailing-zero no-op). Components are
// already validated by NewVersion, so a letter or non-zero number means greater.
func tailSign(comps []component) int {
	for _, c := range comps {
		if c.kind != kindNumber || c.number != 0 {
			return 1
		}
	}
	return 0
}

// String reformats the parsed version (canonical, dot-separated).
func (v Version) String() string {
	ss := make([]string, len(v.components))
	for i, c := range v.components {
		if c.kind == kindLetter {
			ss[i] = string(c.letter)
		} else {
			ss[i] = strconv.FormatUint(c.number, 10)
		}
	}
	return strings.Join(ss, ".")
}
