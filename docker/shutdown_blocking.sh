#!/bin/bash

COMMAND=$1
ARGUMENT=$2

echo shutdown $1 $2

function poweroff {
	echo poweroff $1
	sleep "$1"m
	echo "Shutdown now"
}

function cancel {
	echo cancel
	ps axu | grep "shutdown -P" | grep -v grep | awk '{print $2}' | xargs kill
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
