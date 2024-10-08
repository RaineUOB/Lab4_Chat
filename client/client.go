package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"time"
)

func read(conn net.Conn, close chan bool) {
	//TODO In a continuous loop, read a message from the server and display it.
	reader := bufio.NewReader(conn)
	open := true

	for open {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println("handled: ", msg, " for: ", clientid)
		fmt.Println("from server: ", msg)

		if msg == ".close" {
			open = false
			close <- true
		}

	}
}

func write(conn net.Conn, close chan bool) {
	//TODO Continually get input from the user and send messages to the server.
	open := true
	var msgP string
	for open {
		fmt.Scanln(&msgP)

		//msg := "Mic check Mic check"

		fmt.Fprintln(conn, msgP)
		//fmt.Println(msgP)
		//fmt.Println("sent")

		if msgP == ".close" {
			open = false
			close <- true
		}
		//open = false
	}
}

func main() {
	// Get the server address and port from the commandline arguments.

	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()
	conn, _ := net.Dial("tcp", *addrPtr)
	//msgP := flag.String("msgP", "hello world", "message to send (.close to close)")

	close := make(chan bool)

	go read(conn, close)
	go write(conn, close)
	_ = <-close

	fmt.Println("Goodbye")

	time.Sleep(1 * time.Second)

}
