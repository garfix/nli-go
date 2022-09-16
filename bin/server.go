package main

import (
	"nli-go/lib/server"
	"os"
)

func main2() {

	if len(os.Args) != 2 {
		println("Use: ./server <port>")
		os.Exit(1)
	}

	port := os.Args[1]

	server := server.NewServer(port)

	server.Run()
}
