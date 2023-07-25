package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	rccmd "github.com/jmbenlloch/rccmd/rccmdServer/pkg"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	log.Info("RCCMD server started")

	bcastPtr := flag.String("dst", "0.0.0.0", "broadcast IP receiving RCCMD messages")
	srcPtr := flag.String("src", "0.0.0.0", "UPS IP sending RCCMD messages")
	testPtr := flag.Bool("test", false, "Test mode with docker")
	debugPtr := flag.Bool("debug", false, "Set logging level to DEBUG")
	flag.Parse()

	if *debugPtr {
		log.SetLevel(log.DebugLevel)
	}

	bcastIP := net.ParseIP(*bcastPtr)
	upsIP := net.ParseIP(*srcPtr)

	if *testPtr {
		upsIP = rccmd.FindTestHostIP()
		fmt.Println(upsIP)
	}

	log.WithFields(
		log.Fields{
			"broadcast": bcastIP,
			"ups":       upsIP,
		},
	).Debug("Configuration")

	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 6003,
		IP:   bcastIP,
	}

	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.WithFields(
			log.Fields{
				"error": err,
			},
		).Fatal("Error opening socket")
	}

	for {
		size, remoteaddr, err := server.ReadFromUDP(p)
		if err != nil {
			log.WithFields(
				log.Fields{
					"error": err,
				},
			).Fatal("Error reading from socket")
			continue
		}

		message := string(p[:size])
		go rccmd.ProcessPacket(remoteaddr, upsIP, message)
	}

}
