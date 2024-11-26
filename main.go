package main

import (
	"log"
	"net"
)

func main() {
	server := newServer()
	go server.run()

	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server %s", err.Error())
	}

	defer listen.Close()
	log.Printf("Server started on :8888")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("unable to accept conn of %s", err.Error())
			continue
		}

		go server.newClient(conn)
	}
}
