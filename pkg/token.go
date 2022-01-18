package allot

import (
	"errors"
	"regexp"
	"strings"
)

const (
	definedOptionsPattern    = `\(.*?\)`
	definedParameterPattern  = `<(.*?)>`
	optionalParameterPattern = `<(.*?)[?]>`
	paramterPattern          = definedParameterPattern + "|" + definedOptionsPattern
	numberPattern            = `\d+`
)

const (
	notParameter = iota
	definedParameter
	definedOptionsParameter
	optionalParameter
)

// Token represents the Token object
type Token struct {
	word     string
	tType    int
	position int
}

// Word returns the token word
func (t Token) Word() string {
	return t.word
}

// Type returns the token word
func (t Token) Type() int {
	return t.tType
}

// Position returns the token word
func (t Token) Position() int {
	return t.position
}

func (t Token) IsParameter() bool {
	return t.tType != notParameter
}

// GetParameterFromToken return the parameter object created from token
func (t Token) GetParameterFromToken() (Parameter, error) {
	if t.IsParameter() {
		return Parse(t.Word(), t.Position()), nil
	}

	return Parameter{}, errors.New(t.Word() + " is not a parameter")
}

// tokenize returns array of the Tokens present in the command
func tokenize(line string) []*Token {
	definedParameterRegex := regexp.MustCompile(definedParameterPattern)
	definedOptionsRegex := regexp.MustCompile(definedOptionsPattern)
	optionalParameterRegex := regexp.MustCompile(optionalParameterPattern)
	words := strings.Fields(line)
	tokens := make([]*Token, len(words))
	for i, word := range words {
		switch {
		case optionalParameterRegex.MatchString(word):
			tokens[i] = NewTokenWithType(word[1:len(word)-1], optionalParameter, i)
		case definedParameterRegex.MatchString(word):
			tokens[i] = NewTokenWithType(word[1:len(word)-1], definedParameter, i)
		case definedOptionsRegex.MatchString(word):
			tokens[i] = NewTokenWithType(word[1:len(word)-1], definedOptionsParameter, i)
		default:
			tokens[i] = NewTokenWithType(word, notParameter, i)
		}
	}
	return tokens
}

// NewTokenWithType returns a Token
func NewTokenWithType(word string, tType int, tokenPosition int) *Token {
	return &Token{word, tType, tokenPosition}
}
