package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	dependencyManagerLib "github.com/alf632/goynzane/apps/dependencyManager/lib"
)

func main() {
	dm := dependencyManagerLib.NewDependencyManager()
	rpcServer := rpc.NewServer() // HERE WE NEED TO CREATE NEW RPC SERVER
	rpcServer.Register(dm)       // AND WE SHOULD USE IT TO REEGISTER OUR METHODS
	la, e := net.Listen("tcp", "127.0.0.1:16363")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go rpcServer.Accept(la)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		select {
		case sig := <-sigs:
			fmt.Println()
			fmt.Println(sig)
		}

		done <- true
	}()

	fmt.Println("initialized")
	<-done
}
