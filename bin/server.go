package main

import "nli-go/lib/server"

func main() {
	server := server.NewServer("3333")

	server.Run()
}
