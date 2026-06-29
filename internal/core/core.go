// Package core holds the shared Fortinet version parsing and comparison used by
// the numeric and nonnumeric packages. It is internal; depend on
// github.com/vulsio/go-fortinet-version/numeric or .../nonnumeric instead.
//
// A Fortinet version is a "."-separated list of components, each an unsigned
// integer or — when letters are permitted (the FortiSASE milestone scheme) — a
// single lowercase letter. The two schemes share one comparison algorithm; a
// numeric component meeting a milestone letter at the same position has no
// defined cross-scheme order and yields ErrIncomparable.
package core

import (
	"cmp"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrIncomparable is returned by Compare when the two versions have no defined
// order: a numeric component meets a milestone letter at the same position
// (e.g. a build "1.2.1" against a milestone "1.2.a").
var ErrIncomparable = errors.New("incomparable fortinet versions")

type kind int

const (
	kindNumber kind = iota
	kindLetter
)

type component struct {
	kind   kind
	number uint64
	letter byte
}

// Version is a parsed Fortinet version.
type Version struct {
	components []component
}

// Parse splits ver on "." and validates each component as unsigned digits (so a
// signed "+0"/"-1" or an overflowing number is rejected) or, when allowLetters
// is true, a single lowercase milestone letter. A multi-character / non-[a-z]
// token or an empty component (a stray/leading/trailing dot) is an error.
func Parse(ver string, allowLetters bool) (Version, error) {
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
		if allowLetters && len(s) == 1 && s[0] >= 'a' && s[0] <= 'z' {
			components = append(components, component{kind: kindLetter, letter: s[0]})
			continue
		}
		return Version{}, fmt.Errorf("unexpected fortinet version component %q in %q", s, ver)
	}
	return Version{components: components}, nil
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
	// A valid version (from Parse) always has at least one component; an empty
	// one is the zero value used without NewVersion, which has no meaningful
	// order — fail loudly rather than report a bogus result.
	if len(a) == 0 || len(b) == 0 {
		return 0, fmt.Errorf("compare uninitialized fortinet version (construct with NewVersion)")
	}
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
// already validated by Parse, so a letter or non-zero number means greater.
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
