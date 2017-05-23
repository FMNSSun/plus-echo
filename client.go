package main

import "os"
import "fmt"
import "net"
import "plus"

func main() {
	args := os.Args

	fmt.Println("client")

	if len(args) == 0 {
		showUsage()
		return
	}

	PLUS.LoggerDestination = os.Stdout

	client("localhost:50002", "localhost:50001")
}

func client(laddr string, remoteAddr string) {
	packetConn, err := net.ListenPacket("udp", laddr)

	if err != nil {
		fmt.Printf("[CLIENT] Error: %s\n", err.Error())
		panic("Could not create packet connection!")
	}

	udpAddr, err := net.ResolveUDPAddr("udp4", remoteAddr)

	if err != nil {
		fmt.Printf("[CLIENT] Error: %s\n", err.Error())
		panic("Could not resolve address!")
	}


	connectionManager, connection := PLUS.NewConnectionManagerClient(packetConn, 1989, udpAddr)
	go connectionManager.Listen()

	connection.SetSFlag(true)

	buffer := []byte{0x10, 0x20, 0x30, 0x40, 0x50}
	err = connection.Write(buffer)

	if err != nil {
		fmt.Printf("[CLIENT] Error: %s\n", err.Error())
		return
	}

	fmt.Printf("[CLIENT] Sent: %q\n", buffer)

	n, err := connection.Read(buffer)

	if err != nil {
		fmt.Printf("[CLIENT] Error: %s\n", err.Error())
		return
	}

	buffer = buffer[:n]

	fmt.Printf("[CLIENT] Got: %q\n", buffer)
}

func showUsage() {
}
