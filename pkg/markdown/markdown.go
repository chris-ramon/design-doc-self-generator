package markdown

import (
	"regexp"
	"strings"
)

// StripMarkdown removes Markdown formatting, HTML tags, and blank lines.
func StripMarkdown(input string) string {
	// Remove Markdown links: [text](url)
	re := regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
	output := re.ReplaceAllString(input, "$1")

	// Remove images: ![alt](url)
	re = regexp.MustCompile(`!\[(.*?)\]\(.*?\)`)
	output = re.ReplaceAllString(output, "$1")

	// Remove emphasis, bold, strikethrough, inline code
	replacements := []string{"**", "", "*", "", "__", "", "_", "", "~~", "", "`", ""}
	for i := 0; i < len(replacements); i += 2 {
		output = strings.ReplaceAll(output, replacements[i], replacements[i+1])
	}

	// Remove headers: #, ##, etc.
	re = regexp.MustCompile(`(?m)^#{1,6}\s*`)
	output = re.ReplaceAllString(output, "")

	// Remove blockquotes: >
	re = regexp.MustCompile(`(?m)^>\s?`)
	output = re.ReplaceAllString(output, "")

	// Remove horizontal rules: --- or ***
	re = regexp.MustCompile(`(?m)^[-*_]{3,}\s*$`)
	output = re.ReplaceAllString(output, "")

	// Remove code blocks: ```...```
	re = regexp.MustCompile("(?s)```.*?```")
	output = re.ReplaceAllString(output, "")

	// Remove HTML tags: <tag>...</tag> or <tag/>
	re = regexp.MustCompile(`</?[^>]+?>`)
	output = re.ReplaceAllString(output, "")

	// Split, remove blank lines, and join with space
	lines := strings.Split(output, "\n")
	var nonEmptyLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	return strings.Join(nonEmptyLines, " ")
}
