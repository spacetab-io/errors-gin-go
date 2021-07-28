package ginerrors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func Test_getHTTPCode(t *testing.T) {
	type tc struct {
		name   string
		exp    int
		in     codes.Code
		exists bool
	}

	tcs := []tc{
		{
			name:   "not found",
			in:     codes.NotFound,
			exp:    http.StatusNotFound,
			exists: true,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out, exist := getHTTPCode(tc.in)
			assert.Equal(t, tc.exp, out)
			assert.Equal(t, tc.exists, exist)
		})
	}
}

func Test_getGRPCCode(t *testing.T) {
	type tc struct {
		name   string
		in     int
		exp    codes.Code
		exists bool
	}

	tcs := []tc{
		{
			name:   "not found",
			in:     http.StatusNotFound,
			exp:    codes.NotFound,
			exists: true,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out, exist := getGRPCCode(tc.in)
			assert.Equal(t, tc.exp, out)
			assert.Equal(t, tc.exists, exist)
		})
	}
}
