# Using Docker

- [Using Docker](#using-docker)
  - [Introduction](#introduction)
  - [Docker volumes](#docker-volumes)
  - [Known error messages when starting the pind container](#known-error-messages-when-starting-the-pind-container)
  - [Examples](#examples)
    - [Preamble](#preamble)
    - [Full node without RPC port](#full-node-without-rpc-port)
    - [Full node with RPC port](#full-node-with-rpc-port)
    - [Full node with RPC port running on TESTNET](#full-node-with-rpc-port-running-on-testnet)

## Introduction

With Docker you can easily set up *pind* to run your Bitcoin full node. You can find the official *pind* Docker images on Docker Hub [nyodeco/pind](https://hub.docker.com/r/nyodeco/pind). The Docker source file of this image is located at [Dockerfile](https://github.com/nyodeco/pind/blob/master/Dockerfile).

This documentation focuses on running Docker container with *docker-compose.yml* files. These files are better to read and you can use them as a template for your own use. For more information about Docker and Docker compose visit the official [Docker documentation](https://docs.docker.com/).

## Docker volumes

**Special diskspace hint**: The following examples are using a Docker managed volume. The volume is named *pind-data* This will use a lot of disk space, because it contains the full Bitcoin blockchain. Please make yourself familiar with [Docker volumes](https://docs.docker.com/storage/volumes/).

The *pind-data* volume will be reused, if you upgrade your *docker-compose.yml* file. Keep in mind, that it is not automatically removed by Docker, if you delete the pind container. If you don't need the volume anymore, please delete it manually with the command:

```bash
docker volume ls
docker volume rm pind-data
```

For binding a local folder to your *pind* container please read the [Docker documentation](https://docs.docker.com/). The preferred way is to use a Docker managed volume.

## Known error messages when starting the pind container

We pass all needed arguments to *pind* as command line parameters in our *docker-compose.yml* file. It doesn't make sense to create a *pind.conf* file. This would make things too complicated. Anyhow *pind* will complain with following log messages when starting. These messages can be ignored:

```bash
Error creating a default config file: open /sample-pind.conf: no such file or directory
...
[WRN] BTCD: open /root/.pind/pind.conf: no such file or directory
```

## Examples

### Preamble

All following examples uses some defaults:

- container_name: pind
  Name of the docker container that is be shown by e.g. ```docker ps -a```

- hostname: pind **(very important to set a fixed name before first start)**
  The internal hostname in the docker container. By default, docker is recreating the hostname every time you change the *docker-compose.yml* file. The default hostnames look like *ef00548d4fa5*. This is a problem when using the *pind* RPC port. The RPC port is using a certificate to validate the hostname. If the hostname changes you need to recreate the certificate. To avoid this, you should set a fixed hostname before the first start. This ensures, that the docker volume is created with a certificate with this hostname.

- restart: unless-stopped
  Starts the *pind* container when Docker starts, except that when the container is stopped (manually or otherwise), it is not restarted even after Docker restarts.

To use the following examples create an empty directory. In this directory create a file named *docker-compose.yml*, copy and paste the example into the *docker-compose.yml* file and run it.

```bash
mkdir ~/pind-docker
cd ~/pind-docker
touch docker-compose.yaml
nano docker-compose.yaml (use your favourite editor to edit the compose file)
docker-compose up (creates and starts a new pind container)
```

With the following commands you can control *docker-compose*:

```docker-compose up -d``` (creates and starts the container in background)

```docker-compose down``` (stops and delete the container. **The docker volume pind-data will not be deleted**)

```docker-compose stop``` (stops the container)

```docker-compose start``` (starts the container)

```docker ps -a``` (list all running and stopped container)

```docker volume ls``` (lists all docker volumes)

```docker logs pind``` (shows the log )

```docker-compose help``` (brings up some helpful information)

### Full node without RPC port

Let's start with an easy example. If you just want to create a full node without the need of using the RPC port, you can use the following example. This example will launch *pind* and exposes only the default p2p port 8333 to the outside world:

```yaml
version: "2"

services:
  pind:
    container_name: pind
    hostname: pind
    image: nyodeco/pind:latest
    restart: unless-stopped
    volumes:
      - pind-data:/root/.pind
    ports:
      - 8333:8333

volumes:
  pind-data:
```

### Full node with RPC port

To use the RPC port of *pind* you need to specify a *username* and a very strong *password*. If you want to connect to the RPC port from the internet, you need to expose port 8334(RPC) as well.

```yaml
version: "2"

services:
  pind:
    container_name: pind
    hostname: pind
    image: nyodeco/pind:latest
    restart: unless-stopped
    volumes:
      - pind-data:/root/.pind
    ports:
      - 8333:8333
      - 8334:8334
    command: [
        "--rpcuser=[CHOOSE_A_USERNAME]",
        "--rpcpass=[CREATE_A_VERY_HARD_PASSWORD]"
    ]

volumes:
  pind-data:
```

### Full node with RPC port running on TESTNET

To run a node on testnet, you need to provide the *--testnet* argument. The ports for testnet are 18333 (p2p) and 18334 (RPC):

```yaml
version: "2"

services:
  pind:
    container_name: pind
    hostname: pind
    image: nyodeco/pind:latest
    restart: unless-stopped
    volumes:
      - pind-data:/root/.pind
    ports:
      - 18333:18333
      - 18334:18334
    command: [
        "--testnet",
        "--rpcuser=[CHOOSE_A_USERNAME]",
        "--rpcpass=[CREATE_A_VERY_HARD_PASSWORD]"
    ]

volumes:
  pind-data:
```
