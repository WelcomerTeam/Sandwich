package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/WelcomerTeam/Discord/discord"
	sandwich_protobuf "github.com/WelcomerTeam/Sandwich-Daemon/protobuf"
	sandwich_structs "github.com/WelcomerTeam/Sandwich-Daemon/structs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type GRPCContext struct {
	context.Context

	Logger zerolog.Logger

	SandwichClient sandwich_protobuf.SandwichClient
	GRPCInterface  GRPC

	Session    *discord.Session
	Identifier *sandwich_structs.ManagerConsumerConfiguration
}

type GRPC interface {
	Listen(grpcContext *GRPCContext, identifier string) (client sandwich_protobuf.Sandwich_ListenClient, err error)
	PostAnalytics(grpcContext *GRPCContext, identifier string, data []byte) error

	FetchGuildByID(grpcContext *GRPCContext, guildID discord.Snowflake) (guild *discord.Guild, err error)
	FetchGuildsByName(grpcContext *GRPCContext, query string) (guilds []*discord.Guild, err error)

	FetchChannelByID(grpcContext *GRPCContext, guildID discord.Snowflake, channelID discord.Snowflake) (channel *discord.Channel, err error)
	FetchChannelsByName(grpcContext *GRPCContext, guildID discord.Snowflake, query string) (channels []*discord.Channel, err error)

	FetchRoleByID(grpcContext *GRPCContext, guildID discord.Snowflake, roleID discord.Snowflake) (role *discord.Role, err error)
	FetchRolesByName(grpcContext *GRPCContext, guildID discord.Snowflake, query string) (roles []*discord.Role, err error)

	FetchEmojiByID(grpcContext *GRPCContext, guildID discord.Snowflake, emojiID discord.Snowflake) (emoji *discord.Emoji, err error)
	FetchEmojisByName(grpcContext *GRPCContext, guildID discord.Snowflake, query string) (emojis []*discord.Emoji, err error)

	FetchMemberByID(grpcContext *GRPCContext, guildID discord.Snowflake, memberID discord.Snowflake) (member *discord.GuildMember, err error)
	FetchMembersByID(grpcContext *GRPCContext, guildID discord.Snowflake, memberIDs []discord.Snowflake) (members []*discord.GuildMember, err error)
	FetchMembersByName(grpcContext *GRPCContext, guildID discord.Snowflake, query string) (members []*discord.GuildMember, err error)

	FetchUserByID(grpcContext *GRPCContext, token string, userID discord.Snowflake, createDMChannel bool) (user *discord.User, err error)
	FetchUserByName(grpcContext *GRPCContext, token string, query string, createDMChannel bool) (users []*discord.User, err error)

	FetchConsumerConfiguration(grpcContext *GRPCContext, identifier string) (identifiers *sandwich_structs.SandwichConsumerConfiguration, err error)
	FetchMutualGuilds(grpcContext *GRPCContext, userID discord.Snowflake, expand bool) (guilds []*discord.Guild, err error)

	RequestGuildChunk(grpcContext *GRPCContext, guildID discord.Snowflake) error
	SendWebsocketMessage(grpcContext *GRPCContext, location Location, op int32, data []byte) error
	WhereIsGuild(grpcContext *GRPCContext, guildID discord.Snowflake) (locations []*Location_GuildMember, err error)
}

// Helper structure for SendWebsocketMessage and WhereIsGuild functions.
type Location struct {
	Manager     string
	ShardGroup  int32
	ShardID     int32
}

type Location_GuildMember struct {
	Location
	GuildMember *discord.GuildMember
}

type DefaultGRPCClient struct{}

func NewDefaultGRPCClient() (grpcClient GRPC) {
	grpcClient = &DefaultGRPCClient{}

	return
}

func (grpcClient *DefaultGRPCClient) Listen(grpcContext *GRPCContext, identifier string) (client sandwich_protobuf.Sandwich_ListenClient, err error) {
	client, err = grpcContext.SandwichClient.Listen(grpcContext.Context, &sandwich_protobuf.ListenRequest{
		Identifier: identifier,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	return
}

func (grpcClient *DefaultGRPCClient) PostAnalytics(grpcContext *GRPCContext, identifier string, data []byte) error {
	base, err := grpcContext.SandwichClient.PostAnalytics(grpcContext.Context, &sandwich_protobuf.PostAnalyticsRequest{
		Identifier: identifier,
		Data:       data,
	})
	if err != nil {
		return fmt.Errorf("failed to post analytics: %w", err)
	}

	if !base.Ok {
		return errors.New(base.Error)
	}

	return nil
}

func (grpcClient *DefaultGRPCClient) FetchGuildByID(grpcContext *GRPCContext, guildID discord.Snowflake) (guild *discord.Guild, err error) {
	guildsResponse, err := grpcContext.SandwichClient.FetchGuild(grpcContext.Context, &sandwich_protobuf.FetchGuildRequest{
		GuildIDs: []int64{int64(guildID)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch guilds: %w", err)
	}

	grpcGuild := guildsResponse.Guilds[int64(guildID)]
	if grpcGuild != nil {
		guild, err = sandwich_protobuf.GRPCToGuild(grpcGuild)
		if err != nil {
			return nil, fmt.Errorf("failed to convert protobuf.Guild to Guild: %w", err)
		}
	}

	return guild, nil
}

func (grpcClient *DefaultGRPCClient) FetchGuildsByName(grpcContext *GRPCContext, query string) (guilds []*discord.Guild, err error) {
	guildsResponse, err := grpcContext.SandwichClient.FetchGuild(grpcContext.Context, &sandwich_protobuf.FetchGuildRequest{
		Query: query,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch guilds: %w", err)
	}

	guilds = make([]*discord.Guild, 0, len(guildsResponse.Guilds))

	for _, grpcGuild := range guildsResponse.Guilds {
		guild, err := sandwich_protobuf.GRPCToGuild(grpcGuild)
		if err != nil {
			grpcContext.Logger.Warn().Err(err).Msg("Failed to convert pb.Guild to Guild")

			continue
		}

		guilds = append(guilds, guild)
	}

	return sortGuilds(guilds, query), nil
}

func (grpcClient *DefaultGRPCClient) FetchChannelByID(grpcContext *GRPCContext, guildID discord.Snowflake, channelID discord.Snowflake) (channel *discord.Channel, err error) {
	channelsResponse, err := grpcContext.SandwichClient.FetchGuildChannels(grpcContext.Context, &sandwich_protobuf.FetchGuildChannelsRequest{
		GuildID:    int64(guildID),
		ChannelIDs: []int64{int64(channelID)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channels: %w", err)
	}

	grpcChannel := channelsResponse.GuildChannels[int64(channelID)]
	if grpcChannel != nil {
		channel, err = sandwich_protobuf.GRPCToChannel(grpcChannel)
		if err != nil {
			return nil, fmt.Errorf("failed to convert protobuf.Channel to Channel: %w", err)
		}
	}

	return channel, nil
}

func (grpcClient *DefaultGRPCClient) FetchChannelsByName(grpcContext *GRPCContext, guildID discord.Snowflake, query string) (channels []*discord.Channel, err error) {
	channelsResponse, err := grpcContext.SandwichClient.FetchGuildChannels(grpcContext.Context, &sandwich_protobuf.FetchGuildChannelsRequest{
		GuildID: int64(guildID),
		Query:   query,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channels: %w", err)
	}

	channels = make([]*discord.Channel, 0, len(channelsResponse.GuildChannels))

	for _, grpcChannel := range channelsResponse.GuildChannels {
		channel, err := sandwich_protobuf.GRPCToChannel(grpcChannel)
		if err != nil {
			grpcContext.Logger.Warn().Err(err).Msg("Failed to convert pb.Channel to Channel")

			continue
		}

		channels = append(channels, channel)
	}

	return sortChannels(channels, query), nil
}

func (grpcClient *DefaultGRPCClient) FetchRoleByID(grpcContext *GRPCContext, guildID discord.Snowflake, roleID discord.Snowflake) (role *discord.Role, err error) {
	rolesResponse, err := grpcContext.SandwichClient.FetchGuildRoles(grpcContext.Context, &sandwich_protobuf.FetchGuildRolesRequest{
		GuildID: int64(guildID),
		RoleIDs: []int64{int64(guildID)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}

	grpcRole := rolesResponse.GuildRoles[int64(guildID)]
	if grpcRole != nil {
		role, err = sandwich_protobuf.GRPCToRole(grpcRole)
		if err != nil {
			return nil, fmt.Errorf("failed to convert protobuf.Role to Role: %w", err)
		}
	}

	return role, nil
}

func (grpcClient *DefaultGRPCClient) FetchRolesByName(grpcContext *GRPCContext, guildID discord.Snowflake, query string) (roles []*discord.Role, err error) {
	rolesResponse, err := grpcContext.SandwichClient.FetchGuildRoles(grpcContext.Context, &sandwich_protobuf.FetchGuildRolesRequest{
		GuildID: int64(guildID),
		Query:   query,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}

	roles = make([]*discord.Role, 0, len(rolesResponse.GuildRoles))

	for _, grpcRole := range rolesResponse.GuildRoles {
		role, err := sandwich_protobuf.GRPCToRole(grpcRole)
		if err != nil {
			grpcContext.Logger.Warn().Err(err).Msg("Failed to convert pb.Role to Role")

			continue
		}

		roles = append(roles, role)
	}

	return sortRoles(roles, query), nil
}

func (grpcClient *DefaultGRPCClient) FetchEmojiByID(grpcContext *GRPCContext, guildID discord.Snowflake, emojiID discord.Snowflake) (emoji *discord.Emoji, err error) {
	emojisResponse, err := grpcContext.SandwichClient.FetchGuildEmojis(grpcContext.Context, &sandwich_protobuf.FetchGuildEmojisRequest{
		GuildID:  int64(guildID),
		EmojiIDs: []int64{int64(guildID)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch emojis: %w", err)
	}

	grpcEmoji := emojisResponse.GuildEmojis[int64(guildID)]
	if grpcEmoji != nil {
		emoji, err = sandwich_protobuf.GRPCToEmoji(grpcEmoji)
		if err != nil {
			return nil, fmt.Errorf("failed to convert protobuf.Emoji to Emoji: %w", err)
		}
	}

	return emoji, nil
}

func (grpcClient *DefaultGRPCClient) FetchEmojisByName(grpcContext *GRPCContext, guildID discord.Snowflake, query string) (emojis []*discord.Emoji, err error) {
	emojisResponse, err := grpcContext.SandwichClient.FetchGuildEmojis(grpcContext.Context, &sandwich_protobuf.FetchGuildEmojisRequest{
		GuildID: int64(guildID),
		Query:   query,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch emojis: %w", err)
	}

	emojis = make([]*discord.Emoji, 0, len(emojisResponse.GuildEmojis))

	for _, grpcEmoji := range emojisResponse.GuildEmojis {
		emoji, err := sandwich_protobuf.GRPCToEmoji(grpcEmoji)
		if err != nil {
			grpcContext.Logger.Warn().Err(err).Msg("Failed to convert pb.Emoji to Emoji")

			continue
		}

		emojis = append(emojis, emoji)
	}

	return sortEmojis(emojis, query), nil
}

func (grpcClient *DefaultGRPCClient) FetchMemberByID(grpcContext *GRPCContext, guildID discord.Snowflake, memberID discord.Snowflake) (member *discord.GuildMember, err error) {
	membersResponse, err := grpcContext.SandwichClient.FetchGuildMembers(grpcContext.Context, &sandwich_protobuf.FetchGuildMembersRequest{
		GuildID: int64(guildID),
		UserIDs: []int64{int64(memberID)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch members: %w", err)
	}

	grpcMember := membersResponse.GuildMembers[int64(memberID)]
	if grpcMember != nil {
		member, err = sandwich_protobuf.GRPCToGuildMember(grpcMember)
		if err != nil {
			return nil, fmt.Errorf("failed to convert protobuf.GuildMember to GuildMember: %w", err)
		}
	}

	return member, nil
}

func (grpcClient *DefaultGRPCClient) FetchMembersByID(grpcContext *GRPCContext, guildID discord.Snowflake, memberIDs []discord.Snowflake) (members []*discord.GuildMember, err error) {
	userIDs := make([]int64, 0, len(memberIDs))
	for _, userID := range memberIDs {
		userIDs = append(userIDs, int64(userID))
	}

	membersResponse, err := grpcContext.SandwichClient.FetchGuildMembers(grpcContext.Context, &sandwich_protobuf.FetchGuildMembersRequest{
		GuildID: int64(guildID),
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch members: %w", err)
	}

	members = make([]*discord.GuildMember, 0, len(membersResponse.GuildMembers))

	for _, grpcMember := range membersResponse.GuildMembers {
		member, err := sandwich_protobuf.GRPCToGuildMember(grpcMember)
		if err != nil {
			return nil, fmt.Errorf("failed to convert protobuf.GuildMember to GuildMember: %w", err)
		}

		members = append(members, member)
	}

	return members, nil
}

func (grpcClient *DefaultGRPCClient) FetchMembersByName(grpcContext *GRPCContext, guildID discord.Snowflake, query string) (members []*discord.GuildMember, err error) {
	membersResponse, err := grpcContext.SandwichClient.FetchGuildMembers(grpcContext.Context, &sandwich_protobuf.FetchGuildMembersRequest{
		GuildID: int64(guildID),
		Query:   query,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch members: %w", err)
	}

	members = make([]*discord.GuildMember, 0, len(membersResponse.GuildMembers))

	for _, grpcMember := range membersResponse.GuildMembers {
		member, err := sandwich_protobuf.GRPCToGuildMember(grpcMember)
		if err != nil {
			grpcContext.Logger.Warn().Err(err).Msg("Failed to convert pb.GuildMember to GuildMember")

			continue
		}

		members = append(members, member)
	}

	return sortMembers(members, query), nil
}

func (grpcClient *DefaultGRPCClient) FetchUserByID(grpcContext *GRPCContext, token string, userID discord.Snowflake, createDMChannel bool) (user *discord.User, err error) {
	usersResponse, err := grpcContext.SandwichClient.FetchUsers(grpcContext.Context, &sandwich_protobuf.FetchUsersRequest{
		UserIDs:         []int64{int64(userID)},
		CreateDMChannel: createDMChannel,
		Token:           token,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	grpcUser := usersResponse.Users[int64(userID)]
	if grpcUser != nil {
		user, err = sandwich_protobuf.GRPCToUser(grpcUser)
		if err != nil {
			return nil, fmt.Errorf("failed to convert protobuf.User to User: %w", err)
		}
	}

	return user, nil
}

func (grpcClient *DefaultGRPCClient) FetchUserByName(grpcContext *GRPCContext, token string, query string, createDMChannel bool) (users []*discord.User, err error) {
	usersResponse, err := grpcContext.SandwichClient.FetchUsers(grpcContext.Context, &sandwich_protobuf.FetchUsersRequest{
		Query:           query,
		CreateDMChannel: createDMChannel,
		Token:           token,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	users = make([]*discord.User, 0, len(usersResponse.Users))

	for _, grpcUser := range usersResponse.Users {
		user, err := sandwich_protobuf.GRPCToUser(grpcUser)
		if err != nil {
			grpcContext.Logger.Warn().Err(err).Msg("Failed to convert pb.User to User")

			continue
		}

		users = append(users, user)
	}

	return sortUsers(users, query), nil
}

func (grpcClient *DefaultGRPCClient) FetchConsumerConfiguration(grpcContext *GRPCContext, identifier string) (identifiers *sandwich_structs.SandwichConsumerConfiguration, err error) {
	consumerConfiguration, err := grpcContext.SandwichClient.FetchConsumerConfiguration(grpcContext.Context, &sandwich_protobuf.FetchConsumerConfigurationRequest{
		Identifier: identifier,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch consumer configuration: %w", err)
	}

	identifiers = &sandwich_structs.SandwichConsumerConfiguration{}

	err = json.Unmarshal(consumerConfiguration.File, &identifiers)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal consumer configuration: %w", err)
	}

	return identifiers, nil
}

func (grpcClient *DefaultGRPCClient) FetchMutualGuilds(grpcContext *GRPCContext, userID discord.Snowflake, expand bool) (guilds []*discord.Guild, err error) {
	mutualGuilds, err := grpcContext.SandwichClient.FetchMutualGuilds(grpcContext.Context, &sandwich_protobuf.FetchMutualGuildsRequest{
		UserID: int64(userID),
		Expand: expand,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mutual guilds: %w", err)
	}

	guilds = make([]*discord.Guild, 0, len(mutualGuilds.Guilds))

	for _, grpcGuild := range mutualGuilds.Guilds {
		guild, err := sandwich_protobuf.GRPCToGuild(grpcGuild)
		if err != nil {
			grpcContext.Logger.Warn().Err(err).Msg("Failed to convert pb.Guild to Guild")

			continue
		}

		guilds = append(guilds, guild)
	}

	return guilds, nil
}

func (grpcClient *DefaultGRPCClient) RequestGuildChunk(grpcContext *GRPCContext, guildID discord.Snowflake) error {
	baseResponse, err := grpcContext.SandwichClient.RequestGuildChunk(grpcContext.Context, &sandwich_protobuf.RequestGuildChunkRequest{
		GuildId: int64(guildID),
	})
	if err != nil {
		return fmt.Errorf("failed to request guild chunk: %w", err)
	}

	if baseResponse.Error != "" {
		return errors.New(baseResponse.Error)
	}

	if !baseResponse.Ok {
		return ErrUnknownGRPCError
	}

	return nil
}

func (grpcClient *DefaultGRPCClient) SendWebsocketMessage(grpcContext *GRPCContext, location Location, op int32, data []byte) error {
	baseResponse, err := grpcContext.SandwichClient.SendWebsocketMessage(grpcContext.Context, &sandwich_protobuf.SendWebsocketMessageRequest{
		Manager:       location.Manager,
		ShardGroup:    location.ShardGroup,
		Shard:         location.ShardID,
		GatewayOPCode: op,
		Data:          data,
	})
	if err != nil {
		return fmt.Errorf("failed to send websocket message: %w", err)
	}

	if baseResponse.Error != "" {
		return errors.New(baseResponse.Error)
	}

	if !baseResponse.Ok {
		return ErrUnknownGRPCError
	}

	return nil
}

func (grpcClient *DefaultGRPCClient) WhereIsGuild(grpcContext *GRPCContext, guildID discord.Snowflake) (locations []*Location_GuildMember, err error) {
	locationResponse, err := grpcContext.SandwichClient.WhereIsGuild(grpcContext.Context, &sandwich_protobuf.WhereIsGuildRequest{
		GuildID: int64(guildID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch guild locations: %w", err)
	}

	locations = make([]*Location_GuildMember, 0, len(locationResponse.Locations))

	for _, grpcLocation := range locationResponse.Locations {
		guildMember, _ := sandwich_protobuf.GRPCToGuildMember(grpcLocation.GuildMember)
		locations = append(locations, &Location_GuildMember{
			Location: Location{
				Manager:    grpcLocation.Manager,
				ShardGroup: grpcLocation.ShardGroup,
				ShardID:    grpcLocation.ShardId,
			},
			GuildMember: guildMember,
		})
	}

	return locations, nil
}
