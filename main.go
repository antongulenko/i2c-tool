package main

import (
	"flag"
	"log"
	"strconv"

	"golang.org/x/exp/io/i2c"
)

func main() {
	device := flag.String("d", "/dev/i2c-0", "The i2c device to open")
	addr := flag.Int("a", 03, "The slave device address")
	var send, recv []byte
	numRecv := 0

	flag.Parse()
	for _, arg := range flag.Args() {
		if arg == "?" {
			numRecv++
		} else if numRecv > 0 {
			log.Fatalln("Cannot send bytes after receiving")
		} else {
			b, err := strconv.ParseUint(arg, 0, 8)
			if err != nil {
				log.Fatalln("Illegal send byte ", arg, ":", err)
			}
			send = append(send, byte(b))
		}
	}
	if numRecv > 0 {
		recv = make([]byte, numRecv)
	}

	dev := i2c.Devfs{Dev: *device}
	conn, err := dev.Open(*addr, false)
	if err != nil {
		log.Fatalln("Error opening slave device connection:", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Warning: error closing slave device connection:", err)
		}
	}()

	log.Println("Sending:", send)
	log.Println("Receiving bytes:", len(recv))

	if err := conn.Tx(send, recv); err != nil {
		log.Println("Transmission error:", err)
	}

	log.Println("Received:", recv)
}
