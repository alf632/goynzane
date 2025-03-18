#!/bin/sh

configfile="/perm/home/.config/chromium/Default/Preferences"

echo "disable chromium login"
cat $configfile | jq '.signin.allowed = false' > $configfile.new
mv  $configfile.new $configfile

