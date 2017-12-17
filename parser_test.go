package yaraparser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/nbareil/yaraparser-go"
)

func TestParser_Test(t *testing.T) {
	var tests = []struct {
		s    string
		rule *yaraparser.YaraRule
		err  string
	}{
		// Single field statement
		{
			s: `rule foobar {}`,
			rule: &yaraparser.YaraRule{
				Name: "foobar",
			},
		},
		{
			s: `rule foobar { meta: author = "Roger"}`,
			rule: &yaraparser.YaraRule{
				Name:  "foobar",
				Metas: map[string]string{"author": "Roger"},
			},
		},
		{
			s: `rule foobar { strings: $a = "b"}`,
			rule: &yaraparser.YaraRule{
				Name: "foobar",
				Strings: map[string]yaraparser.YaraPattern{
					"$a": yaraparser.YaraPattern{
						Type:   yaraparser.RegularString,
						String: "b",
					},
				},
			},
		},
		{
			s: `rule foobar { meta: author = "Roger"}`,
			rule: &yaraparser.YaraRule{
				Name:  "foobar",
				Metas: map[string]string{"author": "Roger"},
			},
		},

		// quoted string
		{
			s: `rule foobar { meta: author = "she said \"hello\""}`,
			rule: &yaraparser.YaraRule{
				Name:  "foobar",
				Metas: map[string]string{"author": `she said \"hello\"`},
			},
		},

		// comment
		{
			s: `/* embedded comment */rule foobar {}`,
			rule: &yaraparser.YaraRule{
				Name: "foobar",
			},
		},
		{
			s: "// simple comment\nrule foobar {}",
			rule: &yaraparser.YaraRule{
				Name: "foobar",
			},
		},

		// Errors
		{s: `foo`, err: `found "foo", expected keyword "rule"`},
		{s: `rule * {}`, err: `found "*", expected a valid rule identifier`},
		{s: `rule foo {bar}`, err: `found "}", expecting ':'`},
		//{s: `rule foo {meta: bar}`, err: `found "}", expecting a value assignment`},
		{s: `rule foo : {}`, err: `invalid tag name`},
		{s: `rule foo : * {}`, err: `invalid tag name`},
		{s: `rule foo * {}`, err: `found "*" after rulename, expecting "{"`},
		{s: `rule !`, err: `found "!", expected a valid rule identifier`},
		{s: `rule foo { meta: author="foobar" `, err: `found EOF, expecting '}'`},
		{s: `rule foo { meta: author="foobar" XXX `, err: `found "", expecting a value assignment`},
	}

	for i, tt := range tests {
		rule, err := yaraparser.NewParser(strings.NewReader(tt.s)).Parse()
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.rule, rule) {
			t.Errorf("%d. %q\n\nYaraRule mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.rule, rule)
		}
	}
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
