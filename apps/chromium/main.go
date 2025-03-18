package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alf632/goynzane/apps/common"
	"github.com/alf632/goynzane/apps/dependencyManager/client"
	depman "github.com/alf632/goynzane/apps/dependencyManager/lib"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type ChromiumInstance struct {
	url string

	ready   bool
	exiting bool

	abortWait chan struct{}

	cmd *exec.Cmd
	ctx context.Context
}

type TouchConfig struct {
	Device string `koanf:"device"`
	Matrix string `koanf:"matrix"`
}

type ChromiumConfig struct {
	Url                string      `koanf:"url"`
	DisableScreensaver bool        `koanf:"disableScreensaver"`
	TouchCal           TouchConfig `koanf:"touchCalibration"`
	hasTouchCal        bool
}

var cConfig ChromiumConfig

func NewChromiumInstance(url string) *ChromiumInstance {
	newInstance := ChromiumInstance{ready: false, url: url}
	return &newInstance
}

func (ci *ChromiumInstance) Start() {
	log.Println("waiting for website")
	ci.abortWait = make(chan struct{})
	waitForWebsite(ci.url, ci.abortWait)
	if ci.exiting {
		log.Println("start aborted")
		return
	}
	log.Println("starting instance")
	ci.cmd, ci.ctx = common.StartWithEnv("DISPLAY=:0", "/usr/bin/chromium", "--no-sandbox", "--kiosk", "-a", ci.url)
	ci.ready = true
}

func (ci *ChromiumInstance) Stop() {
	log.Println("stopping instance")
	ci.exiting = true
	if ci.ready {
		ci.cmd.Cancel()
		log.Println("wait for process to end")
		ci.cmd.Process.Wait()
	} else if ci.abortWait != nil {
		log.Println("aborting wait")
		ci.abortWait <- struct{}{}
	}

}

func main() {
	setupConfig()
	mqttClient := depman.NewMQTT("chromium")
	pmClient := client.NewPMClientWithClient("chromium", mqttClient)
	pmClient.WaitFor("rootfs")

	if cConfig.DisableScreensaver {
		log.Println("disabling Screensaver")
		common.RunWithEnv("DISPLAY=:0", "/usr/bin/xset", "-dpms", "s", "off", "s", "noblank", "s", "0", "0", "s", "noexpose")
	}
	if cConfig.hasTouchCal {
		log.Println("calibrating touch with")
		args := []string{"set-prop", cConfig.TouchCal.Device, "--type=float", "Coordinate Transformation Matrix"}
		args = append(args, strings.Split(cConfig.TouchCal.Matrix, " ")...)
		log.Print("xinput ")
		log.Println(args)
		common.RunWithEnv("DISPLAY=:0", "xinput", args...)
	}

	common.Run("/opt/chromium/setup.sh")

	ci := NewChromiumInstance(cConfig.Url)
	mqttClient.Subscribe("apps/chromium/goto", 0, func(c mqtt.Client, m mqtt.Message) {
		log.Println("received request to go to another page")
		log.Println("stopping old instance")
		ci.Stop()
		log.Println("starting new instance")
		ci = NewChromiumInstance(string(m.Payload()))
		go ci.Start()
	})
	go ci.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println("received signal")
		fmt.Println(sig)
		ci.Stop()
		mqttClient.Unsubscribe("apps/chromium/goto")

		done <- true
	}()

	// chrome is not running anymore
	//if ci.cmd.ProcessState != nil {
	//	return
	//}

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

}

func waitForWebsite(url string, abort chan struct{}) {
	for {
		select {
		case <-abort:
			log.Println("wait aborted")
			return
		default:
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
			} else if resp.StatusCode == 200 {
				return
			} else {
				fmt.Println(http.StatusText(resp.StatusCode), "waiting for website", url)
			}
		}
		time.Sleep(time.Second * 5)
	}
}

func setupConfig() {
	var k = koanf.New(".")
	if err := k.Load(file.Provider("/etc/goynzane.yaml"), yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	k.Unmarshal("apps.chromium", &cConfig)
	if k.Exists("apps.chromium.touchCalibration") {
		cConfig.hasTouchCal = true
	}

}
