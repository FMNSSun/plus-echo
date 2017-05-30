package main

import "os"
import "fmt"
import "net"
import "plus"

type cryptoContext struct {
	key byte
}

func (c *cryptoContext) EncryptAndProtect(header []byte, payload []byte) ([]byte, error) {
	for i, v := range payload {
		payload[i] = v ^ c.key
	}

	return payload, nil
}

func (c *cryptoContext) DecryptAndValidate(header []byte, payload []byte) ([]byte, bool, error) {
	for i, v := range payload {
		payload[i] = v ^ c.key
	}

	return payload, true, nil
}

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
		conn.SetCryptoContext(&cryptoContext{key:0x3B})
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

				_, err = connection.Write(buffer)

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
