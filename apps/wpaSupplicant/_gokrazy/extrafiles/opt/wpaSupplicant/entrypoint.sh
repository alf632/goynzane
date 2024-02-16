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

function startDHCP() {
    while true
    do
        echo "starting dhclient on interface $INTERFACE"
        #/gokrazy/dhcp -interface=$INTERFACE
	dhcpcd -B -m 5 $INTERFACE
        sleep 5
    done
}

if [[ ! -d /perm/wpaSupplicant ]]; then
  echo "copying default config"
  cp -r /opt/wpaSupplicant /perm
fi

# collect configs
CONFIGS=""
for i in $(ls /perm/wpaSupplicant/wifis/*.conf); do echo $i; done                                                                                                                                            
    CONFIGS+="-I$i "

startDHCP&
/usr/sbin/wpa_supplicant -c/perm/wpaSupplicant/wpaSupplicant.conf -i$INTERFACE $CONFIGS
