#!/bin/bash

COMMAND=$1
ARGUMENT=$2

function poweroff {
	echo poweroff $1
	mkdir -p /run/systemd/shutdown/
	DATE=$(($(date +%s%N)/1000))
	OFFSET=$(($1 * 60 * 1000000))
	POWEROFF=$(($DATE + $OFFSET))
	cat > /run/systemd/shutdown/scheduled << EOF
USEC=$POWEROFF
WARN_WALL=1
MODE=poweroff
EOF
}

function cancel {
	echo cancel
	rm /run/systemd/shutdown/scheduled
}


case $COMMAND in
    -P)       poweroff $ARGUMENT ;;
    -c)       cancel ;;
    *) echo Unrecognized command: ${COMMAND}
       echo
	   echo Usage:
	   echo 
	   echo "shutdown -P MINUTES"
	   echo "shutdown -c"
	   ;;
esac
