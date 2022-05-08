#构建本地镜像（模块3作业）

##1.编写 Dockerfile，将练习 2.2 编写的 httpserver 容器化（请思考有哪些最佳实践可以引入到 Dockerfile 中来）
##2.将镜像推送至 Docker 官方镜像仓库

```
echo "building httpserver library"
building httpserver library
mkdir -p bin/amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .
echo "building httpserver container"
building httpserver container
docker build -t kengaw/httpserver:v0.0.1 .
Sending build context to Docker daemon  7.419MB
Step 1/9 : FROM ubuntu
 ---> ba6acccedd29
Step 2/9 : ENV MY_SERVICE_PORT=80
 ---> Using cache
 ---> 3090d89b5a93
Step 3/9 : ENV MY_SERVICE_PORT1=80
 ---> Using cache
 ---> c0470455dffb
Step 4/9 : ENV MY_SERVICE_PORT2=80
 ---> Using cache
 ---> 99d8886abc34
Step 5/9 : ENV MY_SERVICE_PORT3=80
 ---> Using cache
 ---> c8792e3bd60e
Step 6/9 : LABEL multi.label1="value1" multi.label2="value2" other="value3"
 ---> Using cache
 ---> 033d3f628a55
Step 7/9 : ADD bin/amd64/httpserver /httpserver
 ---> bde911b270d1
Step 8/9 : EXPOSE 80
 ---> Running in 761c3367c61b
Removing intermediate container 761c3367c61b
 ---> c53d2987769a
Step 9/9 : ENTRYPOINT /httpserver
 ---> Running in e3d36d3b9ff5
Removing intermediate container e3d36d3b9ff5
 ---> e77fdf88df80
Successfully built e77fdf88df80
Successfully tagged kengaw/httpserver:v0.0.1
echo "pushing kengaw/httpserver"
pushing kengaw/httpserver
docker push kengaw/httpserver:v0.0.1
The push refers to repository [docker.io/kengaw/httpserver]
5e813ef2e7e8: Pushed 
9f54eef41275: Mounted from library/ubuntu 
v0.0.1: digest: sha256:928363872f36515443b42db22a31179b45569a3c4ea19867b7a6abef598ae8e2 size: 740 
```

创建容器，并将容器推至仓库，这里容器编号为，地址为[docker](https://hub.docker.com/layers/httpserver/kengaw/httpserver/v0.0.1/images/sha256-928363872f36515443b42db22a31179b45569a3c4ea19867b7a6abef598ae8e2?context=explore)

##3.通过 Docker 命令本地启动 httpserver
###3.1 通过docker images命令查看本地的镜像，是否存在镜像

```
REPOSITORY          TAG       IMAGE ID       CREATED              SIZE
kengaw/httpserver   v0.0.1    e77fdf88df80   About a minute ago   80.2MB
ubuntu              latest    ba6acccedd29   6 months ago         72.8MB
```

###3.2使用docker run命令启动镜像(-d 后台启动），注意启动名称
```
docker run -d kengaw/httpserver
Unable to find image 'kengaw/httpserver:latest' locally

docker run -d kengaw/httpserver:v0.0.1
0474d276993300907a2081ebf62cf569756b2b4ea1333df4e473222df8c1d9d0


docker ps:
CONTAINER ID   IMAGE                      COMMAND                  CREATED              STATUS              PORTS     NAMES
0474d2769933   kengaw/httpserver:v0.0.1   "/bin/sh -c /httpser…"   About a minute ago   Up About a minute   80/tcp    cool_ramanujan
```

##4.通过 nsenter 进入容器查看 IP 配置
###4.1 使用dockers inspect 命令去获取pid号，这里获取的pid号为21788
```
docker inspect 0474d2769933 | grep Pid
            "Pid": 21788,
            "PidMode": "",
            "PidsLimit": null,
```

###4.2 使用nsenter命令去查看docker的ip
```
nsenter -t 21788 -n ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
4: eth0@if5: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever
```
       
4.3 关闭容器
`docker stop 0474d2769933`

