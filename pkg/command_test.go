package allot

import (
	"testing"
)

var resultCommand bool

func BenchmarkMatches(b *testing.B) {
	var r bool

	for n := 0; n < b.N; n++ {
		r = New("command <lorem:integer> <ipsum:string>").Matches("command 12345 abcdef")
	}

	resultCommand = r
}

func TestMatches(t *testing.T) {
	var data = []struct {
		command string
		request string
		matches bool
	}{
		{"command", "example", false},
		{"command", "command", true},
		{"command        ", "command", true},
		{"command", "command         example", false},
		{"command <lorem>", "command", false},
		{"command <lorem>", "command example", true},
		{"command <lorem>", "command 1234567", true},
		{"command <lorem>", "command command", true},
		{"command <lorem>", "example command", false},
		{"command <lorem:integer>", "command example", false},
		{"command <lorem:integer>", "command 1234567", true},
		{"command <lorem>", "command command command", false},
		{"command <ipsum:integer?>", "command", true},
		{"command <ipsum:integer?>", "command 2", true},
		{"command <lorem:string?>", "command", true},
		{"command <lorem:string> <ipsum:string?>", "command 1234567", true},
		{"command <lorem:string> <ipsum:string?>", "command 1234567 test", true},
		{"command <lorem:remaining_string>", "command 1234567 test", true},
	}

	for _, set := range data {
		cmd := New(set.command)

		if cmd.Matches(set.request) != set.matches {
			t.Errorf("Matches() returns unexpected values. Got \"%v\", expected \"%v\"\nExpression: \"%s\" not matching \"%s\"",
				cmd.Matches(set.request), set.matches, cmd.Expression().String(), set.request)
		}
	}
}

func TestPosition(t *testing.T) {
	var data = []struct {
		command  string
		param    Parameter
		position int
	}{
		{"command <lorem>", NewParameterWithType("lorem", "string"), 0},
		{"command <lorem> <ipsum> <dolor> <sit> <amet>", NewParameterWithType("dolor", "string"), 2},
		{"command <lorem> <ipsum> <dolor:integer> <sit> <amet>", NewParameterWithType("dolor", "string"), -1},
		{"command <lorem:integer> <ipsum> <dolor:integer> <sit> <amet>", NewParameterWithType("lorem", "string"), -1},
		{"command <lorem:integer> <ipsum> <dolor:integer> <sit> <amet>", NewParameterWithType("lorem", "integer"), 0},
		{"command <lorem:integer> <ipsum> <lorem:string> <sit> <amet>", NewParameterWithType("lorem", "integer"), 0},
		{"command <lorem:integer> <ipsum> <lorem:string> <sit> <amet>", NewParameterWithType("lorem", "string"), 2},
	}

	var cmd Command
	for _, set := range data {
		cmd = *New(set.command)

		if cmd.Position(set.param) != set.position {
			t.Errorf("Position() should be \"%d\", but is \"%d\"", set.position, cmd.Position(set.param))
		}
	}
}

func TestHas(t *testing.T) {
	var data = []struct {
		command   string
		parameter Parameter
		has       bool
	}{
		{"command <lorem>", NewParameterWithType("lorem", "string"), true},
		{"command <lorem>", NewParameterWithType("lorem", "integer"), false},
	}

	var cmd CommandInterface
	for _, set := range data {
		cmd = New(set.command)

		if cmd.Has(set.parameter) != set.has {
			t.Errorf("HasParameter is \"%v\", expected \"%v\"", cmd.Has(set.parameter), set.has)
		}
	}
}

func TestParameters(t *testing.T) {
	var data = []struct {
		command    string
		parameters []Parameter
	}{
		{"command <lorem>", []Parameter{NewParameterWithType("lorem", "string")}},
		{"cmd <lorem:string>", []Parameter{NewParameterWithType("lorem", "string")}},
		{"command <lorem:integer>", []Parameter{NewParameterWithType("lorem", "integer")}},
		{"example <lorem> <ipsum> <dolor>", []Parameter{NewParameterWithType("lorem", "string"),
			NewParameterWithType("ipsum", "string"),
			NewParameterWithType("dolor", "string")}},
		{"command <lorem> <ipsum> <dolor:string>", []Parameter{NewParameterWithType("lorem", "string"),
			NewParameterWithType("ipsum", "string"),
			NewParameterWithType("dolor", "string")}},
		{"command <lorem> <ipsum:string> <dolor>", []Parameter{NewParameterWithType("lorem", "string"),
			NewParameterWithType("ipsum", "string"),
			NewParameterWithType("dolor", "string")}},
		{"command <lorem:string> <ipsum> <dolor>", []Parameter{NewParameterWithType("lorem", "string"),
			NewParameterWithType("ipsum", "string"),
			NewParameterWithType("dolor", "string")}},
		{"command <lorem:string> <ipsum> <dolor:string>", []Parameter{NewParameterWithType("lorem", "string"),
			NewParameterWithType("ipsum", "string"),
			NewParameterWithType("dolor", "string")}},
		{"command <lorem:string> <ipsum> <dolor:integer>", []Parameter{NewParameterWithType("lorem", "string"),
			NewParameterWithType("ipsum", "string"),
			NewParameterWithType("dolor", "integer")}},
		{"command <lorem:integer> <ipsum:string> <dolor:integer>", []Parameter{NewParameterWithType("lorem", "integer"),
			NewParameterWithType("ipsum", "string"),
			NewParameterWithType("dolor", "integer")}},
	}

	var cmd Command
	for _, set := range data {
		cmd = *New(set.command)

		if cmd.Text() != set.command {
			t.Errorf("cmd.Text() must be \"%s\", but is \"%s\"", set.command, cmd.Text())
		}

		for _, param := range set.parameters {
			if !cmd.Has(param) {
				t.Errorf("\"%s\" missing parameter.Item \"%s\"", cmd.Text(), param.Name())
			}
		}
	}
}

func TestTokenize(t *testing.T) {
	var data = []struct {
		command string
		tokens  []*Token
	}{
		{"command <lorem>", []*Token{NewTokenWithType("command", notParameter, 0), NewTokenWithType("lorem", definedParameter, 1)}},
		{"cmd <lorem:string>", []*Token{NewTokenWithType("cmd", notParameter, 0), NewTokenWithType("lorem:string", definedParameter, 1)}},
		{"cmd <lorem:string?>", []*Token{NewTokenWithType("cmd", notParameter, 0), NewTokenWithType("lorem:string?", optionalParameter, 1)}},
		{"cmd <lorem:integer?>", []*Token{NewTokenWithType("cmd", notParameter, 0), NewTokenWithType("lorem:integer?", optionalParameter, 1)}},
		{"cmd <lorem:?>", []*Token{NewTokenWithType("cmd", notParameter, 0), NewTokenWithType("lorem:?", optionalParameter, 1)}},
	}
	var cmd Command
	for _, set := range data {
		cmd = *New(set.command)

		if cmd.Text() != set.command {
			t.Errorf("cmd.Text() must be \"%s\", but is \"%s\"", set.command, cmd.Text())
		}

		tokens := cmd.Tokenize()
		for index, token := range set.tokens {
			if tokens[index].Word() != token.Word() {
				t.Errorf("\"%s\" missing token \"%s\"", cmd.Text(), token.Word())
			}
			if tokens[index].Type() != token.Type() {
				t.Errorf("for input: %s & test case: %s, %d token type mismatch %d", tokens[index].Word(), token.Word(), tokens[index].Type(), token.Type())
			}
		}
	}
}
