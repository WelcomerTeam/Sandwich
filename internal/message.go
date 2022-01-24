package internal

import (
	discord_structs "github.com/WelcomerTeam/Discord/structs"
)

type Message discord_structs.Message

type MessageChannelMention discord_structs.MessageChannelMention

type MessageReference discord_structs.MessageReference

type MessageReaction discord_structs.MessageReaction

type MessageAllowedMentions discord_structs.MessageAllowedMentions

type MessageAttachment discord_structs.MessageAttachment

type MessageActivity discord_structs.MessageActivity

func (c *Channel) Send(ctx *EventContext, data *discord_structs.MessageParams, files []*discord_structs.File) (message *Message, err error) {
	return ChannelMessageSend(c, ctx, data, files)
}
