package internal

import discord "github.com/WelcomerTeam/Sandwich-Daemon/discord/structs"

type User discord.User

func NewUser(ctx *EventContext, userID discord.Snowflake) User {
	return User{
		ID: userID,
	}
}
