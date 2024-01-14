package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alf632/goynzane/apps/common"
	"github.com/alf632/goynzane/apps/dependencyManager/client"
)

func main() {
	pmClient := client.NewPMClient("chromium")
	pmClient.WaitFor("rootfs")
	waitForWebsite("http://localhost:8123/lovelace-mqtt/0")

	common.RunWithEnv("DISPLAY=:0", "/usr/bin/xset", "-dpms", "s", "off", "s", "noblank", "s", "0", "0", "s", "noexpose")
	chr, chrCTX := common.StartWithEnv("DISPLAY=:0", "/opt/chromium/entrypoint.sh", "--no-sandbox", "--kiosk", "-a", "http://localhost:8123/lovelace-mqtt/0")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		select {
		case sig := <-sigs:
			fmt.Println()
			fmt.Println(sig)
			chr.Cancel()
		case <-chrCTX.Done():
			fmt.Println("chromium closed")
		}

		done <- true
	}()

	// chrome is not running anymore
	if chr.ProcessState != nil {
		return
	}

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

}

func waitForWebsite(url string) {
	done := false
	for !done {
		time.Sleep(time.Second * 5)
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			continue
		}

		if resp.StatusCode == 200 {
			done = true
		} else {
			fmt.Println(http.StatusText(resp.StatusCode), "waiting for website")
		}
	}
}
