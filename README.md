# traefik etcd advertiser

Used for advertising traefik config definitions from peer through etcd KV store.

[![Build Status](https://drone.schleppe.cloud/api/badges/KevinMidboe/traefik-etcd-advertiser/status.svg)](https://drone.schleppe.cloud/KevinMidboe/traefik-etcd-advertiser)

# Install

Install replacing `OS` and `ARCH`, e.g. `linux-amd64`. See [all releases](https://github.com/KevinMidboe/traefik-etcd-advertiser/releases/latest):

```bash
curl -s -L https://github.com/KevinMidboe/traefik-etcd-advertiser/releases/latest/download/traefik-etcd-advertiser-OS-ARCH.tar | tar xvz

# verify install
./traefik-etcd-advertiser-linux-amd64 -version
```

## Configuration

either create `.env` configuration file or prefix environmental variables required for etcd connection:

```bash
ETCD_ENDPOINTS="localhost:2379" traefik-etc-advertiser ...
```

## Usage

Pass either docker-compose or traefik dynamic config:

```bash
traefik-etcd-advertiser -filename web-service.yml -publish
```

