package numeric_test

import (
	"testing"

	"github.com/vulsio/go-fortinet-version/numeric"
)

func TestNewVersion(t *testing.T) {
	tests := []struct {
		name    string
		ver     string
		want    string // canonical String(); checked when !wantErr
		wantErr bool
	}{
		{name: "concrete", ver: "7.4.3", want: "7.4.3"},
		{name: "train minor", ver: "7.2", want: "7.2"},
		{name: "train major", ver: "7", want: "7"},
		{name: "trailing zero kept", ver: "7.2.0", want: "7.2.0"},
		{name: "milestone letter rejected", ver: "25.2.a", wantErr: true},
		{name: "empty", ver: "", wantErr: true},
		{name: "consecutive dots", ver: "7..0", wantErr: true},
		{name: "trailing dot", ver: "7.2.", wantErr: true},
		{name: "signed component", ver: "7.-1", wantErr: true},
		{name: "plus-signed component", ver: "7.+0", wantErr: true},
		{name: "build suffix", ver: "7.1-b5955", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := numeric.NewVersion(tt.ver)
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
		name   string
		v1, v2 string
		want   int
	}{
		{name: "v1 < v2", v1: "7.0.0", v2: "7.0.1", want: -1},
		{name: "equal", v1: "7.2.0", v2: "7.2.0", want: 0},
		{name: "v1 > v2", v1: "7.1.0", v2: "7.0.0", want: 1},
		{name: "train minor < train minor", v1: "7.2", v2: "7.4", want: -1},
		{name: "train major < train major", v1: "7", v2: "8", want: -1},
		{name: "concrete within train lower bound (7.2.0 == 7.2)", v1: "7.2.0", v2: "7.2", want: 0},
		{name: "concrete below next train (7.2.5 < 7.3)", v1: "7.2.5", v2: "7.3", want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := numeric.NewVersion(tt.v1)
			if err != nil {
				t.Fatalf("NewVersion(%q): %v", tt.v1, err)
			}
			v2, err := numeric.NewVersion(tt.v2)
			if err != nil {
				t.Fatalf("NewVersion(%q): %v", tt.v2, err)
			}
			got, err := v1.Compare(v2)
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
	v, err := numeric.NewVersion("7.0.0")
	if err != nil {
		t.Fatalf("NewVersion: %v", err)
	}
	if _, err := (numeric.Version{}).Compare(v); err == nil {
		t.Error("Compare with zero-value receiver: want error, got nil")
	}
	if _, err := v.Compare(numeric.Version{}); err == nil {
		t.Error("Compare against zero-value arg: want error, got nil")
	}
}
