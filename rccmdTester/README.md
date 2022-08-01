# rccmdTester README
Sends UDP packets with RCCMD messages to broadcast using `gopacket` and injecting them in the network interface.

The interface name and broadcast IP are computed by the program.

Usage:
```
Usage of ./rccmdTester:
  -battery int
    	Battery time
  -failure
    	Power failure event
  -restored
    	Power restored event
```

User must choose between power failure or power restore events. Examples:

```
./rccmdTest -failure -battery 120
./rccmdTest -restored
```
