package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alf632/goynzane/apps/common"
	"github.com/alf632/goynzane/apps/dependencyManager/client"
)

func main() {
	common.BindMount("/perm/home", "/home")
	common.BindMount("/perm/log", "/var/log")
	common.BindMount("/perm/cache", "/var/cache")

	//set stage for podman
	common.MountTmpfs("/usr/local", 10)
	os.MkdirAll("/usr/local/bin", os.ModePerm)
	os.Symlink("/usr/bin/podman", "/usr/local/bin/podman")

	common.BindMount("/perm/websockify", "/usr/bin/websockify/")

	// this need /var to exist. gokrazy links /var to /perm/var
	common.Run("/etc/init.d/rcS", "")
	//common.BindMount("/perm/container-storage", "/var/lib/containers")
	os.Mkdir("/tmp/serial-busybox/", os.ModePerm)
	os.Symlink("/bin/bash", "/tmp/serial-busybox/ash")
	//BindMount("/bin/bash", "/tmp/serial-busybox/ash")
	common.Run("/etc/init.d/rc", "5")

	pmClient := client.NewPMClient("rootfs")
	pmClient.Provide("rootfs")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

}
