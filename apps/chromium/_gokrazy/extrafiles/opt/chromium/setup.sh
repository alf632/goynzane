#!/bin/sh

configfile="/perm/home/.config/chromium/Default/Preferences"

echo "disable chromium login"
cat $configfile | jq '.signin.allowed = false' > $configfile.new
mv  $configfile.new $configfile

echo "calibrate touch"
DISPLAY=:0 xinput set-prop "eGalax Inc. USB TouchController" --type=float "Coordinate Transformation Matrix" 0.007714, 1.173833, -0.075632, -1.400702, 0.028842, 1.212083, 0.000000, 0.000000, 1.000000
