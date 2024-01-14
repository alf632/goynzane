package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alf632/goynzane/apps/common"
	"github.com/alf632/goynzane/apps/dependencyManager/client"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type iwdConfig struct {
	Mac string `koanf:"mac"`
}

var dConfig iwdConfig

func main() {
	pmClient := client.NewPMClient("iwd")
	pmClient.WaitFor("rootfs")

	setupConfig()
	iwd, iwdCTX := common.Start("/opt/iwd/entrypoint.sh", dConfig.Mac)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		select {
		case sig := <-sigs:
			fmt.Println()
			fmt.Println(sig)
			iwd.Cancel()
		case <-iwdCTX.Done():
			fmt.Println("iwd closed")
		}

		done <- true
	}()

	// app is not running anymore
	if iwd.ProcessState != nil {
		return
	}

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

}

func setupConfig() {
	var k = koanf.New(".")
	if err := k.Load(file.Provider("/etc/goynzane.yaml"), yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	k.Unmarshal("wifi.iwd", &dConfig)

}
