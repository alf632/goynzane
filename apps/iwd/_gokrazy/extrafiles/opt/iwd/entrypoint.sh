#!/bin/bash

MAC=$1
sleep 1
echo "searching for mac $MAC"
INTERFACE=$(ip a | grep -B1 "$MAC" | head -n1 | cut -f2 -d " " | sed -e "s/://")
echo $INTERFACE | grep -q "wlan"
valid=$?
if [[ valid -ne 0 ]] ;then
    echo "no valid interface found"
    exit 1
fi
echo "found interface $INTERFACE"

function startClient(){
    sleep 5
    echo "starting Client"
    iwctl station $INTERFACE scan
    startDHCP&
}

function startDHCP() {
    while true
    do
        echo "starting dhclient on interface $INTERFACE"
        /gokrazy/dhcp -interface=$INTERFACE
        sleep 5
    done
}

startClient&
/usr/libexec/iwd
