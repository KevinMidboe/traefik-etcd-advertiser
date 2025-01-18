# traefik etcd advertiser

Used for advertising traefik config definitions from peer through etcd KV store.

## Configuration

```bash
wget https://github.com/kevinmidboe/traefik-etcd-advertiser/releases/LATEST > /usr/bin/local/traefik-etcd-advertiser
```

either create `.env` configuration file or prefix environmental variables required for etcd connection:

```bash
ETCD_ENDPOINTS="localhost:2379" traefik-etc-advertiser ...
```

## Usage

Pass either docker-compose or traefik dynamic config:

```bash
traefik-etcd-advertiser -filename web-service.yml -publish
```

