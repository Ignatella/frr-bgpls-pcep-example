sudo ip link delete mv-main-start 
sudo ip link delete mv-top-start 
sudo ip link delete mv-bottom-start 
sudo ip link delete mv-main-end 
sudo ip link delete mv-top-end 
sudo ip link delete mv-bottom-end 

docker rm -f main
docker rm -f start
docker rm -f top
docker rm -f end
docker rm -f bottom
docker rm -f odl
docker rm -f main-peer

docker network rm mv-main-start
docker network rm mv-top-start
docker network rm mv-bottom-start
docker network rm mv-main-end
docker network rm mv-top-end
docker network rm mv-bottom-end
docker network rm br-odl
docker network rm br-main
