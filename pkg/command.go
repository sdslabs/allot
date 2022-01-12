package allot

import (
	"errors"
	"regexp"
	"strings"
)

// CommandInterface describes how to access a Command
type CommandInterface interface {
	Expression() *regexp.Regexp
	Has(name ParameterInterface) bool
	Match(req string) (MatchInterface, error)
	Matches(req string) bool
	Parameters() []Parameter
	Position(param ParameterInterface) int
	Text() string
	Tokenize() []*Token
}

// Command is a Command definition
type Command struct {
	text string
}

// Text returns the command text
func (c Command) Text() string {
	return c.text
}

// Expression returns the regular expression matching the command text
func (c Command) Expression() *regexp.Regexp {
	expr := c.Text()
	expr = strings.TrimSpace(expr)
	expr = removeExtraWhitespaces(expr)

	for _, param := range c.Parameters() {
		newString := param.Expression().String()

		oldString1 := "<" + param.Name() + ":" + param.Datatype() + ">"
		if param.IsOptional() {
			oldString1 = WhitespaceCharacter + oldString1
		}
		expr = strings.Replace(expr, oldString1, newString, -1)

		oldString2 := "<" + param.Name() + ">"
		expr = strings.Replace(expr, oldString2, newString, -1)

		oldString3 := "<" + param.Name() + ":?>"
		if param.IsOptional() {
			oldString3 = WhitespaceCharacter + oldString3
		}
		expr = strings.Replace(expr, oldString3, newString, -1)
	}

	return regexp.MustCompile("^" + expr + "$")
}

// Parameters returns the list of defined parameters
func (c Command) Parameters() []Parameter {
	var list []Parameter
	re := regexp.MustCompile(paramterPattern)
	result := re.FindAllStringSubmatch(c.Text(), -1)

	for listIndex, p := range result {
		if len(p) != 2 {
			continue
		}

		list = append(list, Parse(p[0], listIndex))
	}

	return list
}

// Has checks if the parameter is found in the command
func (c Command) Has(param ParameterInterface) bool {
	return c.Position(param) != -1
}

// Position returns the position of a parameter
func (c Command) Position(param ParameterInterface) int {
	for index, item := range c.Parameters() {
		if item.Equals(param) {
			return index
		}
	}

	return -1
}

// Match returns the parameter matching the expression at the defined position
func (c Command) Match(req string) (MatchInterface, error) {
	req = removeExtraWhitespaces(req)
	req = strings.TrimSpace(req)

	if c.Matches(req) {
		return Match{c, req}, nil
	}

	return nil, errors.New("request does not match command")
}

// Matches checks if a comand definition matches a request
func (c Command) Matches(req string) bool {
	return c.Expression().MatchString(strings.TrimSpace(req))
}

// Tokenize returns Command info as tokens
func (c Command) Tokenize() []*Token {
	return tokenize(c.text)
}

// New returns a new command
func New(command string) *Command {
	return &Command{command}
}
