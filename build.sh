#!/bin/sh
cp kiosk.yml meta-raspberrypi/
kas-container dump meta-raspberrypi/kas-poky-rpi.yml:meta-raspberrypi/kiosk.yml > meta-raspberrypi/bundle.yml
kas-container build meta-raspberrypi/bundle.yml
