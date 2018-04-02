package main

import (
	"context"
	"os/exec"
	"runtime"
)

const (
	appName = "GoScreen"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startWebServer(ctx, 8080)

	open("http://localhost:8080")

	if err := redirectJPEGs(ctx, 4545, 56565); err != nil {
		panic(err)
	}

}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
