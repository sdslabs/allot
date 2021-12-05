package allot

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	definedOptionsPattern    = "\\(.*?\\)"
	definedParameterPattern  = "<(.*?)>"
	optionalParameterPattern = "<(.*?)[?]>"
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

// GetParameterFromToken return the parameter object created from token
func (t Token) GetParameterFromToken() (Parameter, error) {
	if t.IsParameter() {
		return Parse(t.word), nil
	}

	return Parameter{}, errors.New(t.word + " is not a parameter")
}

// tokenize returns array of the Tokens present in the command
func tokenize(line string) []*Token {
	definedParameterRegex := regexp.MustCompile(definedParameterPattern)
	definedOptionsRegex := regexp.MustCompile(definedOptionsPattern)
	optionalParameterRegex := regexp.MustCompile(optionalParameterPattern)
	words := strings.Fields(line)
	tokens := make([]*Token, len(words))
	for i, word := range words {
		tWord := word[1 : len(word)-1]
		fmt.Println(word, definedParameterRegex.MatchString(word))
		switch {
		case optionalParameterRegex.MatchString(word):
			tokens[i] = NewTokenWithType(tWord, optionalParameter)
		case definedParameterRegex.MatchString(word):
			tokens[i] = NewTokenWithType(tWord, definedParameter)
		case definedOptionsRegex.MatchString(word):
			tokens[i] = NewTokenWithType(tWord, definedOptionsParameter)
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
