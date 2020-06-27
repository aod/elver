package flags_test

import (
	"flag"
	"io/ioutil"
	"testing"

	"github.com/aod/elver/flags"
)

func TestIntRangeFlagSet(t *testing.T) {
	fs := flag.NewFlagSet("intrange", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)
	fs.Var(&flags.IntRange{Value: 0, Min: 1, Max: 10}, "num", "select a number between 1 to 10")

	testCases := []struct {
		args []string
		desc string
		ok   bool
	}{
		{
			args: []string{},
			desc: "No args",
			ok:   true,
		},
		{
			args: []string{"-num", ""},
			desc: "Empty string",
			ok:   false,
		},
		{
			args: []string{"-num", "abc"},
			desc: "Not a number",
			ok:   false,
		},
		{
			args: []string{"-num", "0"},
			desc: "Below Min range",
			ok:   false,
		},
		{
			args: []string{"-num", "11"},
			desc: "Above Max range",
			ok:   false,
		},
		{
			args: []string{"-num", "5"},
			desc: "Valid value",
			ok:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := fs.Parse(tc.args)
			if tc.ok && err != nil {
				t.Error(err)
			}
		})
	}
}
