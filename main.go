package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

// Common wake on lan ports
var ports = []int{7, 9}

func constructMagicPacket(macAddress net.HardwareAddr) []byte {
	magic := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	for i := 0; i < 16; i++ {
		magic = append(magic, macAddress...)
	}
	return magic
}

func sendMagicPacket(port int, magic []byte, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		// Broadcast ip
		IP:   net.IPv4(255, 255, 255, 255),
		Port: port,
	})
	if err != nil {
		errChan <- fmt.Errorf("Error setting up UDP connection: %v", err)
		return
	}
	_, err = conn.Write(magic)
	if err != nil {
		errChan <- fmt.Errorf("Error sending magic packet: %v", err)
	}

	_ = conn.Close()
}

func wol(macAddressRaw string) error {

	// Validate mac address
	macAddress, err := net.ParseMAC(macAddressRaw)
	if err != nil {
		fmt.Printf("Error parsing MAC address: %v\n", err)
		return err
	}

	fmt.Println("MAC Address:", macAddress)

	magicPacket := constructMagicPacket(macAddress)

	var wg sync.WaitGroup
	errChan := make(chan error, len(ports))

	for _, port := range ports {
		wg.Add(1)
		go sendMagicPacket(port, magicPacket, &wg, errChan)
	}

	wg.Wait()
	close(errChan)

	// Check for errors from goroutines
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s  <MAC Address>\n", os.Args[0])
		os.Exit(1)
	}

	err := wol(os.Args[1])
	if err != nil {
		fmt.Printf("Error sending magic packet: %v\n", err)
		os.Exit(1)
	}
}
