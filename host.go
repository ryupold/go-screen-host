package main

import (
	"context"
	"fmt"

	"github.com/ryupold/grest"
)

func startServer() {
	serverLife := grest.StartListening(context.Background(), "0.0.0.0", 8080, grest.Choose(
		grest.Path("/click").OK(nil),
		grest.ContentType("text/html").OK([]byte(streamPageHTML)),
	))

	select {
	case err, alive := <-serverLife:
		if err != nil {
			panic(err)
		} else if !alive {
			fmt.Println("server stopped")
		}
	}

	fmt.Println("server listening...")
}
