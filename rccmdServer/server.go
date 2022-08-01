package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Event int

const (
	Failure Event = iota
	Restored
)

func (e Event) String() string {
	return []string{"Failure", "Restored"}[e]
}

type Command struct {
	Evt  Event
	Time float64
}

// https://stackoverflow.com/questions/26028700/write-to-client-udp-socket-in-go
// https://ops.tips/blog/udp-client-and-server-in-go/

var ErrMessageParsing = errors.New("Could not parse event message")
var ErrUnkownEventType = errors.New("Unknown event type")
var ErrBatteryDuration = errors.New("Could not read battery duration on Powerfail event")

func processPacket(addr *net.UDPAddr, upsIP net.IP, message string) {
	log.WithFields(
		log.Fields{
			"ip":      addr.IP,
			"content": message,
		},
	).Info("Packet received")

	if addr.IP.Equal(upsIP) {
		command, err := ParseMessage(message)
		if err != nil {
			log.WithFields(
				log.Fields{
					"error": err,
				},
			).Error("Parsing error")
			return
		}
		switch command.Evt {
		case Failure:
			scheduleShutdown(command.Time)
		case Restored:
			cancelShutdown()
		}
	} else {
		log.WithFields(
			log.Fields{
				"ip":      addr.IP,
				"content": message,
			},
		).Error("Packet not sent from UPS IP address")
	}
}

func scheduleShutdown(time float64) {
	scheduleTime := int(math.Round(time * 0.5))

	log.WithFields(
		log.Fields{
			"time": scheduleTime,
		},
	).Debug("Scheduling shutdown")

	out, err := exec.Command("shutdown", "-P", strconv.Itoa(scheduleTime)).Output()

	if err != nil {
		fmt.Printf("%v", err)
		log.WithFields(
			log.Fields{
				"error": err,
			},
		).Error("Error scheduling shutdown")
	}

	message := string(out)
	log.WithFields(
		log.Fields{
			"output": message,
		},
	).Debug("Command output")
}

func cancelShutdown() {
	log.Debug("Cancelling scheduled shutdown")

	out, err := exec.Command("shutdown", "-c").Output()
	if err != nil {
		fmt.Printf("%v", err)
		log.WithFields(
			log.Fields{
				"error": err,
			},
		).Error("Error cancelling scheduled shutdown")
	}
	message := string(out)
	log.WithFields(
		log.Fields{
			"output": message,
		},
	).Debug("Command output")
}

func findTestHostIP() net.IP {
	ips, err := net.LookupIP("monitor") // must be the hostname used in docker
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		os.Exit(1)
	}
	if len(ips) != 1 {
		fmt.Fprintf(os.Stderr, "Unexpected amount of IPs: %v\n", len(ips))
		os.Exit(1)
	}
	upsIP := net.ParseIP(ips[0].String())
	return upsIP
}

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
		upsIP = findTestHostIP()
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
		go processPacket(remoteaddr, upsIP, message)
	}
}

func ParseMessage(message string) (Command, error) {
	fmt.Println(message)

	r, _ := regexp.Compile(`CONNECT TO RCCMD \(SEND\) MSG_EXECUTE MSG_TEXT (Powerfail) on (.+): (.+)`)
	result := r.FindStringSubmatch(message)

	if len(result) == 4 {
		event := result[1]
		model := result[2]
		rccmdMessage := result[3]

		log.WithFields(
			log.Fields{
				"event": event,
				"model": model,
				"rccmd": rccmdMessage,
			},
		).Debug("RCCMD message")

		restoreRegexp := regexp.MustCompile(`^restored.$`)
		failureRegexp := regexp.MustCompile(`^Autonomietime ([0-9]+\.[0-0]+) min.$`)

		switch {
		case restoreRegexp.MatchString(rccmdMessage):
			return Command{Restored, 0}, nil
		case failureRegexp.MatchString(rccmdMessage):
			match := failureRegexp.FindStringSubmatch(rccmdMessage)
			time, err := strconv.ParseFloat(match[1], 64)
			// With regexp used before, float conversion can not fail, just to be on the safe side...
			if err != nil {
				return Command{}, ErrBatteryDuration
			}
			return Command{Failure, time}, nil
		default:
			return Command{}, ErrUnkownEventType
		}
	}
	return Command{}, ErrMessageParsing
}
