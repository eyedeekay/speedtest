package main

import (
	"flag"
	"github.com/blang/speedtest"
	"log"
	"net"
	"time"
        "github.com/eyedeekay/sam3"
)

func main() {
	connect := flag.String("connect", "", "connect to address")
	buffersize := flag.Int("buffer", 4096, "Buffer size")
	reportinterval := flag.Duration("report", 5*time.Second, "Report interval")
	send := flag.Bool("send", true, "True for send, false for receive")
        i2p := flag.Bool("i2p", false, "Use I2P address for speedtest")
	flag.Parse()

        var conn net.Conn
        var err error
        if *i2p {
            sess, err := sam.I2PStreamSession("speedtest-client", "127.0.0.1:7656", "speedtest-client")
            if err != nil {
		log.Fatalf("Could not connect: %s", err)
	    }
            conn, err = sess.DialI2P("I2P", *connect)
            if err != nil {
		log.Fatalf("Could not connect: %s", err)
	    }
        } else {
	    conn, err = net.Dial("tcp", *connect)
	    if err != nil {
		log.Fatalf("Could not connect: %s", err)
	    }
        }

	reportCh := make(chan speedtest.BytesPerTime)
	statsCh := make(chan speedtest.BytesPerTime)
	speedtest.SpeedMeter(reportCh, statsCh) // Speedmeter on all connections
	speedtest.SpeedReporter(statsCh, *reportinterval)

	if *send { // Client send mode
		log.Println("Enter Send mode")
		err := speedtest.SendData(conn, *buffersize, reportCh)
		if err != nil {
			log.Printf("Error: %s", err)
			return
		}
	} else { // Receive mode
		log.Println("Enter Receive mode")
		err := speedtest.ReceiveData(conn, *buffersize, reportCh)
		if err != nil {
			log.Printf("Error: %s", err)
			return
		}
	}

}
