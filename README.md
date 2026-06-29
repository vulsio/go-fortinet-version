# go-fortinet-version

Fortinet product version parsing and comparison, split by version scheme:

- [`numeric`](./numeric) — the purely numeric scheme almost every product uses (`7.4.3`, train `7.2`).
- [`nonnumeric`](./nonnumeric) — the FortiSASE milestone-letter scheme (`25.2.a`, `25.1.a.2`).

Each package exposes `NewVersion(string) (Version, error)`, `(Version).Compare(Version) (int, error)`, and `(Version).String()`. The shared core lives in `internal/core`.
