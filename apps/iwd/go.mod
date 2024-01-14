module github.com/alf632/goynzane/apps/iwd

go 1.21.5

replace github.com/alf632/goynzane/apps/common => /home/malte/Projects/goynzane/apps/common

replace github.com/alf632/goynzane/apps/dependencyManager => /home/malte/Projects/goynzane/apps/dependencyManager

replace github.com/alf632/goynzane/apps/dependencyManager/lib => /home/malte/Projects/goynzane/apps/dependencyManager/lib

replace github.com/alf632/goynzane/apps/dependencyManager/client => /home/malte/Projects/goynzane/apps/dependencyManager/client

require (
	github.com/alf632/goynzane/apps/common v0.0.0-00010101000000-000000000000
	github.com/alf632/goynzane/apps/dependencyManager/client v0.0.0-00010101000000-000000000000
	github.com/knadh/koanf/parsers/yaml v0.1.0
	github.com/knadh/koanf/providers/file v0.1.0
	github.com/knadh/koanf/v2 v2.0.1
)

require (
	github.com/alf632/goynzane/apps/dependencyManager/lib v0.0.0-00010101000000-000000000000 // indirect
	github.com/eclipse/paho.mqtt.golang v1.4.3 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
