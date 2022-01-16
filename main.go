package main

import (
	"context"

	"github.com/WelcomerTeam/Sandwich-Daemon/structs"
	sandwich "github.com/WelcomerTeam/Sandwich/internal"
	messaging "github.com/WelcomerTeam/Sandwich/messaging"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.TODO()
	conn, err := grpc.Dial("localhost:15000", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	sandwichClient := sandwich.NewSandwich(ctx, conn)

	welcomerBot := sandwich.NewBot()

	welcomerBot.RegisterOnMessageCreateEvent(func(ctx *sandwich.Context, message sandwich.Message) (err error) {
		println(message.Author.Username, "#", message.Author.Discriminator, message.Content)

		return
	})

	welcomerBot.RegisterOnGuildMemberUpdateEvent(func(ctx *sandwich.Context, before, after sandwich.GuildMember) (err error) {
		println(before.User, after.User)

		return
	})

	sandwichClient.RegisterBot("welcomer", welcomerBot)

	mqC := messaging.NewStanMQClient()
	err = mqC.Connect(ctx, "sdc", map[string]interface{}{
		"Address": "localhost",
		"Cluster": "cluster",
		"Channel": "channel",
	})

	if err != nil {
		panic(err)
	}

	err = mqC.Subscribe(ctx, "sandwich")

	if err != nil {
		panic(err)
	}

	c := mqC.Chan()

	var p structs.SandwichPayload

	for {
		select {
		case m := <-c:
			if err := jsoniter.Unmarshal(m, &p); err == nil {
				err = sandwichClient.DispatchSandwichPayload(p)
				if err != nil {
					println(err.Error())
					println(string(m))
				}
			} else {
				println(err.Error())
				println(string(m))
			}
		}
	}
}
