module github.com/alf632/goynzane/apps/wireguard

go 1.21.5

replace github.com/alf632/goynzane/apps/common => /home/malte/Projects/goynzane/apps/common

replace github.com/alf632/goynzane/apps/dependencyManager => /home/malte/Projects/goynzane/apps/dependencyManager

replace github.com/alf632/goynzane/apps/dependencyManager/lib => /home/malte/Projects/goynzane/apps/dependencyManager/lib

replace github.com/alf632/goynzane/apps/dependencyManager/client => /home/malte/Projects/goynzane/apps/dependencyManager/client

require (
	github.com/alf632/goynzane/apps/common v0.0.0-00010101000000-000000000000
	github.com/gokrazy/gokrazy v0.0.0-20231012060754-0dcab4a1b850
)

require (
	github.com/alf632/goynzane/apps/dependencyManager/lib v0.0.0-00010101000000-000000000000 // indirect
	github.com/eclipse/paho.mqtt.golang v1.4.3 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/alf632/goynzane/apps/dependencyManager/client v0.0.0-00010101000000-000000000000
	github.com/gokrazy/internal v0.0.0-20220129150711-9ed298107648 // indirect
	github.com/google/renameio/v2 v2.0.0 // indirect
	github.com/kenshaw/evdev v0.1.0 // indirect
	github.com/knadh/koanf v1.5.0
	github.com/mdlayher/watchdog v0.0.0-20201005150459-8bdc4f41966b // indirect
	golang.org/x/sys v0.13.0 // indirect
)
