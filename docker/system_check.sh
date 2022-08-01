#!/bin/bash

# Before sending RCMMD message, no scheduled shutdown must be in place
ssh -o StrictHostKeyChecking=no server1 /root/rccmd/docker/check_scheduled_shutdown.sh
result1=$?
echo "Scheduled shutdown for server 1 $result1"
ssh -o StrictHostKeyChecking=no server2 /root/rccmd/docker/check_scheduled_shutdown.sh
result2=$?
echo "Scheduled shutdown for server 2 $result1"
if [ $(( ($result1 != 0) || ($result2 != 0) )) -eq 1 ]; then
	exit 1
fi

# Send the power failure command to broadcast
echo "Sending power failure RCCMD message to broadcast"
/root/rccmd/rccmdTester/rccmdTester -failure -battery 60

# Both servers must have a scheduled shutdown
ssh -o StrictHostKeyChecking=no server1 /root/rccmd/docker/check_scheduled_shutdown.sh
result1=$?
echo "Scheduled shutdown for server 1 $result1"
ssh -o StrictHostKeyChecking=no server2 /root/rccmd/docker/check_scheduled_shutdown.sh
result2=$?
echo "Scheduled shutdown for server 2 $result2"
if [ $(( ($result1 != 1) || ($result2 != 1) )) -eq 1 ]; then
	echo exiting
	exit 1
fi

# Send the power recovered command to broadcast
echo "Sending power restored RCCMD message to broadcast"
/root/rccmd/rccmdTester/rccmdTester -restored

# Both servers must have no scheduled shutdown
ssh -o StrictHostKeyChecking=no server1 /root/rccmd/docker/check_scheduled_shutdown.sh
result1=$?
echo "Scheduled shutdown for server 1 $result1"
ssh -o StrictHostKeyChecking=no server2 /root/rccmd/docker/check_scheduled_shutdown.sh
result2=$?
echo "Scheduled shutdown for server 2 $result2"
if [ $(( ($result1 != 0) || ($result2 != 0) )) -eq 1 ]; then
	exit 1
fi
