package internal

import (
	"encoding/json"

	"github.com/WelcomerTeam/Discord/discord"
	sandwich_protobuf "github.com/WelcomerTeam/Sandwich-Daemon/protobuf"
	sandwich_structs "github.com/WelcomerTeam/Sandwich-Daemon/structs"
	"golang.org/x/xerrors"
)

type GRPC interface {
	Listen(eventCtx *EventContext, identifier string) (client sandwich_protobuf.Sandwich_ListenClient, err error)
	PostAnalytics(eventCtx *EventContext, identifier string, data []byte) (err error)

	FetchGuildByID(eventCtx *EventContext, guildID discord.Snowflake) (guild *Guild, err error)
	FetchGuildsByName(eventCtx *EventContext, query string) (guilds []*Guild, err error)

	FetchChannelByID(eventCtx *EventContext, guildID discord.Snowflake, channelID discord.Snowflake) (channel *Channel, err error)
	FetchChannelsByName(eventCtx *EventContext, guildID discord.Snowflake, query string) (channels []*Channel, err error)

	FetchRoleByID(eventCtx *EventContext, guildID discord.Snowflake, roleID discord.Snowflake) (role *Role, err error)
	FetchRolesByName(eventCtx *EventContext, guildID discord.Snowflake, query string) (roles []*Role, err error)

	FetchEmojiByID(eventCtx *EventContext, guildID discord.Snowflake, emojiID discord.Snowflake) (emoji *Emoji, err error)
	FetchEmojisByName(eventCtx *EventContext, guildID discord.Snowflake, query string) (emojis []*Emoji, err error)

	FetchMemberByID(eventCtx *EventContext, guildID discord.Snowflake, memberID discord.Snowflake) (member *GuildMember, err error)
	FetchMembersByName(eventCtx *EventContext, guildID discord.Snowflake, query string) (members []*GuildMember, err error)

	FetchUserByID(eventCtx *EventContext, userID discord.Snowflake, createDMChannel bool) (user *User, err error)
	FetchUserByName(eventCtx *EventContext, query string, createDMChannel bool) (users []*User, err error)

	FetchConsumerConfiguration(eventCtx *EventContext, identifier string) (identifiers *sandwich_structs.SandwichConsumerConfiguration, err error)
	FetchMutualGuilds(eventCtx *EventContext, userID discord.Snowflake, expand bool) (guilds []*Guild, err error)

	RequestGuildChunk(eventCtx *EventContext, guildID discord.Snowflake) (err error)
	SendWebsocketMessage(eventCtx *EventContext, location Location, op int32, data []byte) (err error)
	WhereIsGuild(eventCtx *EventContext, guildID discord.Snowflake) (locations []*Location, err error)
}

// Helper structure for SendWebsocketMessage and WhereIsGuild functions.
type Location struct {
	Manager    string
	ShardGroup int32
	ShardID    int32
}

type DefaultGRPCClient struct{}

func NewDefaultGRPCClient() (grpcClient GRPC) {
	grpcClient = &DefaultGRPCClient{}

	return
}

func (grpcClient *DefaultGRPCClient) Listen(eventCtx *EventContext, identifier string) (client sandwich_protobuf.Sandwich_ListenClient, err error) {
	client, err = eventCtx.Sandwich.sandwichClient.Listen(eventCtx.Context, &sandwich_protobuf.ListenRequest{
		Identifier: identifier,
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to listen: %v", err)
	}

	return
}

func (grpcClient *DefaultGRPCClient) PostAnalytics(eventCtx *EventContext, identifier string, data []byte) (err error) {
	base, err := eventCtx.Sandwich.sandwichClient.PostAnalytics(eventCtx.Context, &sandwich_protobuf.PostAnalyticsRequest{
		Identifier: identifier,
		Data:       data,
	})
	if err != nil {
		return xerrors.Errorf("Failed to post analytics: %v", err)
	}

	if !base.Ok {
		return xerrors.New(base.Error)
	}

	return nil
}

func (grpcClient *DefaultGRPCClient) FetchGuildByID(eventCtx *EventContext, guildID discord.Snowflake) (guild *Guild, err error) {
	guildsResponse, err := eventCtx.Sandwich.sandwichClient.FetchGuild(eventCtx.Context, &sandwich_protobuf.FetchGuildRequest{
		GuildIDs: []int64{int64(guildID)},
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch guilds: %v", err)
	}

	grpcGuild := guildsResponse.Guilds[int64(guildID)]
	if grpcGuild != nil {
		sandwichGuild, err := sandwich_protobuf.GRPCToGuild(grpcGuild)
		if err != nil {
			return nil, xerrors.Errorf("Failed to convert sandwich_protobuf.Guild to Guild: %v", err)
		}

		guild = (*Guild)(sandwichGuild)
	}

	return guild, nil
}

func (grpcClient *DefaultGRPCClient) FetchGuildsByName(eventCtx *EventContext, query string) (guilds []*Guild, err error) {
	guildsResponse, err := eventCtx.Sandwich.sandwichClient.FetchGuild(eventCtx.Context, &sandwich_protobuf.FetchGuildRequest{
		Query: query,
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch guilds: %v", err)
	}

	guilds = make([]*Guild, 0, len(guildsResponse.Guilds))

	for _, grpcGuild := range guildsResponse.Guilds {
		sandwichGuild, err := sandwich_protobuf.GRPCToGuild(grpcGuild)
		if err != nil {
			eventCtx.Logger.Warn().Err(err).Msg("Failed to convert pb.Guild to Guild")

			continue
		}

		guilds = append(guilds, (*Guild)(sandwichGuild))
	}

	return sortGuilds(guilds, query), nil
}

func (grpcClient *DefaultGRPCClient) FetchChannelByID(eventCtx *EventContext, guildID discord.Snowflake, channelID discord.Snowflake) (channel *Channel, err error) {
	channelsResponse, err := eventCtx.Sandwich.sandwichClient.FetchGuildChannels(eventCtx.Context, &sandwich_protobuf.FetchGuildChannelsRequest{
		GuildID:    int64(guildID),
		ChannelIDs: []int64{int64(channelID)},
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch channels: %v", err)
	}

	grpcChannel := channelsResponse.GuildChannels[int64(channelID)]
	if grpcChannel != nil {
		sandwichChannel, err := sandwich_protobuf.GRPCToChannel(grpcChannel)
		if err != nil {
			return nil, xerrors.Errorf("Failed to convert sandwich_protobuf.Channel to Channel: %v", err)
		}

		channel = (*Channel)(sandwichChannel)
	}

	return channel, nil
}

func (grpcClient *DefaultGRPCClient) FetchChannelsByName(eventCtx *EventContext, guildID discord.Snowflake, query string) (channels []*Channel, err error) {
	channelsResponse, err := eventCtx.Sandwich.sandwichClient.FetchGuildChannels(eventCtx.Context, &sandwich_protobuf.FetchGuildChannelsRequest{
		GuildID: int64(guildID),
		Query:   query,
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch channels: %v", err)
	}

	channels = make([]*Channel, 0, len(channelsResponse.GuildChannels))

	for _, grpcChannel := range channelsResponse.GuildChannels {
		sandwichChannel, err := sandwich_protobuf.GRPCToChannel(grpcChannel)
		if err != nil {
			eventCtx.Logger.Warn().Err(err).Msg("Failed to convert pb.Channel to Channel")

			continue
		}

		channels = append(channels, (*Channel)(sandwichChannel))
	}

	return sortChannels(channels, query), nil
}

func (grpcClient *DefaultGRPCClient) FetchRoleByID(eventCtx *EventContext, guildID discord.Snowflake, roleID discord.Snowflake) (role *Role, err error) {
	rolesResponse, err := eventCtx.Sandwich.sandwichClient.FetchGuildRoles(eventCtx.Context, &sandwich_protobuf.FetchGuildRolesRequest{
		GuildID: int64(guildID),
		RoleIDs: []int64{int64(guildID)},
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch roles: %v", err)
	}

	grpcRole := rolesResponse.GuildRoles[int64(guildID)]
	if grpcRole != nil {
		sandwichRole, err := sandwich_protobuf.GRPCToRole(grpcRole)
		if err != nil {
			return nil, xerrors.Errorf("Failed to convert sandwich_protobuf.Role to Role: %v", err)
		}

		role = (*Role)(sandwichRole)
	}

	return role, nil
}

func (grpcClient *DefaultGRPCClient) FetchRolesByName(eventCtx *EventContext, guildID discord.Snowflake, query string) (roles []*Role, err error) {
	rolesResponse, err := eventCtx.Sandwich.sandwichClient.FetchGuildRoles(eventCtx.Context, &sandwich_protobuf.FetchGuildRolesRequest{
		GuildID: int64(guildID),
		Query:   query,
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch roles: %v", err)
	}

	roles = make([]*Role, 0, len(rolesResponse.GuildRoles))

	for _, grpcRole := range rolesResponse.GuildRoles {
		sandwichRole, err := sandwich_protobuf.GRPCToRole(grpcRole)
		if err != nil {
			eventCtx.Logger.Warn().Err(err).Msg("Failed to convert pb.Role to Role")

			continue
		}

		roles = append(roles, (*Role)(sandwichRole))
	}

	return sortRoles(roles, query), nil
}

func (grpcClient *DefaultGRPCClient) FetchEmojiByID(eventCtx *EventContext, guildID discord.Snowflake, emojiID discord.Snowflake) (emoji *Emoji, err error) {
	emojisResponse, err := eventCtx.Sandwich.sandwichClient.FetchGuildEmojis(eventCtx.Context, &sandwich_protobuf.FetchGuildEmojisRequest{
		GuildID:  int64(guildID),
		EmojiIDs: []int64{int64(guildID)},
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch emojis: %v", err)
	}

	grpcEmoji := emojisResponse.GuildEmojis[int64(guildID)]
	if grpcEmoji != nil {
		sandwichEmoji, err := sandwich_protobuf.GRPCToEmoji(grpcEmoji)
		if err != nil {
			return nil, xerrors.Errorf("Failed to convert sandwich_protobuf.Emoji to Emoji: %v", err)
		}

		emoji = (*Emoji)(sandwichEmoji)
	}

	return emoji, nil
}

func (grpcClient *DefaultGRPCClient) FetchEmojisByName(eventCtx *EventContext, guildID discord.Snowflake, query string) (emojis []*Emoji, err error) {
	emojisResponse, err := eventCtx.Sandwich.sandwichClient.FetchGuildEmojis(eventCtx.Context, &sandwich_protobuf.FetchGuildEmojisRequest{
		GuildID: int64(guildID),
		Query:   query,
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch emojis: %v", err)
	}

	emojis = make([]*Emoji, 0, len(emojisResponse.GuildEmojis))

	for _, grpcEmoji := range emojisResponse.GuildEmojis {
		sandwichEmoji, err := sandwich_protobuf.GRPCToEmoji(grpcEmoji)
		if err != nil {
			eventCtx.Logger.Warn().Err(err).Msg("Failed to convert pb.Emoji to Emoji")

			continue
		}

		emojis = append(emojis, (*Emoji)(sandwichEmoji))
	}

	return sortEmojis(emojis, query), nil
}

func (grpcClient *DefaultGRPCClient) FetchMemberByID(eventCtx *EventContext, guildID discord.Snowflake, memberID discord.Snowflake) (member *GuildMember, err error) {
	membersResponse, err := eventCtx.Sandwich.sandwichClient.FetchGuildMembers(eventCtx.Context, &sandwich_protobuf.FetchGuildMembersRequest{
		GuildID: int64(guildID),
		UserIDs: []int64{int64(memberID)},
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch members: %v", err)
	}

	grpcMember := membersResponse.GuildMembers[int64(memberID)]
	if grpcMember != nil {
		sandwichMember, err := sandwich_protobuf.GRPCToGuildMember(grpcMember)
		if err != nil {
			return nil, xerrors.Errorf("Failed to convert sandwich_protobuf.GuildMember to GuildMember: %v", err)
		}

		member = (*GuildMember)(sandwichMember)
	}

	return member, nil
}

func (grpcClient *DefaultGRPCClient) FetchMembersByName(eventCtx *EventContext, guildID discord.Snowflake, query string) (members []*GuildMember, err error) {
	membersResponse, err := eventCtx.Sandwich.sandwichClient.FetchGuildMembers(eventCtx.Context, &sandwich_protobuf.FetchGuildMembersRequest{
		GuildID: int64(guildID),
		Query:   query,
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch members: %v", err)
	}

	members = make([]*GuildMember, 0, len(membersResponse.GuildMembers))

	for _, grpcMember := range membersResponse.GuildMembers {
		sandwichMember, err := sandwich_protobuf.GRPCToGuildMember(grpcMember)
		if err != nil {
			eventCtx.Logger.Warn().Err(err).Msg("Failed to convert pb.GuildMember to GuildMember")

			continue
		}

		members = append(members, (*GuildMember)(sandwichMember))
	}

	return sortMembers(members, query), nil
}

func (grpcClient *DefaultGRPCClient) FetchUserByID(eventCtx *EventContext, userID discord.Snowflake, createDMChannel bool) (user *User, err error) {
	usersResponse, err := eventCtx.Sandwich.sandwichClient.FetchUsers(eventCtx.Context, &sandwich_protobuf.FetchUsersRequest{
		UserIDs:         []int64{int64(userID)},
		CreateDMChannel: createDMChannel,
		Token:           eventCtx.Identifier.Token,
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch users: %v", err)
	}

	grpcUser := usersResponse.Users[int64(userID)]
	if grpcUser != nil {
		sandwichUser, err := sandwich_protobuf.GRPCToUser(grpcUser)
		if err != nil {
			return nil, xerrors.Errorf("Failed to convert sandwich_protobuf.User to User: %v", err)
		}

		user = (*User)(sandwichUser)
	}

	return user, nil
}

func (grpcClient *DefaultGRPCClient) FetchUserByName(eventCtx *EventContext, query string, createDMChannel bool) (users []*User, err error) {
	usersResponse, err := eventCtx.Sandwich.sandwichClient.FetchUsers(eventCtx.Context, &sandwich_protobuf.FetchUsersRequest{
		Query:           query,
		CreateDMChannel: createDMChannel,
		Token:           eventCtx.Identifier.Token,
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch users: %v", err)
	}

	users = make([]*User, 0, len(usersResponse.Users))

	for _, grpcUser := range usersResponse.Users {
		sandwichUser, err := sandwich_protobuf.GRPCToUser(grpcUser)
		if err != nil {
			eventCtx.Logger.Warn().Err(err).Msg("Failed to convert pb.User to User")

			continue
		}

		users = append(users, (*User)(sandwichUser))
	}

	return sortUsers(users, query), nil
}

func (grpcClient *DefaultGRPCClient) FetchConsumerConfiguration(eventCtx *EventContext, identifier string) (identifiers *sandwich_structs.SandwichConsumerConfiguration, err error) {
	consumerConfiguration, err := eventCtx.Sandwich.sandwichClient.FetchConsumerConfiguration(eventCtx.Context, &sandwich_protobuf.FetchConsumerConfigurationRequest{
		Identifier: identifier,
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch consumer configuration: %v", err)
	}

	identifiers = &sandwich_structs.SandwichConsumerConfiguration{}

	err = json.Unmarshal(consumerConfiguration.File, &identifiers)
	if err != nil {
		return nil, xerrors.Errorf("Failed to unmarshal consumer configuration: %v", err)
	}

	return identifiers, nil
}

func (grpcClient *DefaultGRPCClient) FetchMutualGuilds(eventCtx *EventContext, userID discord.Snowflake, expand bool) (guilds []*Guild, err error) {
	mutualGuilds, err := eventCtx.Sandwich.sandwichClient.FetchMutualGuilds(eventCtx.Context, &sandwich_protobuf.FetchMutualGuildsRequest{
		UserID: int64(userID),
		Expand: expand,
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch mutual guilds: %v", err)
	}

	guilds = make([]*Guild, 0, len(mutualGuilds.Guilds))

	for _, grpcGuild := range mutualGuilds.Guilds {
		sandwichGuild, err := sandwich_protobuf.GRPCToGuild(grpcGuild)
		if err != nil {
			eventCtx.Logger.Warn().Err(err).Msg("Failed to convert pb.Guild to Guild")

			continue
		}

		guilds = append(guilds, (*Guild)(sandwichGuild))
	}

	return guilds, nil
}

func (grpcClient *DefaultGRPCClient) RequestGuildChunk(eventCtx *EventContext, guildID discord.Snowflake) (err error) {
	baseResponse, err := eventCtx.Sandwich.sandwichClient.RequestGuildChunk(eventCtx.Context, &sandwich_protobuf.RequestGuildChunkRequest{
		GuildId: int64(guildID),
	})
	if err != nil {
		return xerrors.Errorf("Failed to request guild chunk: %v", err)
	}

	if baseResponse.Error != "" {
		return xerrors.New(baseResponse.Error)
	}

	if !baseResponse.Ok {
		return ErrUnknownGRPCError
	}

	return nil
}

func (grpcClient *DefaultGRPCClient) SendWebsocketMessage(eventCtx *EventContext, location Location, op int32, data []byte) (err error) {
	baseResponse, err := eventCtx.Sandwich.sandwichClient.SendWebsocketMessage(eventCtx.Context, &sandwich_protobuf.SendWebsocketMessageRequest{
		Manager:       location.Manager,
		ShardGroup:    location.ShardGroup,
		Shard:         location.ShardID,
		GatewayOPCode: op,
		Data:          data,
	})
	if err != nil {
		return xerrors.Errorf("Failed to send websocket message: %v", err)
	}

	if baseResponse.Error != "" {
		return xerrors.New(baseResponse.Error)
	}

	if !baseResponse.Ok {
		return ErrUnknownGRPCError
	}

	return nil
}

func (grpcClient *DefaultGRPCClient) WhereIsGuild(eventCtx *EventContext, guildID discord.Snowflake) (locations []*Location, err error) {
	locationResponse, err := eventCtx.Sandwich.sandwichClient.WhereIsGuild(eventCtx.Context, &sandwich_protobuf.WhereIsGuildRequest{
		GuildID: int64(guildID),
	})
	if err != nil {
		return nil, xerrors.Errorf("Failed to fetch guild locations: %v", err)
	}

	locations = make([]*Location, 0, len(locationResponse.Locations))

	for _, grpcLocation := range locationResponse.Locations {
		locations = append(locations, &Location{
			Manager:    grpcLocation.Manager,
			ShardGroup: grpcLocation.ShardGroup,
			ShardID:    grpcLocation.ShardId,
		})
	}

	return locations, nil
}
