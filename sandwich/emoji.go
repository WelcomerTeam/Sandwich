package internal

import (
	"fmt"

	discord "github.com/WelcomerTeam/Discord/discord"
	sandwich_daemon "github.com/WelcomerTeam/Sandwich-Daemon"
	sandwich_protobuf "github.com/WelcomerTeam/Sandwich-Daemon/proto"
)

func NewEmoji(guildID *discord.Snowflake, emojiID discord.Snowflake) *discord.Emoji {
	return &discord.Emoji{
		ID:      emojiID,
		GuildID: guildID,
	}
}

func FetchEmoji(ctx *GRPCContext, emoji *discord.Emoji) (*discord.Emoji, error) {
	if emoji.Name != "" {
		return emoji, nil
	}

	if emoji.GuildID == nil {
		return emoji, ErrFetchMissingGuild
	}

	gEmojis, err := ctx.SandwichClient.FetchGuildEmoji(ctx, &sandwich_protobuf.FetchGuildEmojiRequest{
		GuildId:  int64(*emoji.GuildID),
		EmojiIds: []int64{int64(emoji.ID)},
	})
	if err != nil {
		return emoji, fmt.Errorf("failed to fetch emoji: %w", err)
	}

	gEmoji, ok := gEmojis.GetEmojis()[int64(emoji.ID)]
	if !ok {
		return nil, ErrEmojiNotFound
	}

	emoji = sandwich_daemon.PBToEmoji(gEmoji)

	if emoji.ID.IsNil() {
		return emoji, ErrEmojiNotFound
	}

	return emoji, nil
}
