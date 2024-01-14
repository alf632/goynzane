package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alf632/goynzane/apps/common"
)

func main() {
	app, appCTX := common.Start("/usr/sbin/mosquitto", "-c", "/etc/mosquitto/mosquitto-overwrite.conf")

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
			fmt.Println("mosquitto closed")
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
