package main

import (
	"context"
	"os/exec"
	"runtime"
)

func main() {
	ctx := context.Background()
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
