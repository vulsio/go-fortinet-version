# go-fortinet-version

Fortinet product version parsing and comparison, split by version scheme into
two packages (each its own type, so the scheme is enforced at parse time):

- [`numeric`](./numeric) — the purely numeric scheme almost every product uses
  (`7.4.3`, train `7.2`). A milestone letter is rejected by `NewVersion`.
- [`nonnumeric`](./nonnumeric) — the FortiSASE milestone-letter scheme
  (`25.2.a`, `25.1.a.2`).

```go
import (
	"github.com/vulsio/go-fortinet-version/numeric"
	"github.com/vulsio/go-fortinet-version/nonnumeric"
)

a, _ := numeric.NewVersion("7.2.0")
b, _ := numeric.NewVersion("7.2")
n, _ := a.Compare(b) // 0  (trailing zero is a no-op)

s, _ := nonnumeric.NewVersion("25.2")
m, _ := nonnumeric.NewVersion("25.2.a")
k, _ := s.Compare(m) // -1 (a bare train precedes its milestone builds)
```

Each package exposes `NewVersion(string) (Version, error)`,
`(Version).Compare(Version) (int, error)`, and `(Version).String()`. The shared
core lives in `internal/core`.
