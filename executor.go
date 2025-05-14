package commands

import (
	"context"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// Executor runs commands.
type Executor struct {
	commands []Command
	log      *zap.Logger
}

// NewCommandExecutor builds a new command executor.
func NewCommandExecutor(commands []Command, log *zap.Logger) *Executor {
	return &Executor{
		commands: commands,
		log:      log,
	}
}

// HasMatch determines whether a message has a match on a registered command.
func (e *Executor) HasMatch(trigger string, message string) bool {
	for _, cmd := range e.commands {
		if e.matches(trigger, message, cmd) {
			return true
		}
	}

	return false
}

// Apply finds and executes a command.
func (e *Executor) Apply(ctx context.Context, trigger string, message *discordgo.Message, s *discordgo.Session) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}

	for _, cmd := range e.commands {
		if !e.matches(trigger, message.Content, cmd) {
			continue
		}

		var commandContext *Context

		if cmd.SkipsPrefix() {
			commandContext = &Context{
				Message: message,
				Args:    []string{message.Content},
				Session: s,
				Content: message.Content,
				Ctx:     ctx,
			}
		} else {
			content := StripPrefix(trigger, cmd.Name())(message.Content)
			args := SplitArgs(content)

			commandContext = &Context{
				Message: message,
				Args:    args,
				Session: s,
				Content: content,
				Ctx:     ctx,
			}
		}

		wg.Add(1)

		go func(cmd Command, commandContext *Context) {
			e.applyCommand(cmd, commandContext)
			wg.Done()
		}(cmd, commandContext)
	}

	wg.Wait()
}

func (e *Executor) applyCommand(cmd Command, ctx *Context) {
	e.log.With(messageZapFields(ctx.Message)...).Info("executing command", zap.String("cmd", cmd.Name()))

	err := cmd.Apply(ctx)
	if err != nil {
		e.log.With(messageZapFields(ctx.Message)...).Error(
			"Command failed",
			zap.String("cmd", cmd.Name()),
			zap.NamedError("err", err),
		)

		_, err1 := ctx.Error(err)
		if err1 != nil {
			e.log.With(messageZapFields(ctx.Message)...).Error(
				"Failed to report command reply error to discord",
				zap.NamedError("orig", err),
				zap.NamedError("err", err1),
			)
		}
	}
}

func (e *Executor) matches(trigger string, message string, cmd Command) bool {
	if cmd.SkipsPrefix() {
		if cmd.Name() == message {
			return true
		}
	} else {
		if HasCommandPrefix(trigger, cmd.Name(), message) {
			return true
		}
	}

	return false
}

func messagePermaURL(guild string, channel string, id string) string {
	return strings.Join(
		[]string{
			"https://discord.com/channels/",
			guild,
			"/",
			channel,
			"/",
			id,
		},
		"",
	)
}

func messageZapFields(message *discordgo.Message) []zap.Field {
	return []zap.Field{
		zap.String("username", message.Author.Username),
		zap.String("guild", message.GuildID),
		zap.String("channel", message.ChannelID),
		zap.String("message", message.ID),
		zap.String("content", message.Content),
		zap.String("url", messagePermaURL(message.GuildID, message.ChannelID, message.ID)),
	}
}
