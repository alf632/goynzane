module github.com/alf632/goynzane/apps/mosquitto

replace github.com/alf632/goynzane/apps/common => ../common

replace github.com/alf632/goynzane/apps/dependencyManager => ../dependencyManager

replace github.com/alf632/goynzane/apps/dependencyManager/lib => ../dependencyManager/lib

replace github.com/alf632/goynzane/apps/dependencyManager/client => ../dependencyManager/client

go 1.21.5

require (
	github.com/alf632/goynzane/apps/common v0.0.0-00010101000000-000000000000
	github.com/alf632/goynzane/apps/dependencyManager/client v0.0.0-00010101000000-000000000000
)

require (
	github.com/alf632/goynzane/apps/dependencyManager/lib v0.0.0-00010101000000-000000000000 // indirect
	github.com/eclipse/paho.mqtt.golang v1.4.3 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
)
