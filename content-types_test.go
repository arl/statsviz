package statsviz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentTypeByUrl(t *testing.T) {
	is := assert.New(t)

	type testCase struct {
		in  string
		exp string
	}
	cases := []testCase{
		{
			"js/app.js",
			"text/javascript",
		},
		{
			in:  "libs/js/tippy.js?asd&fg=h",
			exp: "text/javascript",
		},
		{
			in:  "libs/js/tippy.js@6",
			exp: "text/javascript",
		},
	}
	for _, tc := range cases {
		res := urlContentType(tc.in)
		is.Equal(tc.exp, res, tc.in)
	}
}

func TestStrictExt(t *testing.T) {
	is := assert.New(t)

	type testCase struct {
		in  string
		exp string
	}
	cases := []testCase{
		{
			".js",
			".js",
		},
		{
			".js@6",
			".js",
		},
		{
			".js#asd",
			".js",
		},
		{
			" .js ",
			".js",
		},
	}
	for _, tc := range cases {
		res := strictFileExt(tc.in)
		is.Equal(tc.exp, res, tc.in)
	}
}
