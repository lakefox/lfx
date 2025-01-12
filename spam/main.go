package spam

import (
	"embed"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

//go:embed word_list.txt
var embeddedFiles embed.FS

var wf []string

// NewWordFilter initializes the WordFilter module by reading the banned words from the embedded file.
func Init() {
	// Read the embedded file
	data, err := embeddedFiles.ReadFile("word_list.txt")
	if err != nil {
		return
	}

	// Process the file content into a list of words
	words := []string{}
	for _, line := range strings.Split(string(data), "\n") {
		word := strings.TrimSpace(line)
		if word != "" {
			words = append(words, word)
		}
	}

	wf = words
}

// ContainsBannedWord checks if the given text contains any banned words.
func ContainsBannedWord(text string) bool {
	// Norm text to prent hiding
	text = norm.NFKC.String(text)
	text = removeDiacritics(text)
	text = normalizeString(text)
	for _, word := range wf {
		pattern := `(?i)` + regexp.QuoteMeta(word)
		matched, err := regexp.MatchString(pattern, text)
		if err != nil {
			return false
		}
		if matched {
			return true
		}
	}
	return false
}

// isNonSpacingMark identifies Unicode non-spacing marks (diacritics).
func isNonSpacingMark(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: Nonspacing_Mark
}

func removeDiacritics(input string) string {
	// Normalize to NFD (Normalization Form Decomposed).
	decomposed := norm.NFD.String(input)

	// Filter out non-spacing marks.
	result := make([]rune, 0, len(decomposed))
	for _, r := range decomposed {
		if !isNonSpacingMark(r) {
			result = append(result, r)
		}
	}

	return string(result)
}
func normalizeString(input string) string {
	// Remove spaces and convert to lowercase
	return strings.ReplaceAll(strings.ToLower(input), " ", "")
}
