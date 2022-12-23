package sanitizer

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Sanitize(s string) string {
	// Remove leading/trailing spaces
	s = strings.TrimSpace(s)
	// Converts letters with tildes to normal
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ = transform.String(t, s)
	// Convert string to lowercase
	s = strings.ToLower(s)

	return s
}
