#!/bin/bash

while true; do
    echo "checking devices..."

    if [[ $(lsusb -d 3566:2001) ]]; then
        echo "Huawei E3372-325 detected. Performing mode switch..."
        /usr/sbin/usb_modeswitch -v 3566 -p 2001 -W -R -w 400
        /usr/sbin/usb_modeswitch -v 3566 -p 2001 -W -R
        echo "Mode switch completed."
        break
    fi
    sleep 1  # Wait for 1 second before checking again
done