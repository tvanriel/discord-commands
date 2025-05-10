// Package fx implements FX dependency injection module for commands.
package fx

import (
	commands "github.com/tvanriel/discord-commands"
	"go.uber.org/fx"
)

// GroupCommands is the FX group that commands should be registered under.
const GroupCommands = `group:"commands"`

// Module is the FX module.
//
//nolint:gochecknoglobals // by design
var Module = fx.Module("commands",
	fx.Provide(fx.Annotate(
		commands.NewCommandExecutor,
		fx.ParamTags(GroupCommands),
	)),
)

// AsCommand annotates a command to be a command.
//
//nolint:ireturn // by design.
func AsCommand(in any) any {
	return fx.Annotate(
		in,
		fx.As(new(commands.Command)),
		fx.ResultTags(GroupCommands),
	)
}

// AsCommands annotates many commands to be a command.
func AsCommands(in []any) []any {
	out := make([]any, 0, len(in))

	for _, constructor := range in {
		out = append(out, AsCommand(constructor))
	}

	return out
}
