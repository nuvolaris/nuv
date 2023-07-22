package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func testNoopCheck(p *preflight) {}
func testErrCheck(p *preflight) {
	p.err = errors.New("test error in check")
}

func Test_preflight(t *testing.T) {
	p := preflight{}
	p.check(testNoopCheck)
	require.NoError(t, p.err)

	p.err = errors.New("test error")
	p.check(testNoopCheck)
	require.Error(t, p.err)
	require.Equal(t, "test error", p.err.Error())

	p.err = nil
	p.check(testErrCheck)
	require.Error(t, p.err)
	require.Equal(t, "test error in check", p.err.Error())
}
