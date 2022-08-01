package main

import (
	"errors"
	"net"
)

type NetworkData struct {
	IP        net.IP
	Broadcast net.IP
	Name      string
	Mac       net.HardwareAddr
}

func GetNetworkData() (NetworkData, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return NetworkData{}, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return NetworkData{}, err
		}
		for _, addr := range addrs {
			var ip net.IP
			var mask net.IPMask

			value, ok := addr.(*net.IPNet)
			if ok == false {
				continue // not an IPNet, probably ipv6
			}

			ip = value.IP
			mask = value.Mask

			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			broadcastAddress := computeBroadcastAddress(ip, mask)
			result := NetworkData{ip, broadcastAddress,
				iface.Name, iface.HardwareAddr}
			return result, nil
		}
	}
	return NetworkData{}, errors.New("are you connected to the network?")
}

func computeBroadcastAddress(ip net.IP, mask net.IPMask) net.IP {
	ip = ip.To4()
	tempAddress := [4]byte{}

	for i := 0; i < 4; i++ {
		tempAddress[i] = ip[i] | ^mask[i]
	}

	broadcastAddress := net.IPv4(
		tempAddress[0], tempAddress[1],
		tempAddress[2], tempAddress[3])
	return broadcastAddress
}
