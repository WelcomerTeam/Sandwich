package internal

import (
	"net/http"

	"github.com/WelcomerTeam/Discord/discord"
	discord_structs "github.com/WelcomerTeam/Discord/structs"
	"golang.org/x/xerrors"
)

// ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) (st []*Message, err error) {
// ChannelMessage(channelID, messageID string) (st *Message, err error) {
// ChannelMessageAck(channelID, messageID, lastToken string) (st *Ack, err error) {
// ChannelMessageSend(channelID string, content string) (*Message, error) {
// ChannelMessageSendComplex(channelID string, data *MessageSend) (st *Message, err error) {
// ChannelMessageSendTTS(channelID string, content string) (*Message, error) {
// ChannelMessageSendEmbed(channelID string, embed *MessageEmbed) (*Message, error) {
// ChannelMessageSendEmbeds(channelID string, embeds []*MessageEmbed) (*Message, error) {
// ChannelMessageSendReply(channelID string, content string, reference *MessageReference) (*Message, error) {
// ChannelMessageEdit(channelID, messageID, content string) (*Message, error) {
// ChannelMessageEditComplex(m *MessageEdit) (st *Message, err error) {
// ChannelMessageEditEmbed(channelID, messageID string, embed *MessageEmbed) (*Message, error) {
// ChannelMessageEditEmbeds(channelID, messageID string, embeds []*MessageEmbed) (*Message, error) {
// ChannelMessageDelete(channelID, messageID string) (err error) {
// ChannelMessagesBulkDelete(channelID string, messages []string) (err error) {
// ChannelMessagePin(channelID, messageID string) (err error) {
// ChannelMessageUnpin(channelID, messageID string) (err error) {
// ChannelMessagesPinned(channelID string) (st []*Message, err error) {
// ChannelFileSend(channelID, name string, r io.Reader) (*Message, error) {
// ChannelFileSendWithMessage(channelID, content string, name string, r io.Reader) (*Message, error) {

func ChannelMessageSend(c *Channel, ctx *EventContext, data *discord_structs.MessageParams, files []*discord_structs.File) (message *Message, err error) {
	url := discord.EndpointChannelMessages(c.ID.String())

	if len(data.Files) > 0 {
		var contentType string

		var body []byte

		contentType, body, err = multipartBodyWithJSON(data, data.Files)
		if err != nil {
			return nil, xerrors.Errorf("Failed to create file body: %v", err)
		}

		err = ctx.HTTPSession.FetchBJBot(ctx, http.MethodPost, url, contentType, body, &message)
		println("A", err)
	} else {
		err = ctx.HTTPSession.FetchJJBot(ctx, http.MethodPost, url, data, &message)
		println("B", err)
	}

	if err != nil {
		return nil, xerrors.Errorf("Failed to send message: %v", err)
	}

	return message, nil
}
