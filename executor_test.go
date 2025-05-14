package commands_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tvanriel/discord-commands"
	"go.uber.org/zap"
)

type mockCmd struct {
	mock.Mock
}

// Apply implements commands.Command.
func (m *mockCmd) Apply(ctx *commands.Context) error {
	return m.Called(ctx).Error(0)
}

// Name implements commands.Command.
func (m *mockCmd) Name() string {
	return m.Called().String(0)
}

// SkipsPrefix implements commands.Command.
func (m *mockCmd) SkipsPrefix() bool {
	return m.Called().Bool(0)
}

var _ commands.Command = (*mockCmd)(nil)

func MustCall(t *testing.T, hasprefix bool, incantation string) *mockCmd {
	t.Helper()

	m := new(mockCmd)

	m.Test(t)
	m.On("Apply", mock.Anything).Return(nil).NotBefore(
		m.On("SkipsPrefix").Return(hasprefix),
		m.On("Name").Return(incantation),
	)

	return m
}

func MaybeCall(t *testing.T, hasprefix bool, incantation string) *mockCmd {
	t.Helper()

	m := new(mockCmd)

	m.Test(t)

	m.On("Apply", mock.Anything).Return(nil).Maybe().NotBefore(
		m.On("SkipsPrefix").Return(hasprefix),
		m.On("Name").Return(incantation),
	)

	return m
}

func DontCall(t *testing.T, hasprefix bool, incantation string) *mockCmd {
	t.Helper()

	m := new(mockCmd)

	m.Test(t)
	m.On("SkipsPrefix").Return(hasprefix)
	m.On("Name").Return(incantation)
	m.AssertNotCalled(t, "Apply", mock.Anything)

	return m
}

//nolint:ireturn
func asCommand(m *mockCmd) commands.Command {
	return m
}

func TestExecutor_HasMatch(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		// Named input parameters for receiver constructor.
		commands func(t *testing.T) []commands.Command
		// Named input parameters for target function.
		trigger string
		message string
		want    bool
	}{
		"empty message triggers nothing": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(DontCall(t, false, "I shall never be called")),
				}
			},
			trigger: "!",
			message: "",
			want:    false,
		},
		"not matching prefix doesn't match": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(DontCall(t, false, "test")),
				}
			},
			trigger: "!",
			message: "$test",
			want:    false,
		},
		"not matching command doesn't match": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(DontCall(t, false, "test")),
				}
			},
			trigger: "!",
			message: "!nope",
			want:    false,
		},
		"matching command matches": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(DontCall(t, false, "test")),
				}
			},
			trigger: "!",
			message: "!test",
			want:    true,
		},
		"empty trigger matches okay": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(DontCall(t, false, "test")),
				}
			},
			trigger: "",
			message: "test",
			want:    true,
		},
		"empty trigger matdoesnt match when unmatched command": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(DontCall(t, false, "test")),
				}
			},
			trigger: "",
			message: "nope",
			want:    false,
		},
		"literal command matches": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(DontCall(t, true, "Test")),
				}
			},
			trigger: "",
			message: "Test",
			want:    true,
		},
		"literal command ignores prefix": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(DontCall(t, true, "Test")),
				}
			},
			trigger: "",
			message: "!Test",
			want:    false,
		},
		"literal command doesnt match when text is different": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(DontCall(t, true, "Test")),
				}
			},
			trigger: "",
			message: "Nope",
			want:    false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmds := tt.commands(t)
			e := commands.NewCommandExecutor(cmds, zap.NewNop())
			got := e.HasMatch(tt.trigger, tt.message)

			assert.Equal(t, tt.want, got)

			for _, m := range cmds {
				if cmd, ok := m.(*mockCmd); ok {
					cmd.AssertExpectations(t)
				}
			}
		})
	}
}

func TestExecutor_Apply(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		commands func(t *testing.T) []commands.Command
		log      *zap.Logger
		message  string
	}{
		"trigger command on matching input": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(MustCall(t, false, "test")),
				}
			},
			message: "!test",
		},
		"trigger command on matching with arguments": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(MustCall(t, false, "test")),
				}
			},
			message: "!test 1 2 3",
		},
		"unmatched command doesn't trigger": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(DontCall(t, false, "test")),
				}
			},
			message: "!boop 1 2 3",
		},
		"literalCommand matches": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(MustCall(t, true, "Test")),
				}
			},
			message: "Test",
		},
		"literalCommand ignores unrelated text": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(DontCall(t, true, "Test")),
				}
			},
			message: "Nope!",
		},

		"literalCommand ignores command on trailing text": {
			commands: func(t *testing.T) []commands.Command {
				t.Helper()

				return []commands.Command{
					asCommand(DontCall(t, true, "Test")),
				}
			},
			message: "Test!",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmds := tt.commands(t)
			e := commands.NewCommandExecutor(cmds, zap.NewNop())
			ctx := t.Context()

			e.Apply(ctx, "!", &discordgo.Message{Content: tt.message, Author: &discordgo.User{ID: "1234"}}, nil)

			for _, m := range cmds {
				if cmd, ok := m.(*mockCmd); ok {
					cmd.AssertExpectations(t)
				}
			}
		})
	}
}
