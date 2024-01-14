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

type wpaSupplicantConfig struct {
	Mac string `koanf:"mac"`
}

var dConfig wpaSupplicantConfig

func main() {
	pmClient := client.NewPMClient("wpaSupplicant")
	pmClient.WaitFor("rootfs")

	setupConfig()
	wpaSupplicant, wpaSupplicantCTX := common.Start("/opt/wpaSupplicant/entrypoint.sh", dConfig.Mac)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		select {
		case sig := <-sigs:
			fmt.Println()
			fmt.Println(sig)
			wpaSupplicant.Cancel()
		case <-wpaSupplicantCTX.Done():
			fmt.Println("wpaSupplicant closed")
		}

		done <- true
	}()

	// app is not running anymore
	if wpaSupplicant.ProcessState != nil {
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

	k.Unmarshal("wifi.wpaSupplicant", &dConfig)

}
