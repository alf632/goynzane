#!/bin/bash

INTERFACE="usb0"
ip a | grep -q $INTERFACE
valid=$?
if [[ valid -ne 0 ]] ;then
    echo "running startup"
    /opt/surfstick/startup.sh&
fi

function startDHCP() {
    while true
    do
        echo "starting dhclient on interface $INTERFACE"
        #/gokrazy/dhcp -interface=$INTERFACE -extra_route_priority=99
	dhcpcd -B -m 99 $INTERFACE
        sleep 5
    done
}

# disable IPv6 since there was a bunch of traffic going on
ip6tables -A INPUT -i $INTERFACE -j DROP
ip6tables -A OUTPUT -i $INTERFACE -j DROP
startDHCP
