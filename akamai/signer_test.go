package akamai

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testBuildTime(t *testing.T) {
	var ctx signingCtx
	ctx.buildTime()
	expected := regexp.MustCompile(`^\d{4}[0-1][0-9][0-3][0-9]T[0-2][0-9]:[0-5][0-9]:[0-5][0-9]\+0000$`)
	if assert.Regexp(t, expected, ctx.buildTime, "Fail: Regex do not match") {
		t.Log("Pass: Regex matches")
	}
}
