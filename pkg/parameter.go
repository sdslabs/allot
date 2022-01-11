package allot

import (
	"regexp"
	"strconv"
	"strings"
)

var regexpMapping = map[string]string{
	REMAINING_STRING: "([\\s\\S]*)",
	STRING:           "([^\\s]+)",
	OPTIONAL_STRING:  "([^\\s]+)?",
	INTEGER:          "([0-9]+)",
	OPTIONAL_INTEGER: "([0-9]+)?",
}

// GetRegexpExpression returns the regexp for a data type
func GetRegexpExpression(datatype string) *regexp.Regexp {
	if exp, ok := regexpMapping[datatype]; ok {
		return regexp.MustCompile(exp)
	}

	return nil
}

// ParameterInterface describes how to access a Parameter
type ParameterInterface interface {
	Equals(param ParameterInterface) bool
	Expression() *regexp.Regexp
	Name() string
	Datatype() string
}

// Parameter is the Parameter definition
type Parameter struct {
	name     string
	datatype string
	expr     *regexp.Regexp
}

// Expression returns the regexp behind the type
func (p Parameter) Expression() *regexp.Regexp {
	return p.expr
}

// Name returns the Parameter name
func (p Parameter) Name() string {
	return p.name
}

// Datatype returns the Parameter datatype
func (p Parameter) Datatype() string {
	return p.datatype
}

// Equals checks if two parameter are equal
func (p Parameter) Equals(param ParameterInterface) bool {
	return p.Name() == param.Name() && strings.Contains(p.Expression().String(), param.Expression().String())
}

// NewParameterWithType returns a Parameter
func NewParameterWithType(name string, datatype string) Parameter {
	return Parameter{name, datatype, GetRegexpExpression(datatype)}
}

// Parse parses parameter info
func Parse(token string, paramterPosition int) Parameter {
	definedParameterRegex := regexp.MustCompile(definedParameterPattern)
	definedOptionsRegex := regexp.MustCompile(definedOptionsPattern)
	var name, datatype string

	switch {

	case definedParameterRegex.MatchString(token):
		name, datatype = parseDefinedParameterType(token)
	case definedOptionsRegex.MatchString(token):
		name, datatype = parseDefinedOptionsParameterType(token, paramterPosition)
	default:
		name, datatype = parseParamterType(token)
	}

	return NewParameterWithType(name, datatype)
}

func parseDefinedParameterType(token string) (string, string) {
	tokenWithoutAngleBrackets := token[1 : len(token)-1]
	return parseParamterType(tokenWithoutAngleBrackets)
}

func parseDefinedOptionsParameterType(token string, paramterPosition int) (string, string) {
	tokenWithoutCurlyBrackets := token[1 : len(token)-1]
	numberPatternRegex := regexp.MustCompile(numberPattern)
	datatype := "string"
	name := "option" + strconv.Itoa(paramterPosition)

	if numberPatternRegex.MatchString(tokenWithoutCurlyBrackets) {
		datatype = "integer"
	}

	return name, datatype
}

func parseParamterType(token string) (string, string) {
	datatype := "string"
	name := token
	if strings.Contains(token, ":") {
		splits := strings.Split(token, ":")
		if splits[1] == "?" {
			datatype = "string?"
		} else {
			datatype = splits[1]
		}
		name = splits[0]
	}

	return name, datatype
}
