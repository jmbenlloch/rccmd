#!/bin/bash
cd /root/rccmd/rccmdServer/
mage

echo "Running unit tests"
go test ./pkg
if [ $? -ne 0 ]; then
	echo "Error in unit tests"
	exit 1
fi

cd /root/rccmd/rccmdTester
go build .
