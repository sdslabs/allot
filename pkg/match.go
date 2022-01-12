package allot

import (
	"errors"
	"fmt"
	"strconv"
)

// MatchInterface describes how to access a Match
type MatchInterface interface {
	String(name string) (string, error)
	Integer(name string) (int, error)
	RemainingString(text string) (string, error)
	Match(position int) (string, error)

	Parameter(param ParameterInterface) (string, error)
}

// Match is the Match definition
type Match struct {
	Command CommandInterface
	Request string
}

// String returns the value for a string parameter
func (m Match) String(name string) (string, error) {
	return m.Parameter(NewParameterWithType(name, StringType))
}

// String returns the value for a remaining string parameter
func (m Match) RemainingString(name string) (string, error) {
	return m.Parameter(NewParameterWithType(name, RemaingStringType))
}

// Integer returns the value for an integer parameter
func (m Match) Integer(name string) (int, error) {
	str, err := m.Parameter(NewParameterWithType(name, IntegerType))

	if str == "" {
		return 0, errors.New("value not provided")
	}
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}

// Parameter returns the value for a parameter
func (m Match) Parameter(param ParameterInterface) (string, error) {
	pos := m.Command.Position(param)
	if pos == -1 {
		return "", errors.New("Unknown parameter \"" + param.Name() + "\"")
	}

	matches := m.Command.Expression().FindAllStringSubmatch(m.Request, -1)[0][1:]
	return matches[pos], nil
}

// Match returns the match at given position
func (m Match) Match(position int) (string, error) {
	matches := m.Command.Expression().FindAllStringSubmatch(m.Request, -1)

	if len(matches) != 1 {
		return "", errors.New("unable to parse request")
	}

	if position+1 >= len(matches[0]) {
		return "", fmt.Errorf("no parameter at position %d", position)
	}

	return matches[0][position+1], nil
}
