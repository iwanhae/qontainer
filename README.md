# qontainer

**qontainer = QEMU + Container**

Run QEMU based virtual machine in a Container native way.

# With Docker

```bash
# With Ubuntu (default ID/PW is deploy/deploy)
$ docker run --rm -it\
    --device=/dev/kvm:/dev/kvm --device=/dev/net/tun:/dev/net/tun\
    --cap-add NET_ADMIN \
    -e VM_CPU=2\
    -e VM_MEMORY=2G\
    -e VM_DISK="https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64-disk-kvm.img" \
    -e VM_DISK_SIZE=25G\
    -e NETWORK_TYPE=bridge\
    -v $PWD:/data\
    ghcr.io/iwanhae/qontainer:v1.0.0

# With Amazon Linux 2 (default ID/PW is deploy/deploy)
$ docker run --rm -it\
    --device=/dev/kvm:/dev/kvm --device=/dev/net/tun:/dev/net/tun\
    --cap-add NET_ADMIN \
    -e VM_CPU=2\
    -e VM_MEMORY=2G\
    -e VM_DISK="https://cdn.amazonlinux.com/os-images/2.0.20231101.0/kvm/amzn2-kvm-2.0.20231101.0-x86_64.xfs.gpt.qcow2" \
    -e VM_DISK_SIZE=25G\
    -e NETWORK_TYPE=bridge\
    -v $PWD:/data\
    ghcr.io/iwanhae/qontainer:v1.0.0
```

# With Kubernetes

It is even working on Kubernetes with Cilium.
I believe you know what to do `¯\_(ツ)_/¯`

# Environment Variables

```bash
# CPU Cores
VM_CPU="2"
# Memory
VM_MEMORY="2G" 
# Default disk image path, if URL is provided, will downloads it and save to `/data/disk.img`
VM_DISK="/data/disk.img"
# Resize the disk to this. Can not shrink. Expand only.
VM_DISK_SIZE="25G"

# `user` > (default) Private IP Address will be provided via DHCP
# or 
# `bridge` > Pod IP (or Docker Container IP) will be provided via cloud-init.
NETWORK_TYPE="bridge" 
# Changing this not recommended
NETWORK_INTERFACE="eth0"
# Changing this not recommended; Will use container's IP if NETWORK_TYPE="bridge"
NETWORK_ADDRESS="172.17.0.2/16"
# Changing this not recommended; Will use container's default route if NETWORK_TYPE="bridge"
NETWORK_DEFAULT_GATEWAY="172.17.0.1"
# Changing this not recommended. Will use container's NSs if NETWORK_TYPE="bridge"
NETWORK_NAMESERVERS="10.43.0.10,8.8.8.8,8.8.4.4"
# Changing this not recommended. Will use container's Search if NETWORK_TYPE="bridge"
NETWORK_SEARCH="default.svc.k8s.iwanhae.kr.,svc.k8s.iwanhae.kr.,k8s.iwanhae.kr."
# Changing this not recommended. Will use randomly generated MAC addr if NETWORK_TYPE="bridge"
NETWORK_MAC_ADDRESS="52:de:4d:c0:72:09"

# default shell for the default user
GUEST_SHELL="/bin/bash"
# by default, will use container's hostname 
GUEST_HOSTNAME="c0a47fe410cb"
# name of default user
GUEST_USERNAME="deploy"
# encrypted passworf of default user
# default value is "deploy" (so the default ID and PW is "deploy:deploy", only accessible via console)
GUEST_PASSWORD="$6$rounds=4096$KUjo2cumnYaz0fmk$EsoVV1xP/FXIkv5mm4V26CR3qJrDZhs3Rga8OfBKNBUSsmCM7OHouHMHHz8lApGsD835DqpFvAgqJv1Hq5J.k0"
# default user's SSH authorized keys
GUEST_SSH_AUTHORIZED_KEYS=[]
# default user's SUDO policy
GUEST_SUDO="ALL=(ALL) NOPASSWD:ALL"
# base64 encoded user script that will be run at first boot.
# e.g., "ZWNobyBoZWxsbyB3b3JsZA=="
GUEST_USERSCRIPT_BASE64="" 

# Changing this not recommended
QEMU_EXECUTABLE="/usr/bin/qemu-system-x86_64"
```

