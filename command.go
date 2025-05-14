// Package commands implements a command executor for Discordgo.
package commands

// Command is a command to be executed.
type Command interface {
	Name() string
	Apply(ctx *Context) error
	SkipsPrefix() bool
}

var _ Command = (*cmd)(nil)

type cmd struct {
	literal     bool
	incantation string
	f           func(ctx *Context) error
}

// Apply implements Command.
func (c *cmd) Apply(ctx *Context) error {
	return c.f(ctx)
}

// Name implements Command.
func (c *cmd) Name() string {
	return c.incantation
}

// SkipsPrefix implements Command.
func (c *cmd) SkipsPrefix() bool {
	return c.literal
}

// CommandFunc defines a new command.
//
//nolint:ireturn
func CommandFunc(literal bool, incantation string, f func(*Context) error) Command {
	return &cmd{
		literal:     literal,
		incantation: incantation,
		f:           f,
	}
}

// NewCommand defines a new command.
//
//nolint:ireturn
func NewCommand(incantation string, f func(*Context) error) Command {
	return CommandFunc(false, incantation, f)
}

// NewLiteral defines a new literal command.
//
//nolint:ireturn
func NewLiteral(incantation string, f func(*Context) error) Command {
	return CommandFunc(true, incantation, f)
}
