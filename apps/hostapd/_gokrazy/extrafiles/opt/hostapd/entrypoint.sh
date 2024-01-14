#!/bin/bash

MAC=$1
sleep 1
echo "searching for $MAC"
INTERFACE=$(ip a | grep -B1 "$MAC" | head -n1 | cut -f2 -d " " | sed -e "s/://")
echo $INTERFACE | grep -q "wlan"
valid=$?
if [[ valid -ne 0 ]] ;then
    echo "no valid interface found"
    exit 1
fi
echo "found $INTERFACE"

echo "starting hostapd on $INTERFACE"
/usr/sbin/hostapd -i $INTERFACE /perm/hostapd/hostapd.conf
