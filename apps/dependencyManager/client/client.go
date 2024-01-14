package client

import (
	"fmt"

	dependencyManagerLib "github.com/alf632/goynzane/apps/dependencyManager/lib"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type PMClient struct {
	name string

	mqtt mqtt.Client
}

func NewPMClient(name string) *PMClient {
	pmc := PMClient{name: name, mqtt: dependencyManagerLib.NewMQTT(name)}
	return &pmc
}

func (pmc *PMClient) WaitFor(target string) {
	provided := make(chan struct{})
	fmt.Println("setting up wait for", target)
	targetTopic := fmt.Sprintf("dependencyManager/targets/%s", target)
	pmc.mqtt.Subscribe(targetTopic, 0, func(c mqtt.Client, m mqtt.Message) {
		if string(m.Payload()) == "provided" {
			provided <- struct{}{}
		}
	})
	defer pmc.mqtt.Unsubscribe(targetTopic)

	fmt.Println("waiting for", target)
	<-provided
	fmt.Println("waiting finished for", target)
}

func (pmc *PMClient) Provide(target string) {
	fmt.Println("providing", target)
	targetTopic := fmt.Sprintf("dependencyManager/targets/%s", target)
	pmc.mqtt.Publish(targetTopic, 0, true, "provided")
}
