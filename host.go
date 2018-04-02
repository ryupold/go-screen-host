package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/ryupold/grest"
)

func startWebServer(ctx context.Context, port uint16) {
	fmt.Println("starting server...")
	serverLife := grest.StartListening(ctx, "0.0.0.0", port,
		grest.Choose(
			grest.TypedPath("/click/%d/%d/%s", func(u grest.WebUnit, params []interface{}) *grest.WebUnit {
				x := params[0].(int)
				y := params[1].(int)
				state := params[1].(string)

				fmt.Printf("mouse: (%d, %d) -> %s\n", x, y, state)

				return &u
			}).OK(nil),
			grest.ContentType("text/html").OK([]byte(strings.Replace(streamPageHTML, "{{appName}}", appName, -1))),
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
