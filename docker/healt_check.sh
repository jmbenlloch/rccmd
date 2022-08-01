#!/bin/bash

RCCMD=`lsof -i -P | grep 6003 | wc -l`
SSH=`lsof -i -P | grep 22 | wc -l`

if [ $(( ($RCCMD != 0) && ($SSH != 0) )) -eq 0 ]; then
	exit 1
fi
