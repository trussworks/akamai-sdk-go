package akamai

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTime(t *testing.T) {
	var ctx signingCtx
	ctx.buildTime()
	time := ctx.formattedTime
	expected := regexp.MustCompile(`^\d{4}[0-1][0-9][0-3][0-9]T[0-2][0-9]:[0-5][0-9]:[0-5][0-9]\+0000$`)
	if assert.Regexp(t, expected, time, "Fail: Regex do not match") {
		t.Log("Pass: Regex matches")
	}
}

func TestBuildNonce(t *testing.T) {
	var ctx signingCtx
	ctx.buildNonce()
	first := ctx.nonce
	for i := 0; i < 100; i++ {
		ctx.buildNonce()
		second := ctx.nonce
		assert.NotEqual(t, first, second, "Fail: Nonce matches")
	}
}
