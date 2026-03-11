package main

import (
	"errors"
	"os"
	"testing"

	"github.com/4sp1/snippi/internal/template"
)

type expectedParseFlags struct {
	vars map[template.VarName]template.Value
	file string
	err  error
}

func Test_parseFlags(t *testing.T) {
	for _, test := range []struct {
		name     string
		args     []string
		expected expectedParseFlags
	}{
		{
			name: "unpaired flags",
			args: []string{"", "-template"},
			expected: expectedParseFlags{
				err: errNonMatchingFlags,
			},
		},
		{
			name: "no args",
			args: []string{""},
			expected: expectedParseFlags{
				err: errNotEnoughArgs,
			},
		},
		{
			name: "zero flag",
			args: []string{"", "-template", "hi", "-", "there"},
			expected: expectedParseFlags{
				err: errZeroFlagName,
			},
		},
		{
			name: "zero flag",
			args: []string{"", "-template", "hi", "yo", "there"},
			expected: expectedParseFlags{
				err: errIncorrectflagPrefix,
			},
		},
		{
			name: "zero flag",
			args: []string{"", "-hi", "there"},
			expected: expectedParseFlags{
				err: errMissingTemplateFlag,
			},
		},
		{
			name: "common",
			args: []string{
				"",
				"-template", "some.gotmpl",
				"-hello", "world",
			},
			expected: expectedParseFlags{
				vars: map[template.VarName]template.Value{
					"hello": "world",
				},
				file: "some.gotmpl",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			os.Args = test.args
			vars, file, err := parseFlags()
			if err != nil {
				if !errors.Is(err, test.expected.err) {
					fatal(t, "err", err, test.expected.err)
				}
				return
			}
			if test.expected.err != nil {
				fatal(t, "err", err, test.expected.err)
			}
			if len(vars) != len(test.expected.vars) {
				fatal(t, "vars", vars, test.expected.vars)
			}
			for k := range vars {
				if _, found := test.expected.vars[k]; !found {
					fatal(t, "vars", vars, test.expected.vars)
				}
			}
			for k := range test.expected.vars {
				if _, found := vars[k]; !found {
					fatal(t, "vars", vars, test.expected.vars)
				}
			}
			if file != test.expected.file {
				fatal(t, "file", file, test.expected.file)
			}
		})
	}
}

func fatal(t *testing.T, prefix string, got, expected any) {
	t.Helper()
	t.Fatalf("%ss: \ngot:\n%q\nexpected:\n%q", prefix, got, expected)

}
