package allot

import (
	"regexp"
	"strings"
)

const (
	definedOptionsPattern    = "\\(.*?\\)"
	definedParameterPattern  = "<(.*?)>"
	optionalParameterPattern = "<(.*?)[?]"
)

const (
	notParameter = iota
	definedParameter
	definedOptionsParameter
	optionalParameter
)

// Token represents the Token object
type Token struct {
	word string
	Type int
}

func (t Token) IsParameter() bool {
	return t.Type != notParameter
}

func tokenize(format string) []*Token {
	definedParameterRegex := regexp.MustCompile(definedParameterPattern)
	definedOptionsRegex := regexp.MustCompile(definedOptionsPattern)
	optionalParameterRegex := regexp.MustCompile(optionalParameterPattern)
	words := strings.Fields(format)
	tokens := make([]*Token, len(words))
	for i, word := range words {
		tWord := word[1 : len(word)-1]
		switch {
		case definedOptionsRegex.MatchString(word):
			tokens[i] = NewTokenWithType(tWord, definedOptionsParameter)
		case definedParameterRegex.MatchString(word):
			tokens[i] = NewTokenWithType(tWord, definedParameter)
		case optionalParameterRegex.MatchString(word):
			tokens[i] = NewTokenWithType(tWord, optionalParameter)
		default:
			tokens[i] = NewTokenWithType(word, notParameter)
		}
	}
	return tokens
}

// NewTokenWithType returns a Token
func NewTokenWithType(word string, tType int) *Token {
	return &Token{word, tType}
}
