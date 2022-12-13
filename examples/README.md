# Piazza Shop: A Pizza Status Tracker
Piazza Shop demonstrates how Switchboard works in distributed ecosystems.

To run the demo, go to the `examples/docker-compose` directory and run:
```
docker compose up
```

The main components are:
- [Switchboard](#switchboard)
- [Frontend](#frontend)
- [Backend Services](#backend-services)
    - Store Service
    - Delivery Service
- [Authentication Service](#authentication-service)
- [Postgres](#postgres)
- [Redis](#redis)

## Switchboard
### Setup
Switchboard requires 2 items to run:
1. `config.yaml` file
2. `plugins` directory

#### Config
In this demo, Switchboard is configured to have a single endpoint: `/orders/:id/events`. This endpoint retrives information from the `/orders/:id/events` topic in Redis and gives it to the client. 

#### Plugins
Plugins are defined in `examples/docker-compose/plugins`. While there are 3 samples to experiment with (two middleware plugins and 1 enrichment plugin), the demo only uses `auth.so`, which is compiled from `auth.go`.

##### Compiling for Alpine
The demo plugins are compiled for Alpine. 

**If you are on a Linux machine**, you can use Go's built-in cross compilation infrastructure.
1. Install the `aarch64-linux-gnu-gcc` compiler 
```
sudo apt install gcc-aarch64-linux-gnu
```
2. Find the executable for that compiler, i.e. `/usr/bin/aarch-64-linux-gnu-gcc`
3. Run `CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=gcc-aarch64-linux-gnu go build --buildmode=plugin -o pluginName.so pluginName.go` where CC=[executable name, also available on PATH]

**If you are not on Linux**, there is a workaround where you can compile the plugin file via an alpine docker container.

1. Create a dockerfile with:
```
FROM golang:alpine
RUN apk add build-base bash
WORKDIR /src
```
2. Build `DOCKER_BUILDKIT=0 docker build . -t plugincompiler`  (DOCKER_BUILDKIT may or may not be necessary)
3. Run `docker run -v absolute/path/to/your/plugin/folder:/src -it plugincompiler /bin/bash`
4. This will open a shell for the docker container, then in that shell run: `go build --buildmode=plugin -o pluginName.so pluginName.go`

### Running
In `docker-compose.yaml`, Switchboard is added as a service. You can pull the latest image from [Dockerhub](https://hub.docker.com/r/jettcrowson/switchboard/tags). The `config.yaml` file and the entire `plugins` directory are mounted onto Switchboard as volumes.

## Frontend
A Flask app that uses templated HTML to render the user interface.

### Running
Visit `localhost:5000` after running docker compose to view.

## Backend Services
There are two services -- store status and delivery status. Both are written in Golang using Gin. 

### Running
The services will publish messages to the Redis topic `/orders/:id/events` when a GET request is made to `store-service/store/:id/events:54321` and `delivery-service/delivery/:id/events:12345`.

## Authentication Service
The auth service has 2 endpoints: `/login` and `/validate`.
- `POST /login`: returns the authentication token for valid username and password.
- `POST /validate`: verifies whether an auth token is valid

### Running
The service is available on port 8081.

## Postgres
This database is pulled from the latest `postgres` image on Dockerhub.

### Setup
The database is initialized with a table called `person` and is initialized with one user, `yich7110`. See the [initialization script](loginService/docker_postgres_init.sql) for more details.

## Redis
This database is pulled from the latest `redis` image on Dockerhub.
