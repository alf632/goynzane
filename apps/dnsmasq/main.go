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

type dnsmasqConfig struct {
	IP      string `koanf:"ip"`
	Net     string `koanf:"net"`
	Prefix  string `koanf:"prefix"`
	IPRange string `koanf:"ip-range"`
	DNS     string `koanf:"dns"`
	Domain  string `koanf:"domain"`
	Mac     string `koanf:"mac"`
}

var dConfig dnsmasqConfig

func main() {
	pmClient := client.NewPMClient("dnsmasq")
	pmClient.WaitFor("rootfs")

	setupConfig()
	dnsmasq, dnsmasqCTX := common.Start("/opt/dnsmasq/entrypoint.sh", dConfig.Mac, dConfig.IP, dConfig.Net, dConfig.Prefix)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		select {
		case sig := <-sigs:
			fmt.Println()
			fmt.Println(sig)
			dnsmasq.Cancel()
		case <-dnsmasqCTX.Done():
			fmt.Println("dnsmasq closed")
		}

		done <- true
	}()

	// app is not running anymore
	if dnsmasq.ProcessState != nil {
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

	k.Unmarshal("wifi.dnsmasq", &dConfig)

	//if _, err := os.Stat("/perm/dnsmasq/dnsmasq.conf"); errors.Is(err, os.ErrNotExist) {
	content := []byte(fmt.Sprintf(`
server=%s
no-resolv
listen-address=127.0.0.1,%s
domain-needed
bogus-priv
filterwin2k
domain=%s
local=/%s/
dhcp-range=%s
dhcp-option=option:dns-server,%s`, dConfig.DNS, dConfig.IP, dConfig.Domain, dConfig.Domain, dConfig.IPRange, dConfig.IP))
	os.Mkdir("/perm/dnsmasq", os.ModePerm)
	os.WriteFile("/perm/dnsmasq/dnsmasq.conf", content, 0644)
	//}
}
