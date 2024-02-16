module github.com/alf632/goynzane/apps/x11vnc

replace github.com/alf632/goynzane/apps/common => /home/malte/Projects/goynzane/apps/common

replace github.com/alf632/goynzane/apps/dependencyManager => /home/malte/Projects/goynzane/apps/dependencyManager

replace github.com/alf632/goynzane/apps/dependencyManager/lib => /home/malte/Projects/goynzane/apps/dependencyManager/lib

replace github.com/alf632/goynzane/apps/dependencyManager/client => /home/malte/Projects/goynzane/apps/dependencyManager/client

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
