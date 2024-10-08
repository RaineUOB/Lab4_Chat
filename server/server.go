package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strings"
	"time"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
	fmt.Println("oopsie: ", err)
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
	for true {
		conn, err := ln.Accept()
		if err != nil {
			handleError(err)
		}
		conns <- conn
	}

}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	reader := bufio.NewReader(client)

	open := true

	for open {
		msg, err := reader.ReadString('\n')
		if err != nil {
			handleError(err)
		}

		if strings.Compare(msg, ".close\n") == 0 {
			//client.Close()
			fmt.Println("closed")
			open = false

		} else {
			fmt.Fprintln(client, "ack")
			fmt.Println("handled: ", msg, " for: ", clientid)

			msgs <- Message{clientid, msg}
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.
	ln, _ := net.Listen("tcp", *portPtr)

	//conn, _ := ln.Accept()
	//reader := bufio.NewReader(conn)
	//msg, _ := reader.ReadString('\n')
	//fmt.Println(msg)

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	clientids := 0

	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:

			fmt.Println("new connect: ", conn.RemoteAddr())
			//fmt.Println(clients)
			go handleClient(conn, clientids, msgs)

			clients[clientids] = conn
			clientids++

			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients map
			// - start to asynchronously handle messages from this client
		case msg := <-msgs:

			fmt.Println(msg.message, " From ", msg.sender)

			for _, conn := range clients {
				if conn != clients[msg.sender] && conn != nil {
					fmt.Fprintln(conn, msg.message)
				}
				//fmt.Fprintln(conn, msg.message)
			}

			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
		}
	}

}
