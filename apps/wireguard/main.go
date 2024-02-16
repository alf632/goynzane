package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/alf632/goynzane/apps/common"
	"github.com/alf632/goynzane/apps/dependencyManager/client"
	depman "github.com/alf632/goynzane/apps/dependencyManager/lib"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type WireguardInstance struct {
	WireguardConfig

	ready bool
	cmd   *exec.Cmd
	ctx   context.Context
}

type WireguardPeer struct {
	Name       string `koanf:"name"`
	Endpoint   string `koanf:"endpoint"`
	AllowedIPs string `koanf:"allowed-ips"`
	PubKey     string `koanf:"public-key"`
}

type WireguardConfig struct {
	Name        string          `koanf:"name"`
	PrivKeyPath string          `koanf:"private-key-path"`
	Net         string          `koanf:"net"`
	Peers       []WireguardPeer `koanf:"peers"`
}

func NewWireguardInstance(wgConfig WireguardConfig) *WireguardInstance {
	newInstance := WireguardInstance{ready: false, WireguardConfig: wgConfig}
	return &newInstance
}

func (wg *WireguardInstance) Start() {
	log.Println("creating interface")
	common.Run("/sbin/ip", "link", "add", wg.Name, "type", "wireguard")
	log.Println("adding ip")
	common.Run("/sbin/ip", "addr", "add", wg.Net, "dev", wg.Name)
	log.Println("setting key")
	common.Run("/usr/bin/wg", "set", wg.Name, "listen-port", "51820", "private-key", wg.PrivKeyPath)
	log.Println("adding peers")
	for _, peer := range wg.Peers {
		wg.addPeer(peer)
	}
	log.Println("set interface up")
	common.Run("/sbin/ip", "link", "set", "up", "dev", wg.Name)
	for _, peer := range wg.Peers {
		if peer.Endpoint != "" {
			common.Run("/bin/ping", peer.Endpoint)
		}
	}
	log.Println("starting watch")
	wg.cmd, wg.ctx = common.Start("/bin/watch", "/usr/bin/wg")
	wg.ready = true
}

func (wg *WireguardInstance) addPeer(peer WireguardPeer) {
	cmd := []string{"/usr/bin/wg", "set", wg.Name, "peer", peer.PubKey, "allowed-ips", peer.AllowedIPs, "persistent-keepalive", "120"}
	if peer.Endpoint != "" {
		cmd = append(cmd, []string{"endpoint", peer.Endpoint}...)
	}
	common.Run(cmd...)
}

func (wg *WireguardInstance) AddPeer(peer WireguardPeer) {
	wg.addPeer(peer)
	wg.Peers = append(wg.Peers, peer)
}

func (wg *WireguardInstance) Stop() {
	log.Println("stopping instance")
	if wg.ready {
		wg.cmd.Cancel()
		log.Println("wait for process to end")
		wg.cmd.Process.Wait()
	}
}

func main() {
	mqttClient := depman.NewMQTT("wireguard")
	pmClient := client.NewPMClientWithClient("wireguard", mqttClient)
	pmClient.WaitFor("rootfs")

	wgConfig := WireguardConfig{}
	setupConfig(&wgConfig)

	wg := NewWireguardInstance(wgConfig)
	go wg.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println("received signal")
		fmt.Println(sig)
		wg.Stop()

		done <- true
	}()

	// chrome is not running anymore
	//if wg.cmd.ProcessState != nil {
	//	return
	//}

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

}

func setupConfig(wgConfig *WireguardConfig) {
	var k = koanf.New(".")
	if err := k.Load(file.Provider("/etc/goynzane.yaml"), yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	k.Unmarshal("net.wireguard", wgConfig)

}
