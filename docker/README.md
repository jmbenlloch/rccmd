# Docker configuration
The `go` folder contains the `Dockerfile` for building the image used to compile and test the code. All the containers in the `docker-compose` file use that image.

Docker compose builds the network with two servers and one client to send the RCCMD messages.

`shutdown_*blocking.sh` simulates the behaviour of Linux shutdown schedule in a (non-)blocking way.

`healt_check.sh` checks that both ssh server and rccmdServer are running. Required before running the `system_test.sh`

`supervisord.conf` is the configuration file to define the different processes to be run in the containers that require more than one (both servers).
