package lib

import (
	"fmt"
	"log"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type WaitReq struct {
	Name   string
	Target string
}

type WaitRet struct{}

type ProvideReq struct {
	Name   string
	Target string
}

type ProvideRet struct{}

type watcher struct {
	target string
	done   chan struct{}
}

type DependencyManager struct {
	watchers map[string][]*watcher
	provided []string

	mu   sync.Mutex
	mqtt mqtt.Client
}

func NewDependencyManager() *DependencyManager {
	newDM := DependencyManager{watchers: map[string][]*watcher{}, provided: []string{}, mqtt: NewMQTT("depman")}
	return &newDM
}

func (dm *DependencyManager) WaitFor(wreq WaitReq, wret *WaitRet) error {
	target := wreq.Target
	log.Println(wreq.Name, "is waiting for", target)
	done := make(chan struct{})
	dm.mu.Lock()
	newWatcher := watcher{target: target, done: done}
	// test if target is already watched and add to - or create - list of watchers
	_, exists := dm.watchers[target]
	if !exists {
		dm.watchers[target] = []*watcher{&newWatcher}
	} else {
		dm.watchers[target] = append(dm.watchers[target], &newWatcher)
	}

	for _, prov := range dm.provided {
		if prov == target {
			return nil
		}
	}
	dm.mu.Unlock()

	return nil //since stuff is not working..
	<-done
	log.Println(wreq.Name, "is done waiting for", target)
	return nil
}

func (dm *DependencyManager) Provide(preq ProvideReq, pret ProvideRet) error {
	target := preq.Target
	log.Println(preq.Name, "is providing", target)

	go func() {
		dm.mu.Lock()
		watchers, exists := dm.watchers[target]
		if exists {
			for _, watcher := range watchers {
				watcher.done <- struct{}{}
			}
		}

		dm.provided = append(dm.provided, target)
		dm.mu.Unlock()
	}()

	return nil
}

// MQTT

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func NewMQTT(clientID string) mqtt.Client {
	var broker = "127.0.0.1"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientID)
	//opts.SetUsername("emqx")
	//opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}
