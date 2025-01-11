package spam

import (
	"embed"
	"regexp"
	"strings"
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
