package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/blang/speedtest"
	"github.com/eyedeekay/sam3/helper"
)

const CHECK_INTERVAL = 3 // Times BUFFER_SIZE

func main() {
	listen := flag.String("listen", ":8000", "Address to listen on")
	buffersize := flag.Int("buffer", 4096, "Buffer size")
	reportinterval := flag.Duration("report", 5*time.Second, "Report interval")
	send := flag.Bool("send", false, "True for send, false for receive")
	i2p := flag.Bool("i2p", false, "Run over an I2P Streaming Service, overrides -listen")
	flag.Parse()

	var ln net.Listener
	var err error

	if *i2p {
		ln, err = sam.I2PListener("speedtest", "127.0.0.1:7656", "speedtest")
		if err != nil {
			log.Fatalf("Could not listen on %s: %s", *listen, err)
		}
	} else {
		ln, err = net.Listen("tcp", *listen)
		if err != nil {
			log.Fatalf("Could not listen on %s: %s", *listen, err)
		}
	}
	buffer := make(chan speedtest.BytesPerTime)
	output := make(chan speedtest.BytesPerTime)
	speedtest.SpeedMeter(buffer, output) // Speedmeter on all connections
	speedtest.SpeedReporter(output, *reportinterval)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error on connection: %s", err)
			continue
		}
		log.Printf("Accepted connection from %s", conn.RemoteAddr())
		go handleConnection(conn, *send, buffer, *buffersize)
	}

}

func handleConnection(conn net.Conn, send bool, reportCh chan speedtest.BytesPerTime, buffersize int) {
	if send {
		log.Println("Enter Send mode")
		err := speedtest.SendData(conn, buffersize, reportCh)
		if err != nil {
			log.Printf("Error: %s", err)
			return
		}
	} else {
		log.Println("Enter Receive mode")
		err := speedtest.ReceiveData(conn, buffersize, reportCh)
		if err != nil {
			log.Printf("Error: %s", err)
			return
		}
	}
}
