---
title: "Docker"
date: "2024-08-03"
tags: ["Programming"]
---

## Images
- An image is the underlying definition of what gets reconstituted into a running container. Every Linux container is based on an image.
- Each image consists of one or more linked filesystem layers that generally have a one-to-one mapping to each build step used to create the image.
- Docker relies on the Linux kernel providing drivers to run the storage backend which communicates with the underlying Linux filesystem to build and manage the layers that combine into a single usable image.
  - The storage backend provides a fast copy-on-write (CoW) system for image management.
- A Dockerfile describes all the steps that are required to create an image. Each line in a Dockerfile creates a new image layer containing all the changes that are a result of the command being issued.
- Docker only needs to build layers that deviate from previous builds and will reuse all the layers that have not changed.
  - When you rebuild an image, every single layer after the first introduced change will need to be rebuilt. 
- By default, Docker runs all processes as root within the container. Containers provide some isolation from the underlying operating system but they still run on the host kernel. Production containers should almost always be run in the context of an unprivileged user.
- It is not recommended to run commands like `apt-get -y update` in a Dockerfile as this crawls the repository index during each build with the consequence that the build is not guaranteed to be repeatable since package versions might change between builds.
- Instead, an application image can be based on another image that already has the required updates applied to it and where the versions are in a known state.
- Each instruction creates a new Docker image layer so it often makes sense to combine a few logically grouped commands into a single instruction. It is also possible to use `COPY` in combination with `RUN` to copy and execute a script.
- The order of commands in a Dockerfile can have a significant impact on build times. Try to order commands so that things that change between every build are closer to the bottom.
- It is generally considered a best practice to run only a single process within a container such that the container provides a single function and remains easy to individually horizontally scale.
- The `.dockerignore` file allows for excluding files and directories from an image.
- Base images are the lowest level images that other images will build upon. These are often based on minimal installs of Linux distributions.

```Dockerfile
# start from a community-maintained base image providing Node 18.13.0
FROM node:18.13.0

# build argument with a default value which is only available during the image build process
ARG email="anna@example.com"

# image metadata referencing a build argument
LABEL "maintainer"=$email
LABEL "rating"="Five Stars" "class"="First Class"

# user to run all processes as within the container (defaults to root)
USER root

# shell variable that can be used by the running application and during the build process
ENV AP="/data/app"
ENV SCPATH="/etc/supervisor/conf.d"

# run instructions to install dependencies and create directories
RUN apt-get -y update
RUN apt-get -y install supervisor
RUN mkdir -p /var/log/supervisor

# copy files from the local filesystem into the image
COPY ./supervisord/conf.d/* $SCPATH/
COPY *.js* $AP/

# change working directory in the image for remaining build instructions and the default process that launches with any resulting containers
WORKDIR $AP

RUN npm install

# define command that launches the process to be run within the container
CMD ["supervisord", "-n"]
```

```bash
# build and tag an image based on the files in the current directory
docker image build -t example/docker-node-hello:latest .

# disable local caching
docker image build --no-cache -t example/docker-node-hello:latest .

# override build argument
docker image build --build-arg email=me@example.com -t example/docker-node-hello:latest .

# inspect image metadata
docker image inspect example/docker-node-hello:latest 

# run container, mapping port 8080 in the container to port 8080 on the host
docker container run --rm -d -p 8080:8080 example/docker-node-hello:latest

# set environment variable
docker container run --rm -d -p 8080:8080 -e WHO="Sean and Karl" example/docker-node-hello:latest

# list running containers
docker container ls
docker ps

# format output using template string
docker container ls --format "table {{.ID}}\t{{.Image}}\t{{.Status}}"

# stop container
docker container stop 6b2fe6e3d205
docker stop 6b2fe6e3d205

# start bash shell in container
docker exec -it 6b2fe6e3d205 bash
```

## Registries
- A registry stores images.
- Deployment is the process of pulling an image from a repository and running it on one or more Docker hosts.
- Docker provides a centralized image registry called Docker Hub for private and public images. A Docker image registry can be hosted internally.
- Docker can store registry login information to use on your behalf.

```bash
# login to Docker Hub
docker login
cat ~/.docker/config.json

# login to specified registry
docker login docker.io

# logout
docker logout

# build image explicitly specifying image registry hostname and repository
docker image build -t docker.io/example/docker-node-hello:latest .

# edit image tags
docker image tag example/docker-node-hello:latest docker.io/reissenzhn/docker-node-hello:latest

# push image to registry
docker image push reissenzahn/docker-node-hello:latest

# pull image
docker image pull ${<myuser>}/docker-node-hello:latest

# search for images
docker search node
```

## Image Optimization
- For convenience, many Linux containers inherit from a base image that contains a minimal Linux distribution. However, this isn't required as containers only need to contain the files that are required to tun the application on the host kernel.
- It is often useful to troubleshooting to have access to a working shell in a container so images are often build from a very lightweight Linux distribution.

```bash
# export the files in a container as a tarball
docker container export ddc3f61f311b -o web-app.tar
tar -tvf web-app.tar

# run shell in a container, allocating a pseudo-TTY and keeping stdin open
docker run -it alpine:latest /bin/sh
```

- Multistage builds 
- This encourages doing builds insider Docker.
- The special image name `scratch` indicates to start from an empty image which includes no additional files.
- Multiple stages can be used and these stages do not need to even be related to one another.

```Dockerfile
# name image during the build phase
FROM docker.io/golang:alpine as builder

RUN apk update && \
  apk add git && \
  CGO_ENABLED=0 go install -a -ldflags '-s' \
  github.com/spkane/scratch-helloworld@latest

FROM scratch

# copy binary from build image into the current image
COPY --from=builder /go/bin/scratch-helloworld /helloworld

# document that this container listens on port 8080
EXPOSE 8080
CMD ["/helloworld"]
```

```bash
# build the image
docker image build .

# show image history and build steps
docker image history size1
```

- Filesystem layers that make up an image are strictly additive. While it is possible to mask files in previous layers, those files cannot be deleted. An image cannot be made smaller by simply deleting files that were generated in earlier steps.
- `docker image build --squash` can be used to squash multiple layer into a single layer. This will cause files that were deleted in the intermediate layers to disappear from the final image but it also means that the whole layer must be downloaded by every system that requires it.
- Earlier layers in an image cannot be made smaller by deleting files in subsequent layers. The only way to make a layer smaller is to remove files before saving the layer. This is commonly performed by stringing commands together on a single line with the && and \ operators.
- Package managers create large package index files.

```Dockerfile
FROM docker.io/fedora

# clean package index (does not decrease the size of the layer above!)
# RUN dnf install -y httpd
# RUN dnf clean all

RUN dnf install -y httpd && \
    dnf clean all

CMD ["/usr/sbin/httpd", "-DFOREGROUND"]


CMD ["/usr/sbin/httpd", "-DFOREGROUND"]
```

- Docker uses a layer cache to try to avoid rebuilding any image layers that it has already built and that do not contain any noticeable changes. As a result, the order of instructions in a Dockerfile can have a dramatic impact on how long a build takes on average.
- In general, the most stable and time-consuming portions of a build process should happen first while the code is added as late as possible.

```Dockerfile
FROM docker.io/fedora
RUN dnf install -y httpd && \
    dnf clean all
RUN mkdir -p /var/www && \
    mkdir -p /var/www/html
ADD index.html /var/www/html
CMD ["/usr/sbin/httpd", "-DFOREGROUND"]
```

- Directory caching allows for saving the contents of a directory inside an image in a special layer that can be bind-mounted at build time and then unmounted before the image snapshot is made.
- This is often used to handle directories where tools like package managers and language dependency manager download their databases and archive files.
- A caching directory will be removed from the resulting image and also be remounted in consecutive builds.

```Dockerfile
# syntax=docker/dockerfile:1
FROM python:3.9.15-slim-bullseye

RUN mkdir /app

WORKDIR /app

COPY . /app

# mount a caching layer into the container at /root/.cache for the duration of this build step
RUN --mount=type=cache,target=/root/.cache pip install -r requirements.txt

WORKDIR /app/mastermind

CMD ["python", "mastermind.py"]
```

## Troubleshooting Builds
- With BuildKit, none of the intermediate build layers are exported from the build container to the Docker daemon.
- One approach that works is to leverage multistage builds and the --target argument of docker image build.

```Dockerfile
# create multistage build
FROM node:18.13.0 as deploy

ARG email="anna@example.com"

LABEL "maintainer"=$email
LABEL "rating"="Five Stars" "class"="First Class"

USER root

ENV AP="/data/app"
ENV SCPATH="/etc/supervisor/conf.d"

RUN apt-get -y update
RUN apt-get -y install supervisor
RUN mkdir -p /var/log/supervisor

COPY ./supervisord/conf.d/* $SCPATH/
COPY *.js* $AP/

WORKDIR $AP

# just before failing command
FROM deploy
RUN npm installer

CMD ["supervisord", "-n"]
```


```bash
# target specific stage
docker image build -t example/docker-node-hello:debug --target deploy .

# and debug
docker container run --rm -ti docker.io/example/docker-node-hello:debug /bin/bash
```

## Multi-architecture Builds
- buildx is a Docker plugin.
- By default, buildx will leverage QEMU-based virtualization and binfmt_misc to support architectures that differ from the underlying system.
- BuildKit can utilize a build container when it builds images.
- When ENTRYPOINT is missing from the Dockerfile, the CMD instruction is expected to contain both the  process and all the required command-line arguments.
- Platform information is provided to the Docker server via an image manifest.

```bash
docker buildx version

# check qemu are properly registered
docker container run --rm --privileged multiarch/qemu-user-static --reset -p yes

# create a default buildx container called builder
docker buildx create --name builder --driver docker-container --use
docker buildx inspect --bootstrap

# build the image and side-load it into the local Docker server
docker buildx build --tag wordchain:test --load .

# build image for multiple architectures and leave the results in the build cache
docker buildx build --platform linux/amd64,linux/arm64 --tag wordchain:test .

# push resulting images
docker  buildx  build  --platform  linux/amd64,linux/arm64 --tag docker.io/spkane/wordchain:latest --push .

# show image manifest
docker manifest inspect docker.io/spkane/wordchain:latest
```

```Dockerfile
FROM golang:1.18-alpine3.15 AS build
RUN apk --no-cache add \
    bash \
    gcc \
    musl-dev \
    openssl
ENV CGO_ENABLED=0
COPY . /build
WORKDIR /build
RUN go install github.com/markbates/pkger/cmd/pkger@latest && \
    pkger -include /data/words.json && \
    go build .

FROM alpine:3.15 AS deploy
WORKDIR /
COPY --from=build /build/wordchain /
USER 500
EXPOSE 8080

# default process that is run by the container
ENTRYPOINT ["/wordchain"]

# command-line arguments passed to that process
CMD ["listen"]
```


## Containers
- A container is a self-contained execution environment that shares the kernel of the host system and is (optionally) isolated from other containers in the system.
- Container share a single kernel and isolation between workloads is implemented entirely within that one kernel.
- The `docker container run` command is a convenience command that wraps `docker container create` and `docker container start`.
- Create a container from an underlying image and execute the container.
- Many of the setting specified in the Dockerfile can be overridden when creating a container.
- When a container is started, Docker copies certain system files on the host into the configuration directory of the container on the host and then uses a bind mount to link that copy of the file into the container.
- If there is no ENTRYPOINT defined in the image then the final argument in the run command is the executable and command-line arguments to run in the container.
- If there is an ENTRYPOINT defined in the image, then the final argument is passed to the ENTRYPOINT process as a list of command-line arguments to that command.

```bash
# provide container name
docker container create --name="awesome-service" ubuntu:latest sleep 120

# execute container
docker container start awesome-service

# stop container
docker container stop awesome-service

# delete container
docker container rm awesome-service

# additional labels
docker container run --rm -d --name has-some-labels -l deployer=Joe -l tester=Asako ubuntu:latest sleep 1000

# filter components based on labels
docker container ls -a -f label=deployer=Joe

# inspect container
docker container inspect has-some-labels

# map network port in the underlying container to the host
-p 8080:8080

# pass environment variables into the container
-e

# delete the container when it exists
docker container run --rm container-name

# allocate a pseudo-TTY and keep STDIN open
docker container run -ti ubuntu:latest /bin/bash

# specify user
--user

# specify hostname
docker container run --rm -ti --hostname="mycontainer.example.com" ubuntu:latest /bin/bash

# override resolv.conf
docker container run --rm -ti --dns=8.8.8.8 --dns=8.8.4.4 --dns-search=example1.com --dns-search=example2.com ubuntu:latest /bin/bash
```

- Container name: Uniquely identifies the container.
- Labels: key/value pairs that can be applied to images and container as metadata. Container inherit all the labels from their parent image and additional labels can be added.
- Hostname: A bind mount is created for /etc/hostname which links to the hostname file Docker has prepared for the container, which contains the container ID by default.
- DNS Resolution: The resolv.conf file that configures DNS resolution is also managed via a bind mount. By default, this is an exact copy of the host resolv.conf file.
- MAC Address: By default, a container will receive a calculated MAC address starting with the 02:42:ac:11 prefix.

## Storage Volumes
- Storage that can persist between container deployments.
- Mounting storage from the Docker host is not generally advisable because it ties the container to a particular host for its persistent state. For cases like temporary cache files or other semi-ephemeral states, it can make sense.
  - Mount directories and individual files from the host server into the container. Fully qualified paths are required.
  - Volumes are mounted read-write by default.
  - Neither the host mount point nor the mount point in the container needs to preexist. If the host mount point does not exist already, then it will be created as a directory.
- It is possible to tell Docker that the root volume of a container should be mounted read-only. This prevents things like log files from filling up the allocated disk. When used in conjunction with a mounted volume, you can ensure that data is written only into expected locations.
- The tmpfs mount type allows for mounting a tmpfs filesystem into the container. This filesystem is completely in-memory and will be very fast, but also ephemeral and lost when the container is stopped.

```bash
# mount /mnt/session_data to /data within the container
docker container run --rm -ti --mount type=bind,target=/mnt/session_data,source=/data ubuntu:latest /bin/bash

# more concise
docker container run --rm -ti -v /mnt/session_data:/data:ro ubuntu:latest /bin/bash

# mount read-only
docker container run --rm -ti -v /mnt/session_data:/data:ro ubuntu:latest /bin/bash

# mount root volume as readonly
docker container run --rm -ti --read-only=true -v /mnt/session_data:/data ubuntu:latest /bin/bash

# 
docker container run --rm -ti --read-only=true --mount type=tmpfs,destination=/tmp,tmpfs-size=256M ubuntu:latest /bin/bash
```


## Resource Quotas
- Constraints are normally applied at the time of container creation. If  you  need  to  change  them,  you  can  use  the  docker container update command or deploy a new container with the adjustments.
- While  Docker  supports  various  resource  limits, you must have these capabilities enabled in your kernel for Docker to take advantage of them. You might need to add these as command-line parameters to your kernel on startup.

```bash
# show whether resource limits are supported
docker system info
```

- There are several ways to limit CPU usage:
  - CPU shares: The total computing power of all cores is considered to be the full pool of 1024 shares. This is a hint to the schedular about how long each container should be able to run each time it is scheduled.
    - The cgroup-based constraints on CPU shares are not hard limits; they are relative limits.
  - CPU pinning: Pin a container to one or more CPU cores. Work for this container will be scheduled only on the cores that have been assigned to this container. Additional CPU sharing restrictions on the container only take into account other containers running on the same set of cores.
  - The --cpus command can be set to a floating-point number between 0.01 and the number of CPU cores on the Docker server.
- While constraining the CPU only impacts the priority of an application for CPU time, the memory limit is a hard limit.

```bash
# cpu shares
docker container run --rm -ti --cpu-shares 512 spkane/train-os stress -v --cpu 2 --io 1 --vm 2 --vm-bytes 128M --timeout 120s

# cpu pinning
docker container run --rm -ti --cpu-shares 512 --cpuset-cpus=0-2 spkane/train-os stress -v --cpu 2 --io 1 --vm 2 --vm-bytes 128M --timeout 120s

# more convenient
docker container run --rm -ti --cpus=".25" spkane/train-os stress -v --cpu 2 --io 1 --vm 2 --vm-bytes 128M --timeout 60s

# dynamically adjust cpu allocation on two containers
docker container update --cpus="1.5" 092c5dc85044 92b797f12af1

# constrain container to 512m of RAM and 512m of swap
docker container run --rm -ti --memory 512m spkane/train-os stress -v --cpu 2 --io 1 --vm 2 --vm-bytes 128M --timeout 10s

# specify swap separately: 512m of memory and 256mb of addition swap space
docker container run --rm -ti --memory 512m --memory-swap=768m spkane/train-os stress -v --cpu 2 --io 1 --vm 2 --vm-bytes 128M  --timeout 10s

# show system events
docker system events
```

- Docker also supports limiting block I/O in a few different ways via the cgroups mechanism:
  - Apply prioritization to a container's use of block device I/O using the blkio.weight cgroup attribute. This can have a value of 0 (disabled) or a number between 10 and 1000 with the default being 500. To  set  this  weight  on  a  container,  you  need  to  pass  the  --blkio-weight  to  your docker container run  command  with  a  valid  value.  You  can  also  target  a  specific device using the --blkio-weight-device option.
  - As  with  CPU  shares,  tuning  the  weights  is  hard  to  get  right  in  practice,  but  we  can make  it  vastly  simpler  by  limiting  the  maximum  number  of  bytes  or  operations  per second  that  are  available  to  a  container  via  its  cgroup.  The  following  settings  let  us control that:
--device-read-bps     Limit read rate (bytes per second) from a device
--device-read-iops    Limit read rate (IO per second) from a device
--device-write-bps    Limit write rate (bytes per second) to a device
--device-write-iops   Limit write rate (IO per second) to a device

```bash
time docker container run -ti --rm --device-write-iops /dev/vda:256  spkane/train-os:latest bonnie++ -u 500:500 -d /tmp -r 1024 -s 2048 -x 1

time docker container run -ti --rm --device-write-bps /dev/vda:5mb  spkane/train-os:latest bonnie++ -u 500:500 -d /tmp -r 1024 -s 2048 -x 1
```

- The  following  code  is  a  list  of  the  types  of  system  resources  that  you  can  usually
constrain by setting soft and hard limits via the ulimit command:

```bash
ulimit -a
core file size (blocks, -c) 0
data seg size (kbytes, -d) unlimited
scheduling priority (-e) 0
file size (blocks, -f) unlimited
pending signals (-i) 5835
max locked memory (kbytes, -l) 64
max memory size (kbytes, -m) unlimited
open files (-n) 1024
pipe size (512 bytes, -p) 8
POSIX message queues (bytes, -q) 819200
real-time priority (-r) 0
stack size (kbytes, -s) 10240
cpu time (seconds, -t) unlimited
max user processes (-u) 1024
virtual memory (kbytes, -v) unlimited
file locks (-x) unlimited

docker container run --rm -d --ulimit nofile=150:300 nginx
```

## Start/Stop/Kill

- If  it  is  set  to  always,  the  container  will  restart  whenever  it exits,  with  no  regard  to  the  exit  code.  If  restart  is  set  to  on-failure,  whenever  the container  exits  with  a  nonzero  exit  code,  Docker  will  try  to  restart  the  container.  If we set restart to on-failure:3, Docker will try and restart the container three times before  giving  up.  unless-stopped  is  the  most  common  choice  and  will  restart  the container  unless  it  is  intentionally  stopped  with  something  like  docker container stop.
- When  stopped,  the  process  is  not  paused;  it  exits. Although  memory  and  temporary file  system  (tmpfs)  contents  will  have  been  lost,  all  of  the  contain‐ er’s other filesystem contents and metadata, including environment variables  and  port  bindings,  are  saved  and  will  still  be  in  place when you restart the container.
- On  reboot,  Docker  will  attempt  to  start  all  of  the  containers  that  were running  at  shutdown.
- In the previous docker container stop example, we’re sending the container a SIGTERM signal and waiting for the container to exit gracefully. Containers follow the same process group signal propagation that any other process group would receive on Linux.
-   normal  docker container stop  sends  a  SIGTERM  to  the  process.  If  you  want  to force  a  container  to  be  killed  if  it  hasn’t  stopped  after  a  certain  amount  of  time,  you can use the -t argument, like this. This tells Docker to initially send a SIGTERM signal as before, but if the container has not  stopped  within  25  seconds  (the  default  is  10),  it  tells  Docker  to  send  a  SIGKILL signal to forcefully kill it.

```bash
docker container create -p 6379:6379 redis:7.0

# list containers filtering by image
docker container ls -a --filter ancestor=redis:7.0

# start container
docker container start 092c5dc85044

# restart on failure three times
docker container run -ti --restart=on-failure:3 --memory 100m spkane/train-os stress -v --cpu 2 --io 1 --vm 2 --vm-bytes 128M --timeout 120s

# stop container
docker container stop 092c5dc85044

# list all containers
docker container ls -a

# terminate
docker container stop -t 25 092c5dc85044

# force kill
docker container kill 092c5dc85044
```

- There  are  a  few  reasons  why  we  might  not  want  to  completely  stop  our  container. We  might  want  to  pause  it,  leave  its  resources  allocated,  and  leave  its  entries  in  the process table.
- Pausing  leverages  the  cgroup  freezer,  which  essentially  just  prevents  your  process from  being  scheduled  until  you  unfreeze  it.  This  will  prevent  the  container  from doing anything while maintaining its overall state, including memory contents.

```bash
docker container pause 092c5dc85044

docker container unpause 092c5dc85044
```

## Cleaning up

```bash
docker container stop 092c5dc85044

docker container rm 092c5dc85044

# force remove running container
docker container rm -f 092c5dc85044

# list images
docker image ls

# delete image and all associated filesystem layers
docker image rm 0256c63af7db

# purge all images and container
docker system prune

# purge all unused images
docker system prune -a

# delete all container
docker container rm $(docker container ls -a -q)

# delete all images
docker image rm $(docker images -q)

# To remove all containers that exited with a nonzero state, you can use this filter:
$ docker container rm $(docker container ls -a -q --filter 'exited!=0')

# And to remove all untagged images, you can type this:
docker image rm $(docker images -q -f "dangling=true")
```

- We must stop all containers that are using an image before removing the image itself.


## Debugging








