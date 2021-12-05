package allot

import (
	"testing"
)

func TestGetParameterFromToken(t *testing.T) {
	var data = []struct {
		token     *Token
		parameter Parameter
	}{
		{NewTokenWithType("cmd", notParameter), NewParameterWithType("cmd", "")},
		{NewTokenWithType("lorem", definedParameter), NewParameterWithType("lorem", "string")},
		{NewTokenWithType("lorem:string", definedParameter), NewParameterWithType("lorem", "string")},
		{NewTokenWithType("lorem:string?", optionalParameter), NewParameterWithType("lorem", "string?")},
		{NewTokenWithType("lorem:integer?", optionalParameter), NewParameterWithType("lorem", "integer?")},
		{NewTokenWithType("lorem:?", optionalParameter), NewParameterWithType("lorem", "string?")},
	}
	for _, set := range data {
		param, err := set.token.GetParameterFromToken()
		if err != nil && set.token.Type != notParameter {
			t.Errorf("Cannot parse token: %s into parameter: %s", set.token.Word, set.parameter.name)
		} else if err == nil && param.datatype != set.parameter.datatype {
			t.Errorf("Expected parsed token type to be: %s but got: %s", set.parameter.datatype, param.datatype)
		}
	}
}
