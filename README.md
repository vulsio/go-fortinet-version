# go-fortinet-version

Fortinet product version parsing and comparison, split by version scheme into
two packages (each its own type, so the scheme is enforced at parse time):

- [`numeric`](./numeric) — the purely numeric scheme almost every product uses
  (`7.4.3`, train `7.2`). A milestone letter is rejected by `NewVersion`.
- [`nonnumeric`](./nonnumeric) — the FortiSASE milestone-letter scheme
  (`25.2.a`, `25.1.a.2`).

```go
package main

import (
	"fmt"
	"log"

	"github.com/vulsio/go-fortinet-version/nonnumeric"
	"github.com/vulsio/go-fortinet-version/numeric"
)

func main() {
	a, err := numeric.NewVersion("7.4.3")
	if err != nil {
		log.Fatal(err)
	}
	b, err := numeric.NewVersion("7.4.10")
	if err != nil {
		log.Fatal(err)
	}
	n, err := a.Compare(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n) // -1  (7.4.3 < 7.4.10)

	// trailing zeros are a no-op
	c, err := numeric.NewVersion("7.2.0")
	if err != nil {
		log.Fatal(err)
	}
	d, err := numeric.NewVersion("7.2")
	if err != nil {
		log.Fatal(err)
	}
	if n, err = c.Compare(d); err != nil {
		log.Fatal(err)
	}
	fmt.Println(n) // 0

	// FortiSASE milestone letters: a bare train precedes its milestone builds
	s, err := nonnumeric.NewVersion("25.2")
	if err != nil {
		log.Fatal(err)
	}
	m, err := nonnumeric.NewVersion("25.2.a")
	if err != nil {
		log.Fatal(err)
	}
	if n, err = s.Compare(m); err != nil {
		log.Fatal(err)
	}
	fmt.Println(n) // -1
}
```

Each package exposes `NewVersion(string) (Version, error)`,
`(Version).Compare(Version) (int, error)`, and `(Version).String()`. The
`nonnumeric` package additionally exposes the `ErrIncomparable` sentinel,
returned by `(nonnumeric.Version).Compare` when a numeric component meets a
milestone letter at the same position (e.g. `1.2.1` vs `1.2.a`); test for it with
`errors.Is`. `numeric` versions are totally ordered, so
`(numeric.Version).Compare` never returns it. The shared core lives in
`internal/core`.
