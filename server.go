package main

import "os"
import "fmt"
import "net"
import "plus"

func main() {
	args := os.Args

	fmt.Println("server")

	if len(args) == 0 {
		showUsage()
		return
	}

	PLUS.LoggerDestination = os.Stdout

	server("localhost:50001")
}

func server(addr string) {
	packetConn, err := net.ListenPacket("udp", addr)

	if err != nil {
		panic("Could not create packet connection!")
	}

	connectionManager := PLUS.NewConnectionManager(packetConn)
	connectionManager.SetInitConn(func(conn *PLUS.Connection) error {
		conn.SetSFlag(true)
		return nil
	})

	go connectionManager.Listen()

	for {
		connection := connectionManager.Accept()



		go func() {
			for {
				fmt.Println("[SERVER] Waiting for client")
				buffer := make([]byte, 8129)
				n, err := connection.Read(buffer)
				buffer = buffer[:n]

				fmt.Printf("[SERVER] Received: %q\n", buffer)

				if err != nil {
					fmt.Printf("[SERVER] Error: %s\n", err.Error())
					return
				}

				err = connection.Write(buffer)

				if err != nil {
					fmt.Printf("[SERVER] Error: %s\n", err.Error())
					return
				}
			}
		}()
	}
}

func showUsage() {
}
