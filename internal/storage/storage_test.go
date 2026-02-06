package storage

import (
	"strings"
	"testing"
)

func TestSanitizeFilename(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in   string
		want string
	}{
		{in: " book.epub ", want: "book.epub"},
		{in: "../a/b/c.mobi", want: "c.mobi"},
		{in: "a\\b\\c.azw3", want: "a_b_c.azw3"},
		{in: "a/b:c.pdf", want: "b_c.pdf"},
		{in: "....", want: ""},
		{in: "\x00evil.epub", want: "evil.epub"},
		{in: "a\nb.epub", want: "ab.epub"},
	}

	for _, tc := range cases {
		got := SanitizeFilename(tc.in)
		if got != tc.want {
			t.Fatalf("SanitizeFilename(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestSanitizeFilename_MaxLen(t *testing.T) {
	t.Parallel()

	input := strings.Repeat("a", 300) + ".epub"
	got := SanitizeFilename(input)
	if got == "" {
		t.Fatalf("expected non-empty sanitized filename")
	}
	if len(got) > 200 {
		t.Fatalf("expected sanitized filename <= 200, got %d", len(got))
	}
}
