package commands_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	commands "github.com/tvanriel/discord-commands"
)

func TestStripPrefix(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "test", commands.StripPrefix("!", "ps")("!ps test"))
	assert.Equal(t, "test", commands.StripPrefix("!", "ps")("!ps       test"))
	assert.Equal(t, "test", commands.StripPrefix("!", "ps")("!pstest"))
}

func TestHasPrefix(t *testing.T) {
	t.Parallel()

	assert.True(t, commands.HasCommandPrefix("!", "ps", "!ps test"))
	assert.True(t, commands.HasCommandPrefix("!", "ps", "!ps"))
	assert.True(t, commands.HasCommandPrefix("!", "ps", "!ps      test"))
	assert.False(t, commands.HasCommandPrefix("!", "ps", "!pstest"))

	assert.False(t, commands.HasCommandPrefix("!", "ps", "!pslist"))
}

func TestSplitArgs(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		s    string
		want []string
	}{
		"empty string": {
			s:    "",
			want: []string{},
		},
		"one arg": {
			s:    "test",
			want: []string{"test"},
		},

		"two arg": {
			s:    "test one",
			want: []string{"test", "one"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := commands.SplitArgs(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
