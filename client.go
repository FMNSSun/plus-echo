package main

import "os"
import "fmt"
import "net"
import "github.com/mami-project/plus-lib"
import "flag"

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
	localAddr := flag.String("local-addr", "", "Local address:port to listen on.")
	remoteAddr := flag.String("remote-addr", "", "Remote address:port to connect to.")

	fmt.Println("[CLIENT]")

	flag.Parse()

	if *localAddr == "" || *remoteAddr == "" {
		flag.Usage()
		return
	}

	PLUS.LoggerDestination = os.Stdout

	client(*localAddr, *remoteAddr)
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
	connection.SetCryptoContext(&cryptoContext{key:0x3B})

	connectionManager.SetTransparentMode()

	go connectionManager.Listen()

	var buffer []byte

	for i := 0; i < 128; i++ {

		buffer = []byte{0x65, 0x66, 0x67, 0x68}
		_, err = connection.Write(buffer)

		if err != nil {
			fmt.Printf("[CLIENT] Error: %s\n", err.Error())
			return
		}

		fmt.Printf("[CLIENT] Sent: %q\n", buffer)

		connection.QueuePCFRequest(0x01, 0, []byte{0x00})

		n, err := connection.Read(buffer)

		if err != nil {
			fmt.Printf("[CLIENT] Error: %s\n", err.Error())
			return
		}

		buffer = buffer[:n]

		fmt.Printf("[CLIENT] Got: %q\n", buffer)
	}

	connection.Close()
	connectionManager.Close()
}

func showUsage() {
}
