#!/bin/bash

PLATFORM="raspberrypi4-64"
WORKDIR="$PWD/packageGoKrazy/$PLATFORM"
SRCPATH="$WORKDIR/../../build/tmp/deploy/images/$PLATFORM"
SRCIMG="$SRCPATH/hapi-image-$PLATFORM.rootfs.tar.bz2"
GOPACKAGEROOT="github.com/alf632/goynzane/packageGoKrazy/$PLATFORM"

sudo rm -r $WORKDIR/kernel
mkdir -p $WORKDIR/kernel
sudo rm -r $WORKDIR/firmware
mkdir -p $WORKDIR/firmware
sudo rm -rf $WORKDIR/rootfs
mkdir -p $WORKDIR/rootfs
sudo rm -rf $WORKDIR/rootfs-tmp
mkdir -p $WORKDIR/rootfs-tmp

#cd $WORKDIR/firmware
#tar -xvf $SRCIMG lib/firmware/brcm
#tar --delete 'lib/firmware/brcm*' -f $SRCIMG > deleted.tar
#mv deleted.tar $SRCIMG
#
#cd $WORKDIR/kernel
#tar -xvf $SRCIMG lib/modules
#tar -xvf $SRCIMG boot/*Image vmlinuz
#tar --delete 'lib/modules*' --delete 'boot/*Image' -f $SRCIMG > deleted.tar
#mv deleted.tar $SRCIMG
#
#cp $SRCIMG $WORKDIR/extrafiles_arm64.tar


cd $WORKDIR/rootfs-tmp
sudo tar --same-owner -xpf $SRCIMG

sudo rm etc/hostname
sudo rm etc/hosts

sudo rm etc/rc*.d/*20hostapd
sudo rm etc/rc*.d/*20dnsmasq
sudo rm etc/rc*.d/*30mosquitto

#sudo rm -r ../firmware/brcm
#sudo mv lib/firmware/brcm ../firmware
#sudo rm -r ../kernel/lib/modules
#mkdir ../kernel/lib
#sudo mv lib/modules ../kernel/lib
sudo rm ../kernel/vmlinuz
sudo cp boot/*Image ../kernel/vmlinuz
echo | sudo tee -a etc/fstab
echo "tmpfs /var/volatile tmpfs defaults 0 0" | sudo tee -a etc/fstab
sudo tar -cf extrafiles_arm64.tar *

mkdir -p ../rootfs/_gokrazy/
sudo mv extrafiles_arm64.tar ../rootfs/_gokrazy/
cd $WORKDIR/rootfs && go mod init $GOPACKAGEROOT/rootfs
cp $WORKDIR/../../static/launch.go $WORKDIR/rootfs/

cd $WORKDIR
cp -a $SRCPATH/{bcm2711-rpi-400.dtb,bcm2711-rpi-4-b.dtb,bcm2711-rpi-cm4.dtb,bcm2711-rpi-cm4s.dtb,overlay_map.dtb} kernel/
cp -a ../../static/{config.txt,cmdline.txt} kernel/
mkdir kernel/overlays
find $SRCPATH -depth -type f | grep -E vc4.*\.dtbo | xargs -I '{}' cp -a '{}' kernel/overlays/

cd $WORKDIR/kernel && go mod init $GOPACKAGEROOT/kernel
echo "package kernel" > $WORKDIR/kernel/kernel.go

#sudo rm -rf $WORKDIR/rootfs-tmp

cd $WORKDIR/..

#bmaptool copy $SRCIMG $WORKDIR/sdcard.img
#part1=$(fdisk -lu $WORKDIR/sdcard.img | sed -n "s:^$WORKDIR/sdcard.img1 \** *\([0-9]*\) .*$:\1:p")
#part2=$(fdisk -lu $WORKDIR/sdcard.img | sed -n "s:^$WORKDIR/sdcard.img2 \** *\([0-9]*\) .*$:\1:p")
#sudo mount -o loop,offset=$part1 $WORKDIR/sdcard.img $WORKDIR/kernel
#sudo mount -o loop,offset=$part2 $WORKDIR/sdcard.img $WORKDIR/rootfs
