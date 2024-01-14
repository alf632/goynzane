this is goynzane - what gokrazy never wanted to become


TL:DR
this a base oparating system that has basic wifi and container functionality focused on raspberry pi. in contrast to the commonly found approach to raspi projects - hack together some raspbian until it works and never update again - this OS is designed to decouple the base system from it's payload so the base image can be updated independen from the apps that are running on it.


Backround:
at some point it felt not only boring but outright painfull to start up the raspi imager again, wait the same ammount of time to download and install the image, clicking the same checkboxes, filling out the same forms just to get a general purpose operating system that needs additional packages installed and configured to even get started. so quickly firing up a raspi project meant spending the first few hours with setup.
with this frustration in mind i was verry open to the ideas of Michael Stapelberg with his gokrazy project.

GoKrazy:
"we had a crazy idea: what if we massively reduced the overall system complexity by getting rid of all software we don’t strictly need, and instead built up a minimal system from scratch entirely in Go, a memory safe programming language?" - gokrazy.org
gokrazy and perticularly it's utility program "gok" enable the engeneer to build a sd-card image and boot a raspberry pi from it. the image consists of a minimal set of components:
"""
Your app(s) + only 4 moving parts
 - the Linux kernel
   new versions are typically available < 24h after upstream release!
 - the Raspberry Pi firmware files
 - the Go compiler and standard library
 - the gokrazy system (minimal init system, updater, DHCP, NTP, …)
""" - gokrazy.org
the main goal is to build a statically linked go binary and package it together with a kernel and boilerplate.

GoYnzane:
while gokrazy is great at enabling you to deploy your golang binary and artifacts, it lacks support in areas where you would normally just ssh to the pi and use a pre installed tool or install a package and configure it. escpecially since most of the libs those tools and packages would require are not included in gokrazy instances. setting up a wifi client on a usb-wifi-interface together with a wifi ap on the build-in wifi chip including dhcp/dns server and uplink route setup didn't seem like a straight forward task to stay locked in the golang ecosystem.
fortunately gokrazy has a practical aproach to solve this:
"""
If your program needs extra files to be present in gokrazy’s root file system image at a specific location, you can add them with the extrafiles mechanism
""" - gokazy.org
to solve the dependencies for the wifi setup (described above) with a traditional stack (wpa_supplicant, hostapd, dnsmasq) there are a lot of extra files needed. finding and copying them is a tedious job. fortunately this job is called "making a linux distribution" and there are a lot of tools available to automate it.

Yocto:
"""
The Yocto Project
It’s not an embedded Linux distribution,
it creates a custom one for you.
""" - yoctoproject.org
yocto is a lot like gokrazy: it compiles a static operating system ready to be flashed onto a sd-card and booted on the raspberry pi.
only that it is not purposely build for the raspberry pi but supports a wide range of embedded devices and receives support from major hardware manufacturers. it's receipie book has all *most* the packages we would expect when searching for a tool and is further extensible.
so with yocto we can build a base set of extrafiles (binaries, libaries, configfiles...).
gokrazy then includes these in the instance's filesystem ready to be used.


Practical Example:
We want to realize a homeautomation system.
we need:
 - a broker
 - an automation framework
 - sensors
 - actors
 - connectivity between all components

we configure yocto to build us an image including everything nessesary to run the broker (mosquitto in most cases) and the automation framework (homeassistant/nodered) in containers. it also needs to include a wifi client if the raspi is not connected via ethernet and a wifi ap and dhcp-server if the sensors are supposed to directly connect to the pi (e.g. esp32-instances). <<build.sh>>
then we wait. a long time. a verry long time. proposals on the internet as to how to spend the waiting time range from drinking some hot bewerages to going on a weekend trip to paris. depending on your hardware. it took about 4h on my dell xps but i had trouble with the internet connection.

while we wait we write an init go program where we initialize the pi.
a most basic approach would be to call "exec.Command(("/etc/init.d/rcS", "")", since the default init system for yocto is sysVinit. this will bring the system to a point where mountpoints for the r/w overlay are wired-up, udev has detected external hardware and an x-server is started.

when yocto has finished compiling we package the result together with our init and some additional configuration (boot.txt, cmdline.txt...). <<package.sh>>

then we can add the newly generated rootfs and kernel packages to gokrazy and prepare the sd-card <<gok add>> <<gok overwrite --full=/dev/mmblk0>> <<gok update>>
when booting this sd-card we should be greeted with a terminal emulator running in an X window.

to make this base system do useful things we can now add additional gokrazy packages. like a launcher for the broker, wifi-setup, containers and possibly sensors running on the pi. 
