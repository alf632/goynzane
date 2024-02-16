package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alf632/goynzane/apps/common"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type surfstickConfig struct {
	Mac string `koanf:"mac"`
}

var dConfig surfstickConfig

func main() {
	//pmClient := client.NewPMClient("surfstick")
	//pmClient.WaitFor("rootfs")

	//setupConfig()
	surfstick, surfstickCTX := common.Start("/opt/surfstick/entrypoint.sh", dConfig.Mac)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		select {
		case sig := <-sigs:
			fmt.Println()
			fmt.Println(sig)
			surfstick.Cancel()
		case <-surfstickCTX.Done():
			fmt.Println("surfstick closed")
		}

		done <- true
	}()

	// app is not running anymore
	if surfstick.ProcessState != nil {
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

	k.Unmarshal("wifi.surfstick", &dConfig)

}
