
```shell
docker build -t "base_event_service:v1" -f docker/base.dockerfile ./
docker build -t "event_service:v1" -f docker/user.dockerfile ./
docker run -t -d \
  --cpuset-cpus="0-35" \
  -m 220G \
  -v $(pwd)/data:/data \
  -p 17115:17115 \
  -p 17112:17112 \
  --name event_monitoring \
  event_service:v1
```


```shell
docker exec -it event_monitoring bash
```

```shell
mkdir /data/images
bash start_tor.sh
cd monitoring/
```

```shell
screen -S front
cd front
npm install
nano src/sagas/fetchData.js 
# change ip address of server variable on ip of backend, save and exit
npm start
```
Ctrl + A; Ctrl + D 

```shell
screen -S storage
cd storage
./data_storage
```
Ctrl + A; Ctrl + D 

```shell
screen -S backend
cd backend
./backend
```
Ctrl + A; Ctrl + D 