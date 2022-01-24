package internal

import (
	"sort"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// sortGuilds returns Guilds in order of most similar to the query passed.
func sortGuilds(source []*Guild, query string) (sorted []*Guild) {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Name)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted = make([]*Guild, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortChannels returns Channels in order of most similar to the query passed.
func sortChannels(source []*Channel, query string) (sorted []*Channel) {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, channel := range source {
		rank := accurateString(query, channel.Name)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted = make([]*Channel, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortRoles returns Roles in order of most similar to the query passed.
func sortRoles(source []*Role, query string) (sorted []*Role) {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Name)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted = make([]*Role, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortEmojis returns Emojis in order of most similar to the query passed.
func sortEmojis(source []*Emoji, query string) (sorted []*Emoji) {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Name)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted = make([]*Emoji, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortMembers returns Members in order of most similar to the query passed.
func sortMembers(source []*GuildMember, query string) (sorted []*GuildMember) {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Nick, source.User.Username+"#"+source.User.Discriminator, source.User.Username)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted = make([]*GuildMember, 0, len(source))

	sort.Sort(ranks)

	for _, rank := range ranks {
		sorted = append(sorted, source[rank.OriginalIndex])
	}

	return sorted
}

// sortUsers returns Users in order of most similar to the query passed.
func sortUsers(source []*User, query string) (sorted []*User) {
	if len(source) < 2 {
		return source
	}

	var ranks fuzzy.Ranks

	for index, source := range source {
		rank := accurateString(query, source.Username+"#"+source.Discriminator, source.Username)
		rank.OriginalIndex = index
		ranks = append(ranks, rank)
	}

	sorted = make([]*User, 0, len(source))

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
