# rccmdServer README
UDP server listening for RCCMD messages. Must be run using the `-src` flag to specify the UPS IP address.

Writes logs in json format to stdout.

A new go routine is launched for each UDP packet read. This has to be this way, if not, the server would block and it could loose following RCCMD messages, breaking completely the program in the case of blocking shutdowns.

Usage:
```
Usage of ./rccmdServer:
  -debug
    	Set logging level to DEBUG
  -dst string
    	broadcast IP receiving RCCMD messages (default "0.0.0.0")
  -src string
    	UPS IP sending RCCMD messages (default "0.0.0.0")
  -test
    	Test mode with docker
```
