package allot

import (
	"testing"
)

func TestGetParameterFromToken(t *testing.T) {
	var data = []struct {
		token     *Token
		parameter Parameter
	}{
		{NewTokenWithType("cmd", notParameter, 0), NewParameterWithType("cmd", "")},
		{NewTokenWithType("lorem", definedParameter, 0), NewParameterWithType("lorem", "string")},
		{NewTokenWithType("lorem:string", definedParameter, 0), NewParameterWithType("lorem", "string")},
		{NewTokenWithType("lorem:string?", optionalParameter, 0), NewParameterWithType("lorem", "string?")},
		{NewTokenWithType("lorem:integer?", optionalParameter, 0), NewParameterWithType("lorem", "integer?")},
		{NewTokenWithType("lorem:?", optionalParameter, 0), NewParameterWithType("lorem", "string?")},
	}
	for _, set := range data {
		param, err := set.token.GetParameterFromToken()
		if err != nil && set.token.Type() != notParameter {
			t.Errorf("Cannot parse token: %s into parameter: %s", set.token.Word(), set.parameter.Name())
		} else if err == nil && param.Datatype() != set.parameter.Datatype() {
			t.Errorf("Expected parsed token type to be: %s but got: %s", set.parameter.Datatype(), param.Datatype())
		}
	}
}
