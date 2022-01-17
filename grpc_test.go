package ginerrors_test

import (
	"database/sql"
	"testing"

	ginerrors "github.com/spacetab-io/errors-go-gin"
	"github.com/stretchr/testify/assert"
)

const (
	unknownErrMessage        = "rpc error: code = InvalidArgument desc = unknown error value"
	recordNotFountMessage    = "rpc error: code = NotFound desc = record not found"
	unavailableMethodMessage = "rpc error: code = Unavailable desc = method not allowed"
)

func TestWrapErrWithRPC(t *testing.T) {
	cases := []struct {
		caseName string
		err      error
		lang     string
		expected string
	}{
		{
			caseName: "unknown err",
			err:      ginerrors.ErrUnknownErrVal,
			lang:     "en",
			expected: unknownErrMessage,
		},
		{
			caseName: "sql err",
			err:      sql.ErrNoRows,
			lang:     "en",
			expected: recordNotFountMessage,
		},
		{
			caseName: "unavailable method",
			err:      ginerrors.ErrNoMethod,
			lang:     "en",
			expected: unavailableMethodMessage,
		},
	}

	t.Parallel()

	for _, cc := range cases {
		cc := cc
		t.Run(cc.caseName, func(t *testing.T) {
			t.Parallel()

			err := ginerrors.WrapErrorWithStatus(cc.err, cc.lang)
			assert.Error(t, err)
			assert.Equal(t, cc.expected, err.Error())
		})
	}
}
