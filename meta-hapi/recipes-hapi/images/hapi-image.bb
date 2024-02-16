SUMMARY = "ha-pi appliance"
DESCRIPTION = "config i use to run my home automation and HA Webinterface kiosk"
LICENSE = "MIT"

python do_display_banner() {
    bb.plain("***********************************************");
    bb.plain("*                                             *");
    bb.plain("*  Example recipe created by bitbake-layers   *");
    bb.plain("*                                             *");
    bb.plain("***********************************************");
}

addtask display_banner before do_build

IMAGE_FEATURES += "splash package-management x11 read-only-rootfs"

inherit core-image 
#features_check

REQUIRED_DISTRO_FEATURES = "x11 alsa"
QB_MEM = '${@bb.utils.contains("DISTRO_FEATURES", "opengl", "-m 512", "-m 256", d)}'
QB_MEM:qemuarmv5 = "-m 256"
QB_MEM:qemumips = "-m 256"

IMAGE_INSTALL:append = " \
kernel-modules linux-firmware-mediatek linux-firmware-bcm43455 \
packagegroup-core-x11-base \
x11vnc xf86-input-mouse xf86-input-evdev \
novnc \
alsa-lib alsa-utils alsa-plugins \
matchbox-wm chromium-x11 git \
podman podman-compose \
mosquitto dnsmasq hostapd dhcpcd wpa-supplicant iw \
htop jq procps \
usb-modeswitch \
wireguard-tools \
tcpdump iproute2-ss \
"

PACKAGECONFIG:append:chromium = " proprietary-codecs"
PACKAGECONFIG:append:mosquitto = " ssl websockets dns-srv dlt"
