#!/bin/bash

# This script uses the return code as a boolean value
# If a shutdown is scheduled it returns 1, if not, 0

function check_non_blocking_scheduled_shutdown {
	if [ -f /run/systemd/shutdown/scheduled ]; then 
		return 1
	fi
	return 0
}

function check_blocking_scheduled_shutdown {
	blockingProcess=`ps axu | grep "shutdown -P" | grep -v grep  | wc -l`
	if [ $blockingProcess -gt 0 ]; then
		return 1
	fi
	return 0
}

check_non_blocking_scheduled_shutdown
result1=$?
check_blocking_scheduled_shutdown
result2=$?

result=$(($result1 || $result2))
exit $result
