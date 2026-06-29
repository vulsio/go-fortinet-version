package version_test

import (
	"errors"
	"testing"

	version "github.com/vulsio/go-fortinet-version/nonnumeric"
)

func TestNewVersion(t *testing.T) {
	tests := []struct {
		name    string
		ver     string
		want    string // canonical String(); checked when !wantErr
		wantErr bool
	}{
		{name: "numeric", ver: "7.4.3", want: "7.4.3"},
		{name: "milestone letter", ver: "25.2.a", want: "25.2.a"},
		{name: "milestone letter patch", ver: "25.1.a.2", want: "25.1.a.2"},
		{name: "empty", ver: "", wantErr: true},
		{name: "consecutive dots", ver: "7..0", wantErr: true},
		{name: "trailing dot", ver: "7.2.", wantErr: true},
		{name: "signed component", ver: "7.-1", wantErr: true},
		{name: "plus-signed component", ver: "7.+0", wantErr: true},
		{name: "multi-char milestone", ver: "25.2.alpha", wantErr: true},
		{name: "letter+digits milestone", ver: "25.1.a10", wantErr: true},
		{name: "build suffix", ver: "7.1-b5955", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := version.NewVersion(tt.ver)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewVersion(%q) error = %v, wantErr %v", tt.ver, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if got := v.String(); got != tt.want {
				t.Errorf("NewVersion(%q).String() = %q, want %q", tt.ver, got, tt.want)
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		name    string
		v1, v2  string
		want    int
		wantErr bool // expect ErrIncomparable
	}{
		{name: "numeric v1 < v2", v1: "7.0.0", v2: "7.0.1", want: -1},
		{name: "trailing zero stays equal (7.2.0 == 7.2)", v1: "7.2.0", v2: "7.2", want: 0},
		{name: "bare train below its milestone build (25.2 < 25.2.a)", v1: "25.2", v2: "25.2.a", want: -1},
		{name: "milestone above its bare train (25.2.a > 25.2)", v1: "25.2.a", v2: "25.2", want: 1},
		{name: "milestone below next train (25.2.a < 25.3)", v1: "25.2.a", v2: "25.3", want: -1},
		{name: "milestone above prev train (25.2.a > 25.1)", v1: "25.2.a", v2: "25.1", want: 1},
		{name: "sequential letters (25.2.a < 25.2.b)", v1: "25.2.a", v2: "25.2.b", want: -1},
		{name: "letter patch after letter (25.1.a < 25.1.a.2)", v1: "25.1.a", v2: "25.1.a.2", want: -1},
		{name: "numeric letter-patch ordering (25.1.a.2 < 25.1.a.10)", v1: "25.1.a.2", v2: "25.1.a.10", want: -1},
		{name: "nested milestone above bare train (25.1.a.2 > 25.1)", v1: "25.1.a.2", v2: "25.1", want: 1},
		{name: "numeric vs milestone at same position → incomparable (1.2.1 vs 1.2.a)", v1: "1.2.1", v2: "1.2.a", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := version.NewVersion(tt.v1)
			if err != nil {
				t.Fatalf("NewVersion(%q): %v", tt.v1, err)
			}
			v2, err := version.NewVersion(tt.v2)
			if err != nil {
				t.Fatalf("NewVersion(%q): %v", tt.v2, err)
			}
			got, err := v1.Compare(v2)
			if tt.wantErr {
				if !errors.Is(err, version.ErrIncomparable) {
					t.Fatalf("Compare(%q, %q) error = %v, want ErrIncomparable", tt.v1, tt.v2, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Compare(%q, %q) unexpected error: %v", tt.v1, tt.v2, err)
			}
			if got != tt.want {
				t.Errorf("Compare(%q, %q) = %d, want %d", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

func TestVersion_CompareZeroValue(t *testing.T) {
	v, err := version.NewVersion("25.2.a")
	if err != nil {
		t.Fatalf("NewVersion: %v", err)
	}
	if _, err := (version.Version{}).Compare(v); err == nil {
		t.Error("Compare with zero-value receiver: want error, got nil")
	}
	if _, err := v.Compare(version.Version{}); err == nil {
		t.Error("Compare against zero-value arg: want error, got nil")
	}
}
