package internal

import (
	"sort"

	discord_structs "github.com/WelcomerTeam/Discord/structs"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// sortGuilds returns Guilds in order of most similar to the query passed.
func sortGuilds(source []*discord_structs.Guild, query string) (sorted []*discord_structs.Guild) {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Name)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted = make([]*discord_structs.Guild, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortChannels returns Channels in order of most similar to the query passed.
func sortChannels(source []*discord_structs.Channel, query string) (sorted []*discord_structs.Channel) {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, channel := range source {
		rank := accurateString(query, channel.Name)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted = make([]*discord_structs.Channel, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortRoles returns Roles in order of most similar to the query passed.
func sortRoles(source []*discord_structs.Role, query string) (sorted []*discord_structs.Role) {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Name)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted = make([]*discord_structs.Role, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortEmojis returns Emojis in order of most similar to the query passed.
func sortEmojis(source []*discord_structs.Emoji, query string) (sorted []*discord_structs.Emoji) {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Name)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted = make([]*discord_structs.Emoji, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortMembers returns Members in order of most similar to the query passed.
func sortMembers(source []*discord_structs.GuildMember, query string) (sorted []*discord_structs.GuildMember) {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Nick, source.User.Username+"#"+source.User.Discriminator, source.User.Username)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted = make([]*discord_structs.GuildMember, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortUsers returns Users in order of most similar to the query passed.
func sortUsers(source []*discord_structs.User, query string) (sorted []*discord_structs.User) {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Username+"#"+source.Discriminator, source.Username)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted = make([]*discord_structs.User, 0, len(source))

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
