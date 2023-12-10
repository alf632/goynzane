package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	
	"github.com/gokrazy/gokrazy"
)

func main() {
	BindMount("/perm/home", "/home")
	BindMount("/perm/log", "/var/log")
	BindMount("/perm/websockify", "/usr/bin/websockify/")
	run("/etc/init.d/rcS", "")
	run("/etc/init.d/rc", "5")
	runWithEnv("DISPLAY=:0", "/usr/bin/xset", "-dpms", "s", "off", "s", "noblank", "s", "0", "0", "s", "noexpose")
	gokrazy.WaitForClock()
	runWithEnv("DISPLAY=:0", "chromium", "--kiosk", "-a", "floor796.com")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		run("/etc/init.d/rc", "6")
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

}

func BindMount(src, tgt string) {
	err := syscall.Mount(src, tgt, "", syscall.MS_BIND, "")
	if err != nil {
		fmt.Printf("Unable to bind mount %s to %s\t%s\n",
			src, tgt, err)
	}
}

func run(command ...string) {
	name := command[0]
	args := []string{}
	if len(command) > 1 {
		args = command[1:]
	}
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("%v: %v", cmd.Args, err)
	}
}

func runWithEnv(env, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, env)
	if err := cmd.Run(); err != nil {
		fmt.Printf("%v: %v", cmd.Args, err)
	}
}

