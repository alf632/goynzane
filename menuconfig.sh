#!/bin/sh
kas-container shell meta-raspberrypi/bundle.yml -c 'bitbake -c menuconfig linux-raspberrypi'
