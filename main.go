package main

import (
	"context"

	sandwich "github.com/WelcomerTeam/Sandwich/internal"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.TODO()
	conn, err := grpc.Dial("localhost:15000")
	if err != nil {
		panic(err)
	}

	_ = sandwich.NewSandwich(ctx, conn)
	println("Hello World")
}
