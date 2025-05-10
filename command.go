// Package commands implements a command executor for Discordgo.
package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Command is a command to be executed.
type Command interface {
	Name() string
	Apply(ctx *Context) error
	SkipsPrefix() bool
}

// Context is the context that a command is running in.
type Context struct {
	Message *discordgo.Message
	Content string
	Args    []string
	Session *discordgo.Session
}

// Reply replies to a message.
func (ctx *Context) Reply(s string) (*discordgo.Message, error) {
	msg, err := ctx.Session.ChannelMessageSendEmbedReply(
		ctx.Message.ChannelID,

		&discordgo.MessageEmbed{
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Result",
					Value: s,
				},
			},
		},
		ctx.Reference(),
	)
	if err != nil {
		return nil, fmt.Errorf("send reply to command: %w", err)
	}

	return msg, nil
}

// errRed is the red that is used for the error embed sidebar.
const errRed = 0xFF0000

func (ctx *Context) Error(err error) (*discordgo.Message, error) {
	msg, err2 := ctx.Session.ChannelMessageSendEmbedReply(
		ctx.Message.ChannelID,
		&discordgo.MessageEmbed{
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Error",
					Value: err.Error(),
				},
			},
			Color: errRed,
		},
		ctx.Reference(),
	)

	if err2 != nil {
		return nil, fmt.Errorf("send err reply: %w", err2)
	}

	return msg, nil
}

// Reference reports a reference to the message.
func (ctx *Context) Reference() *discordgo.MessageReference {
	return &discordgo.MessageReference{
		MessageID: ctx.Message.ID,
		ChannelID: ctx.Message.ChannelID,
		GuildID:   ctx.Message.GuildID,
	}
}

// DiscordMessageMaxLength is the maximum length a discord message may have.
const DiscordMessageMaxLength = 2000

// ReplyList replies to a message using a list.
func (ctx *Context) ReplyList(s []string) ([]*discordgo.Message, error) {
	if len(s) == 0 {
		return []*discordgo.Message{}, nil
	}

	itemTpl := "`%s`\n"
	templated := make([]string, 0, len(s))

	for i := range s {
		templated = append(templated, fmt.Sprintf(itemTpl, s[i]))
	}

	var (
		sb           strings.Builder
		sentMessages []*discordgo.Message
		err          error
	)

	totalLength := 0
	for i := range templated {
		if totalLength+len(templated[i]) > DiscordMessageMaxLength {
			msg, msgerr := ctx.Session.ChannelMessageSendReply(
				ctx.Message.ChannelID,
				sb.String(),
				ctx.Reference(),
			)
			sentMessages = append(sentMessages, msg)
			err = errors.Join(err, msgerr)

			sb.Reset()

			totalLength = 0
		}

		totalLength += len(templated[i])
		sb.WriteString(templated[i])
	}

	if totalLength > 0 {
		msg, msgerr := ctx.Session.ChannelMessageSendReply(
			ctx.Message.ChannelID,
			sb.String(),
			ctx.Reference(),
		)
		sentMessages = append(sentMessages, msg)
		err = errors.Join(err, msgerr)
	}

	return sentMessages, err
}
