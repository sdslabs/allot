package allot

import "regexp"

// removeExtraWhitespaces converts two or more whitespaces to one whitespace
func removeExtraWhitespaces(text string) string {
	return regexp.MustCompile(WhitespaceRegex).ReplaceAllString(text, WhitespaceCharacter)

}
