package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_regex(t *testing.T) {
	tests := []struct {
		pattern string
		s       string
		expect  bool
	}{
		{
			`\[backend\]`,
			"[szpp-judge] [backend] DB 設計会",
			true,
		},
		{
			`\[backend\]`,
			"[szpp-judge] [backed] DB 設計会",
			false,
		},
		{
			`\[judge\]`,
			"[szpp-judge] [backend] DB 設計会",
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.pattern, func(t *testing.T) {
			matched, err := regexp.MatchString(test.pattern, test.s)
			require.NoError(t, err)
			assert.Equal(t, test.expect, matched)
		})
	}
}
