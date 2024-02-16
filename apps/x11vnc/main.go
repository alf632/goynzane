package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alf632/goynzane/apps/common"
	"github.com/gokrazy/gokrazy"
)

func main() {
	gokrazy.DontStartOnBoot()

	app, appCTX := common.StartWithEnv("DISPLAY=:0", "/usr/bin/x11vnc")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		select {
		case sig := <-sigs:
			fmt.Println()
			fmt.Println(sig)
			app.Cancel()
		case <-appCTX.Done():
			fmt.Println("x11vnc closed")
		}

		done <- true
	}()

	// app is not running anymore
	if app.ProcessState != nil {
		return
	}

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
