# `dsb`: DID-authentication based storage server

`dsb` is a PoC of storage server based on the [`didcomauth`](https://github.com/commercionetwork/didcomauth) library, providing
on-disk persistence of data sent by authorized peers.

## Usage and building

`dsb` provides the standard DID-based authentication endpoints, as well as:

 - a protected `/upload/{id}` endpoint, where authorized peers upload data via POST form - the `dsb` admin can 
 decide whether a DID can upload to a given id by means of the `cmd/dsb-add-didres` software,
 
 - an unprotected `/get/{id}` endpoint, where users can retrieve data contained at `id`.
 
The resource/DID mapping is held either in a Redis instance or in a in-memory map.
 
The standard procedure described in `didcomauth` applies to `dsb` too.

To build `dsb`, the standard Go procedure must be followed:

```shell script
$ go build -o dsb -ldflags="-s -w"
```

## Configuration

| Environment variable | Required | Default | Meaning |
|---|---|---|---|
|DSB_STORAGE_PATH | no | `./dsb` | path where `dsb` will save data |
|DSB_LOG_PATH     | no | `./dsb.log` | path where `dsb` will save its logs |
|DSB_DEBUG        | no | `false` | enables debug output|
|DSB_REDIS_ADDR   | no | `localhost:6379` | Redis host for caching and storage|
|DSB_CACHE_TYPE   | no | `0` | cache/storage type, `0` for Redis, `1` for in-memory map |
|DSB_JWT_SECRET   | yes | none | JWT secret used to sign `didcomauth` tokens |
|DSB_COMMERCIO_LCD| no | `http://localhost:1317` | commercio.network LCD REST server used for DDO resolution |