# General idea

The UPS used has a network interface that monitors the system and is able to send alerts to external servers. The card used is a [CS141](https://www.salicru.com/files/documentacion/salicru_cs141_graphical_quickstart_cmyk_190122.pdf).

The supported events are *Powerfail* and *Power restore*. The idea is that ths UPS will send an UDP packet to broadcast in an internal network (separated VLAN) that will be received by all the necessary servers.  On these servers, a service will be listening for those messages and scheduling shutdowns.

The UPS is able to send the expected battery time. In case of failure, a shutdown will be scheduled for half the battery duration. If the power is restored before that moment, the shutdown will be cancelled.

# Structure of the repository
This repository is organised in the following way:

- **rccmdServer**: Go server to listen for RCCMD messages from the UPS.
- **rccmdTester**:  Go client that sends RCCMD messages to test the server
- **docker**: Docker files and scripts to test the server without a real UPS.
- **compile.sh**: Script to compile both packages and run server's unit tests.
- **automate.sh**: Script to build the docker image, compile the code, deploy the network wtih docker compose and run the system tests.

## Testing the system
Just run `bash automate.sh`

## Setup for development
Run `bash automate.sh debug`  to compile and deploy using docker compose. Then attach a terminal to the desired container. For example:

`docker exec -ti docker exec -ti docker-server1-1 bash`

To check whether UDP packets are arriving or not, tcpdump can be used within the containers: `tcpdump -u -i eth2 udp port 6003 -X`

### Docker setup
There are three containers defined:

- **server1**: Runs `rccmdServer`, listening for UDP packets to the port 6003. Simulates the shutdown in a *non-blocking* way.
- **server2**: Runs `rccmdTester` simulating the behaviour of the UPS, sending RCCMD messages to broadcast. Simulates the shutdown in *blocking way*.
- **monitor**: Additional server than can be used in development to send packets, analyze the network, etc. In the automated tests it does not have any function.
- **sut**: System under test. Ephemeral container that runs the script with the desired system test.

## Shutdown scheduling mechanisms
Currently only implemented for linux systems. To schedule the shutdown the command used is: `shutdown -P [time]`.  There are to possible behaviours that we call:
- **Blocking shutdown**: The shutdown command stays in execution until the timeout, then the shutdown happens. The process runs during the whole time. This happens in older systems like the ones running DATE.
- **Non-blocking shutdown**: The shutdown command exits immediately and creates a special file with the shutdown time in `/run/systemd/shutdown/scheduled`.

# Installation
TODO: Describe how to create the service in the different cases.

# RCCMD messages
CS141 allows for different actions to be configured for many different events. The relevant ones for this system are these:

- Power failure: `CONNECT TO RCCMD (SEND) MSG_EXECUTE MSG_TEXT Powerfail on SLC-20-CUBE3 (5): Autonomietime %v.00 min.`
- Power restore: `CONNECT TO RCCMD (SEND) MSG_EXECUTE MSG_TEXT Powerfail on SLC-20-CUBE3 (5): restored.`
