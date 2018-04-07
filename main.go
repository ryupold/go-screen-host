package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"runtime"

	"github.com/getlantern/systray"
	_ "github.com/qodrorid/godaemon"
	"github.com/sqweek/dialog"
)

//go:generate go run internal/resources.go

const (
	appName = "GoScreen Host"
	//Version of the application
	Version       = "1.0.0"
	wwwPort       = 8080
	dataPort      = 4545
	streamingPort = 56565
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	go startWebServer(ctx, wwwPort)

	ip, err := externalIP()
	if err != nil {
		systray.SetTooltip(err.Error())
	}
	open(fmt.Sprintf("http://localhost:8080/%s:%d", ip, streamingPort))
	go func() {
		if err := redirectJPEGs(ctx, dataPort, streamingPort); err != nil {
			systray.SetTooltip(err.Error())
		}
	}()

	systray.Run(onReady, cancel)
}

func onReady() {
	if runtime.GOOS == "darwin" {
		systray.SetTitle("")
	} else {
		systray.SetTitle(appName)
	}

	systray.SetIcon(binICOIcon)

	//menu items
	showStreamMenuItem := systray.AddMenuItem("Show Stream", "Open browser window to show the stream")
	systray.AddSeparator()
	aboutMenuItem := systray.AddMenuItem("About", "")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the host")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			case <-showStreamMenuItem.ClickedCh:
				ip, err := externalIP()
				if err != nil {
					systray.SetTooltip(err.Error())
				}
				open(fmt.Sprintf("http://localhost:8080/%s:%d", ip, streamingPort))
			case <-aboutMenuItem.ClickedCh:
				log("about clicked")
				dialog.Message("%s %s", appName, Version).Title(appName).Info()
			}
		}
	}()
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

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
