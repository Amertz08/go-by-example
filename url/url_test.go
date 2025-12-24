package url

import (
	"fmt"
	"strings"
	"testing"
)

var parseTests = []struct {
	name string
	uri  string
	want *URL
}{
	{
		name: "with_data_scheme",
		uri:  "data:text/plain;base64:R28gYnkgRXhhbXBsZQ==",
		want: &URL{Scheme: "data"},
	},
	{
		name: "full",
		uri:  "https://github.com/inancgumus",
		want: &URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "inancgumus",
		},
	},
	{
		name: "without_path",
		uri:  "https://github.com",
		want: &URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "",
		},
	},
}

func TestURLString(t *testing.T) {
	tests := []struct {
		name string
		uri  *URL
		want string
	}{
		{
			name: "valid_case",
			uri: &URL{
				Scheme: "https",
				Host:   "github.com",
				Path:   "inancgumus",
			},
			want: "https://github.com/inancgumus",
		},
		{
			name: "nil",
			uri:  nil,
			want: "",
		},
		{name: "empty", uri: new(URL), want: ""},
		{
			name: "with_data_scheme",
			uri:  &URL{Scheme: "data"},
			want: "data:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.uri.String()
			if got != tt.want {
				t.Errorf("\ngot %q\nwant %q\n for %#v", got, tt.want, tt.uri)
			}
		})
	}
}

func TestParseSubtests(t *testing.T) {
	for _, tt := range parseTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.uri)
			if err != nil {
				t.Fatalf("Parse (%q) err = %q, want <nil>", tt.uri, err)
			}
			if *got != *tt.want {
				t.Errorf("Parse (%q)\ngot %#v\nwant %#v", tt.uri, got, tt.want)
			}
		})
	}
}

func TestParseError(t *testing.T) {
	tests := []struct {
		name string
		uri  string
	}{
		{name: "without_scheme", uri: "github.com"},
		{name: "empty_scheme", uri: "://github.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.uri)
			if err == nil {
				t.Errorf("Parse (%q) err=nil; wan an error", tt.uri)
			}
		})
	}
}

func BenchmarkURLString(b *testing.B) {
	u := &URL{Scheme: "https", Host: "github.com", Path: "inancgumus"}
	for b.Loop() {
		_ = u.String()
	}
}

func BenchmarkURLStringLong(b *testing.B) {
	for _, n := range []int{10, 100, 1_000} {
		u := &URL{
			Scheme: strings.Repeat("x", n),
			Host:   strings.Repeat("y", n),
			Path:   strings.Repeat("z", n),
		}
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			for b.Loop() {
				_ = u.String()
			}
		})
	}
}
