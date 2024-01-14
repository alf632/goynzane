#!/bin/sh

configfile="/perm/home/.config/chromium/Default/Preferences"

cat $configfile | jq '.signin.allowed = false' > $configfile.new
mv  $configfile.new $configfile

/usr/bin/chromium "$@"