package allot

import "testing"

func TestMatch(t *testing.T) {
	var data = []struct {
		command  string
		request  string
		position int
		value    string
	}{
		{"command <param1:integer>", "command 1234", 0, "1234"},
		{"command <param1:remaining_string>", "command lorem ipsum", 0, "lorem ipsum"},
		{"(comm|Comm) <param1:string> <param2:integer>", "comm hello 2", 0, "comm"},
		{"(comm|Comm) <param1:string> <param2:integer>", "comm hello 2", 1, "hello"},
		{"(comm|Comm) <param1:string> <param2:integer>", "comm hello 2", 2, "2"},
		{"revert from <project:string> last <commits:integer> commits", "revert from example last 51 commits", 1, "51"},
		{"revert from <project:string> last <commits:integer> commits on (stage|prod)", "revert from example last 51 commits on stage", 2, "stage"},
		{"revert from <project:string> last <commits:integer> commits on (stage|prod)", "revert from example last 51 commits on prod", 2, "prod"},
		{"revert from <project:string?> last <commits:integer> commits", "revert from last 51 commits", 1, "51"},
	}

	for _, set := range data {
		match, err := New(set.command).Match(set.request)

		if err != nil {
			t.Errorf("Request [%s] does not match Command [%s]\n => %s", set.request, set.command, New(set.command).Expression().String())
		}

		value, err := match.Match(set.position)

		if err != nil {
			t.Errorf("Parsing parameter returned error: %v", err)
		}

		if value != set.value {
			t.Errorf("GetString() returned incorrect value. Got \"%v\", expected \"%v\"", value, set.value)
		}
	}
}

func TestMatchAndInteger(t *testing.T) {
	var data = []struct {
		command   string
		request   string
		parameter string
		value     int
	}{
		{"command <param1:integer>", "command 1234", "param1", 1234},
		{"revert from <project:string> last <commits:integer> commits", "revert from example last 51 commits", "commits", 51},
	}

	for _, set := range data {
		match, err := New(set.command).Match(set.request)

		if err != nil {
			t.Errorf("Request [%s] does not match Command [%s]\n => %s", set.request, set.command, New(set.command).Expression().String())
		}

		value, err := match.Integer(set.parameter)

		if err != nil {
			t.Errorf("Parsing parameter returned error: %v", err)
		}

		if value != set.value {
			t.Errorf("GetString() returned incorrect value. Got \"%d\", expected \"%d\"", value, set.value)
		}
	}
}

func TestMatchAndOptionalInteger(t *testing.T) {
	var data = []struct {
		command   string
		request   string
		parameter string
		value     int
	}{
		{"command <param1:integer?>", "command ", "param1", 0},
		{"command <param1:integer?>", "command 1234", "param1", 1234},
		{"revert from <project:string> last <commits:integer?> commits", "revert from example last commits", "commits", 0},
		{"revert from <project:string> last <commits:integer?> commits", "revert from example last 51 commits", "commits", 51},
	}

	for _, set := range data {
		match, err := New(set.command).Match(set.request)

		if err != nil {
			t.Errorf("Request [%s] does not match Command [%s]\n => %s", set.request, set.command, New(set.command).Expression().String())
		}

		value, _ := match.Integer(set.parameter)

		if value != set.value {
			t.Errorf("GetString() returned incorrect value. Got \"%d\", expected \"%d\"", value, set.value)
		}
	}
}

func TestMatchAndString(t *testing.T) {
	var data = []struct {
		command   string
		request   string
		parameter string
		value     string
	}{
		{"command <param1>", "command lorem", "param1", "lorem"},
		{"deploy <project:string> to <environment:string>", "deploy example to stage", "project", "example"},
		{"deploy <project:string> to <environment:string>", "deploy example to stage", "environment", "stage"},
		{"deploy <project:string> to <environment:string> at <host>", "deploy example to stage at api.exa!@#$%^&*()_mple.com", "host", "api.exa!@#$%^&*()_mple.com"},
		{"deploy <project:string> to <environment:string> at <host>", "deploy exam/ple to stage at api-prod.example.com", "host", "api-prod.example.com"},
		{"deploy <project:string> to <environment:string> at <host>", "deploy @klaus to stage at api-<test>.example.com", "project", "@klaus"},
		{"deploy <project:string> to <environment:string> at <host>", "deploy \"klaus\" to stage at api-<test>.example.com", "project", "\"klaus\""},
		{"deploy <project:string>-<stage:string> to <host>", "deploy klaus-prod to example", "project", "klaus"},
		{"deploy <project:string>-<stage:string> to <host>", "deploy klaus-prod to example", "stage", "prod"},
		{"deploy <project:string> to (stage|prod)", "deploy klaus to stage", "project", "klaus"},
		{"deploy <project:string> to (stage|prod)+", "deploy klaus to prod", "project", "klaus"},
		{"deploy to (stage|prod) at <host>", "deploy to stage at localhost", "option0", "stage"},
		{"deploy to (stage|prod) at <host>", "deploy to prod at localhost", "option0", "prod"},
		{"deploy to (stage|prod) at <host>", "deploy to prod at localhost", "host", "localhost"},
		{"deploy to (stage|prod) at (localhost|server)", "deploy to prod at server", "option1", "server"},
	}

	for _, set := range data {
		match, err := New(set.command).Match(set.request)

		if err != nil {
			t.Errorf("Request [%s] does not match Command [%s]", set.request, set.command)
		}

		value, err := match.String(set.parameter)

		if err != nil {
			t.Errorf("Parsing parameter returned error: %v", err)
		}

		if value != set.value {
			t.Errorf("GetString() returned incorrect value. Got \"%s\", expected \"%s\"", value, set.value)
		}
	}
}

func TestMatchAndOptionalString(t *testing.T) {
	var data = []struct {
		command   string
		request   string
		parameter string
		value     string
	}{
		{"command <param1:?>", "command ", "param1", ""},
		{"command <param1:?>", "command lorem", "param1", "lorem"},
		{"command <param1:string?>", "command 1", "param1", "1"},
		{"deploy <project:string> to <environment:string?>", "deploy example to stage", "environment", "stage"},
	}

	for _, set := range data {
		match, err := New(set.command).Match(set.request)

		if err != nil {
			t.Errorf("Request [%s] does not match Command [%s]", set.request, set.command)
		}

		value, err := match.String(set.parameter)

		if err != nil {
			t.Errorf("Parsing parameter returned error: %v", err)
		}

		if value != set.value {
			t.Errorf("GetString() returned incorrect value. Got \"%s\", expected \"%s\"", value, set.value)
		}
	}
}
