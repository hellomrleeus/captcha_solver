#!/bin/bash
CONN=$1
IP=$2
chmod +x client
./client -conn=$CONN ip=$IP
