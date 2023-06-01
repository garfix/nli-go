package main

import (
	"nli-go/lib/server"
	"os"
)

func main() {

	if len(os.Args) != 4 {
		println("Use: ./server <port> <application-base-dir> <work-dir>")
		os.Exit(1)
	}

	port := os.Args[1]
	baseDir := os.Args[2]
	workDir := os.Args[3]

	server := server.NewServer(port, baseDir, workDir)

	server.Run()
}
