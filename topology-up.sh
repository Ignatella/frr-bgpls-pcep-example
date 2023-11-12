#!/bin/sh


# links
sudo ip link add mv-main-start type dummy
sudo ip link add mv-top-start type dummy
sudo ip link add mv-bottom-start type dummy
sudo ip link add mv-main-end type dummy
sudo ip link add mv-top-end type dummy
sudo ip link add mv-bottom-end type dummy

# interface up
sudo ifconfig mv-main-start up
sudo ifconfig mv-top-start up
sudo ifconfig mv-bottom-start up
sudo ifconfig mv-main-end up
sudo ifconfig mv-top-end up
sudo ifconfig mv-bottom-end up

# assign ip
sudo ifconfig mv-main-start 172.18.0.1/29
sudo ifconfig mv-top-start 172.18.2.1/29
sudo ifconfig mv-bottom-start 172.18.4.1/29
sudo ifconfig mv-main-end 172.18.1.1/29
sudo ifconfig mv-top-end 172.18.3.1/29
sudo ifconfig mv-bottom-end 172.18.5.1/29

# enable mpls
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
docker network create -d macvlan --subnet=172.18.0.1/29 --gateway=172.18.0.1 -o parent=mv-main-start 	mv-main-start 
docker network create -d macvlan --subnet=172.18.2.1/29 --gateway=172.18.2.1 -o parent=mv-top-start 	mv-top-start 
docker network create -d macvlan --subnet=172.18.4.1/29 --gateway=172.18.4.1 -o parent=mv-bottom-start 	mv-bottom-start 
docker network create -d macvlan --subnet=172.18.1.1/29 --gateway=172.18.1.1 -o parent=mv-main-end 	mv-main-end 
docker network create -d macvlan --subnet=172.18.3.1/29 --gateway=172.18.3.1 -o parent=mv-top-end  	mv-top-end
docker network create -d macvlan --subnet=172.18.5.1/29 --gateway=172.18.5.1 -o parent=mv-bottom-end 	mv-bottom-end 

# containers
docker run --name main --hostname main --cap-add NET_ADMIN --cap-add SYS_ADMIN \
-v ./vol/main/sysctl.conf:/etc/sysctl.conf \
-v ./vol/main/frr.log:/frr.log \
-v ./vol/main/conf:/etc/frr \
--privileged -d -it quay.io/frrouting/frr:9.0.1 

docker run --name start --hostname start --cap-add NET_ADMIN --cap-add SYS_ADMIN \
-v ./vol/start/sysctl.conf:/etc/sysctl.conf \
-v ./vol/start/frr.log:/frr.log \
-v ./vol/start/conf:/etc/frr \
--privileged -d -it quay.io/frrouting/frr:9.0.1

docker run --name top --hostname top --cap-add NET_ADMIN --cap-add SYS_ADMIN \
-v ./vol/top/sysctl.conf:/etc/sysctl.conf \
-v ./vol/top/frr.log:/frr.log \
-v ./vol/top/conf:/etc/frr \
--privileged -d -it quay.io/frrouting/frr:9.0.1

docker run --name end --hostname end --cap-add NET_ADMIN --cap-add SYS_ADMIN \
-v ./vol/end/sysctl.conf:/etc/sysctl.conf \
-v ./vol/end/frr.log:/frr.log \
-v ./vol/end/conf:/etc/frr \
--privileged -d -it quay.io/frrouting/frr:9.0.1

docker run --name bottom --hostname bottom --cap-add NET_ADMIN --cap-add SYS_ADMIN \
-v ./vol/bottom/sysctl.conf:/etc/sysctl.conf \
-v ./vol/bottom/frr.log:/frr.log \
-v ./vol/bottom/conf:/etc/frr \
--privileged -d -it quay.io/frrouting/frr:9.0.1

# disconnect container from default network
docker network disconnect bridge `docker ps -aqf "name=main"`
docker network disconnect bridge `docker ps -aqf "name=start"`
docker network disconnect bridge `docker ps -aqf "name=top"`
docker network disconnect bridge `docker ps -aqf "name=end"`
docker network disconnect bridge `docker ps -aqf "name=bottom"`

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
