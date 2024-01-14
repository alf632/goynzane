package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/alf632/goynzane/apps/common"
	"github.com/alf632/goynzane/apps/dependencyManager/client"
)

var hConfig HostapdConf

type HostapdConf struct {
	SSID       string `koanf:"ssid"`
	Passphrase string `koanf:"passphrase"`
	Mac        string `koanf:"mac"`
}

func main() {
	pmClient := client.NewPMClient("hostapd")
	pmClient.WaitFor("rootfs")

	setupConfig()

	hostapd, hostapdCTX := common.Start("/opt/hostapd/entrypoint.sh", hConfig.Mac)
	pmClient.Provide("hostapd")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		select {
		case sig := <-sigs:
			fmt.Println()
			fmt.Println(sig)
			hostapd.Cancel()
		case <-hostapdCTX.Done():
			fmt.Println("hostapd closed")
		}

		done <- true
	}()

	// app is not running anymore
	if hostapd.ProcessState != nil {
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

	k.Unmarshal("wifi.hostapd", &hConfig)

	//if _, err := os.Stat("/perm/hostapd/hostapd.conf"); errors.Is(err, os.ErrNotExist) {
	content := []byte(fmt.Sprintf(`
# "g" simply means 2.4GHz band
hw_mode=g
# the channel to use
channel=10
# limit the frequencies used to those allowed in the country
ieee80211d=1
# the country code
country_code=DE
# 802.11n support
ieee80211n=1
# QoS support, also required for full speed on 802.11n/ac/ax
wmm_enabled=1

# the name of the AP
ssid=%s
# 1=wpa, 2=wep, 3=both
auth_algs=1
# WPA2 only
wpa=2
wpa_key_mgmt=WPA-PSK
rsn_pairwise=CCMP
wpa_passphrase=%s`, hConfig.SSID, hConfig.Passphrase))
	os.Mkdir("/perm/hostapd", os.ModePerm)
	os.WriteFile("/perm/hostapd/hostapd.conf", content, 0644)
	//}
}
