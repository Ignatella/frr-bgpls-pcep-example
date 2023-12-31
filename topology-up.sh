#!/bin/sh

# ensure script exits on error
set -e

# links
sudo ip link add mv-main-start	 type dummy
sudo ip link add mv-top-start 	 type dummy
sudo ip link add mv-bottom-start type dummy
sudo ip link add mv-main-end 	 type dummy
sudo ip link add mv-top-end 	 type dummy
sudo ip link add mv-bottom-end 	 type dummy

# interface up
sudo ifconfig mv-main-start   up
sudo ifconfig mv-top-start    up
sudo ifconfig mv-bottom-start up
sudo ifconfig mv-main-end     up
sudo ifconfig mv-top-end      up
sudo ifconfig mv-bottom-end   up

# assign ip
sudo ifconfig mv-main-start   172.18.0.1/29
sudo ifconfig mv-top-start    172.18.2.1/29
sudo ifconfig mv-bottom-start 172.18.4.1/29
sudo ifconfig mv-main-end     172.18.1.1/29
sudo ifconfig mv-top-end      172.18.3.1/29
sudo ifconfig mv-bottom-end   172.18.5.1/29

# enable host mpls
sudo modprobe mpls_router
sudo modprobe mpls_gso
sudo modprobe mpls_iptunnel

sudo sysctl -w net.mpls.conf.mv-main-start.input=1
sudo sysctl -w net.mpls.conf.mv-top-start.input=1
sudo sysctl -w net.mpls.conf.mv-bottom-start.input=1
sudo sysctl -w net.mpls.conf.mv-main-end.input=1
sudo sysctl -w net.mpls.conf.mv-top-end.input=1
sudo sysctl -w net.mpls.conf.mv-bottom-end.input=1
sudo sysctl -w net.mpls.conf.lo.input=1
sudo sysctl -w net.mpls.platform_labels=1048575

# create docker network
docker network create -d macvlan --subnet=172.18.0.1/29 --gateway=172.18.0.1 -o parent=mv-main-start 			mv-main-start 
docker network create -d macvlan --subnet=172.18.2.1/29 --gateway=172.18.2.1 -o parent=mv-top-start 			mv-top-start 
docker network create -d macvlan --subnet=172.18.4.1/29 --gateway=172.18.4.1 -o parent=mv-bottom-start 			mv-bottom-start 
docker network create -d macvlan --subnet=172.18.1.1/29 --gateway=172.18.1.1 -o parent=mv-main-end 			mv-main-end 
docker network create -d macvlan --subnet=172.18.3.1/29 --gateway=172.18.3.1 -o parent=mv-top-end  			mv-top-end
docker network create -d macvlan --subnet=172.18.5.1/29 --gateway=172.18.5.1 -o parent=mv-bottom-end 			mv-bottom-end 
docker network create -d bridge  --subnet=172.20.0.1/29 --gateway=172.20.0.1 -o com.docker.network.bridge.name=br-odl	br-odl
docker network create -d bridge  --subnet=172.19.0.0/23 --gateway=172.19.1.1 -o com.docker.network.bridge.name=br-main	br-main

# containers
docker run --name main --hostname main \
-v ./vol/main/frr.log:/frr.log \
-v ./vol/main/conf:/etc/frr \
--privileged -d -it quay.io/frrouting/frr:9.0.1 

docker run --name start --hostname start \
-v ./vol/start/frr.log:/frr.log \
-v ./vol/start/conf:/etc/frr \
--privileged -d -it ignatella/frr:alpine-fa082128f9-odl-compatible-pcep

docker run --name top --hostname top \
-v ./vol/top/frr.log:/frr.log \
-v ./vol/top/conf:/etc/frr \
--privileged -d -it quay.io/frrouting/frr:9.0.1

docker run --name end --hostname end \
-v ./vol/end/frr.log:/frr.log \
-v ./vol/end/conf:/etc/frr \
--privileged -d -it ignatella/frr:alpine-fa082128f9-odl-compatible-pcep

docker run --name bottom --hostname bottom \
-v ./vol/bottom/frr.log:/frr.log \
-v ./vol/bottom/conf:/etc/frr \
--privileged -d -it quay.io/frrouting/frr:9.0.1

docker run --name bgp-ls-connector \
-d -it ignatella/bgp-ls-connector

docker run --name odl --hostname odl \
-v ./vol/odl/requests:/requests \
-v ./vol/odl/entrypoint.sh:/entrypoint.sh \
--entrypoint /entrypoint.sh -d -it opendaylight/odl

docker run --name pathman-sr \
-p 127.0.0.1:8020:8020 \
-d -it ignatella/pathman-sr

# disconnect container from default network
docker network disconnect bridge `docker ps -aqf "name=main$"`
docker network disconnect bridge `docker ps -aqf "name=start$"`
docker network disconnect bridge `docker ps -aqf "name=top$"`
docker network disconnect bridge `docker ps -aqf "name=end$"`
docker network disconnect bridge `docker ps -aqf "name=bottom$"`
docker network disconnect bridge `docker ps -aqf "name=bgp-ls-connector$"`

# connect to containers to network
docker network connect mv-main-start main
docker network connect mv-main-start start
docker network connect mv-top-start top
docker network connect mv-top-start start
docker network connect mv-bottom-start bottom
docker network connect mv-bottom-start start
docker network connect mv-main-end main
docker network connect mv-main-end end
docker network connect mv-top-end top
docker network connect mv-top-end end
docker network connect mv-bottom-end bottom
docker network connect mv-bottom-end end

docker network connect --ip 172.20.0.3 br-odl odl
docker network connect --ip 172.20.0.2 br-odl start
docker network connect --ip 172.20.0.4 br-odl end
docker network connect --ip 172.20.0.5 br-odl bgp-ls-connector
docker network connect --ip 172.20.0.6 br-odl pathman-sr

docker network connect --ip 172.19.0.1 br-main main
docker network connect --ip 172.19.0.2 br-main start
docker network connect --ip 172.19.0.3 br-main top
docker network connect --ip 172.19.0.4 br-main end
docker network connect --ip 172.19.0.5 br-main bottom
docker network connect --ip 172.19.1.2 br-main bgp-ls-connector

# enable mpls inside containers
docker exec main   sysctl -w net.mpls.conf.lo.input=1
docker exec main   sysctl -w net.mpls.conf.eth1.input=1
docker exec main   sysctl -w net.mpls.conf.eth2.input=1
docker exec main   sysctl -w net.mpls.platform_labels=1048575
docker exec main   sysctl -w net.ipv4.ip_forward=1

docker exec start  sysctl -w net.mpls.conf.lo.input=1
docker exec start  sysctl -w net.mpls.conf.eth1.input=1
docker exec start  sysctl -w net.mpls.conf.eth2.input=1
docker exec start  sysctl -w net.mpls.conf.eth3.input=1
docker exec start  sysctl -w net.mpls.platform_labels=1048575
docker exec start  sysctl -w net.ipv4.ip_forward=1

docker exec top    sysctl -w net.mpls.conf.lo.input=1
docker exec top    sysctl -w net.mpls.conf.eth1.input=1
docker exec top    sysctl -w net.mpls.conf.eth2.input=1
docker exec top    sysctl -w net.mpls.platform_labels=1048575
docker exec top    sysctl -w net.ipv4.ip_forward=1

docker exec end    sysctl -w net.mpls.conf.lo.input=1
docker exec end    sysctl -w net.mpls.conf.eth1.input=1
docker exec end    sysctl -w net.mpls.conf.eth2.input=1
docker exec end    sysctl -w net.mpls.conf.eth3.input=1
docker exec end    sysctl -w net.mpls.platform_labels=1048575
docker exec end    sysctl -w net.ipv4.ip_forward=1

docker exec bottom sysctl -w net.mpls.conf.lo.input=1
docker exec bottom sysctl -w net.mpls.conf.eth1.input=1
docker exec bottom sysctl -w net.mpls.conf.eth2.input=1
docker exec bottom sysctl -w net.mpls.platform_labels=1048575
docker exec bottom sysctl -w net.ipv4.ip_forward=1

# add lo0 interfaces
docker exec start ip link add lo0 type dummy
docker exec start ifconfig lo0 172.191.2.1 netmask 255.255.255.0 

docker exec end	  ip link add lo0 type dummy
docker exec end   ifconfig lo0 172.191.4.1 netmask 255.255.255.0 

# add routes frr admin prefixes 
docker exec --privileged odl ip route add 172.21.0.4/32 via 172.20.0.4
docker exec --privileged odl ip route add 172.21.0.2/32 via 172.20.0.2
