package internal

import (
	"sort"

	discord "github.com/WelcomerTeam/Discord/discord"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// sortGuilds returns Guilds in order of most similar to the query passed.
func sortGuilds(source []*discord.Guild, query string) []*discord.Guild {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Name)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted := make([]*discord.Guild, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortChannels returns Channels in order of most similar to the query passed.
func sortChannels(source []*discord.Channel, query string) []*discord.Channel {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, channel := range source {
		rank := accurateString(query, channel.Name)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted := make([]*discord.Channel, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortRoles returns Roles in order of most similar to the query passed.
func sortRoles(source []*discord.Role, query string) []*discord.Role {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Name)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted := make([]*discord.Role, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortEmojis returns Emojis in order of most similar to the query passed.
func sortEmojis(source []*discord.Emoji, query string) []*discord.Emoji {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Name)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted := make([]*discord.Emoji, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortMembers returns Members in order of most similar to the query passed.
func sortMembers(source []*discord.GuildMember, query string) []*discord.GuildMember {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Nick, source.User.Username+"#"+source.User.Discriminator, source.User.Username)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted := make([]*discord.GuildMember, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortUsers returns Users in order of most similar to the query passed.
func sortUsers(source []*discord.User, query string) []*discord.User {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Username+"#"+source.Discriminator, source.Username)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted := make([]*discord.User, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// accurateString returns most accurate string closest to query.
func accurateString(query string, source ...string) fuzzy.Rank {
	ranks := fuzzy.RankFindNormalizedFold(query, source)

	if len(ranks) > 0 {
		return ranks[0]
	}

	return fuzzy.Rank{}
}
