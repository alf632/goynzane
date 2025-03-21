#!/bin/sh
MAC=$1
IP=$2
NET=$3
PREFIX=$4

if [ -z "$MAC" ]; then
  INTERFACE=""
else
  INTERFACE=$(ip a | grep -B1 "$MAC" | head -n1 | cut -f2 -d " " | sed -e "s/://")
fi

if [ -z "$INTERFACE" ]; then

  echo "starting dnsmasq"
  /usr/bin/dnsmasq -d --conf-file=/perm/dnsmasq/dnsmasq.conf

else

  ip link set up dev $INTERFACE
  ip a a $IP$PREFIX dev $INTERFACE
  
  sysctl -w net.ipv4.ip_forward=1
  iptables -t nat -A POSTROUTING -s $NET$PREFIX -j MASQUERADE
  iptables -P FORWARD ACCEPT
  iptables -F FORWARD
  
  echo "starting dnsmasq on interface $INTERFACE"
  /usr/bin/dnsmasq -d -i $INTERFACE --conf-file=/perm/dnsmasq/dnsmasq.conf

fi
