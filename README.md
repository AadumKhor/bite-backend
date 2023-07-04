# BiteSpeed Backend Task - Aayush Malhotra

## Overview

This repo contains the implementation of the [backend task](https://bitespeed.notion.site/Bitespeed-Backend-Task-Identity-Reconciliation-53392ab01fe149fab989422300423199) shared via email. The following tech stack is used: <br>

1. Gin framework in Golang
2. Docker
3. `docker-compose`
4. Postgres for DB
5. JSON for configs that can be exported to Vault (or similar tools)
6. `Bash` for scripting along with `Make`

I have not worked with NodeJS in depth, hence the choice of going with Golang.

## Folder structure explanation

<details>
  <summary>Folder Structure using `tree`</summary>
  
```bash
├── cmd
│   └── main.go
├── config
│   └── config.json
├── docker-compose.yaml
├── dockerfiles
│   └── Dockerfile
├── go.mod
├── go.sum
├── Makefile
├── makefiles
│   ├── banner.sh
│   ├── build.mk
│   ├── help.mk
│   └── local.mk
├── README.md
├── scripts
│   └── container_bite_build.sh
└── src
    ├── app
    │   ├── bite.go
    │   ├── handlers
    │   │   └── identify_handler.go
    │   └── middleware
    │       └── identify_middleware.go
    └── pkg
        ├── database
        │   ├── contact.go
        │   └── database.go
        ├── models
        │   ├── constants.go
        │   ├── contact.go
        │   └── identify.go
        └── utils
            ├── config.go
            ├── logger.go
            ├── string.go
            └── timezone.go

```
</details>

`cmd` -> Contains the `main.go` file which would be used to generate the binary <br>
`config` -> Contains the config for this project <br>
`dockerfiles` -> Self explanatory <br>
`makefiles` -> Using Make in this project to make it easier to use and replicate therefore some makefiles <br>
`scripts` -> Any random scripting that I have done for this project <br>
`src` -> Main application folder. Contains the `app` folder which contains our service code & `pkg` folder for peripheral services like db, logging, utils etc. <br>

## Local Setup

### Setup DB

Before starting the app, it is important to setup the DB. Run the following command to start just the DB for some pre-liminary operations. 

```bash
make build.db
```

This would start the DB in a detached state. You can use tools like DBeaver to connect to the DB using the username and password given in the config (see next section) <br> 

To make DB ready for operations, we need to create the table and possibly seed some data in it for testing. 

(Complete this section further)

### Starting `docker` containers

Copy the following config (if different) to `config/config.json`

```json
{
  "DefaultTimezone": "Asia/Kolkata",
  "Port": 8080,
  "Mode": "debug",
  "Databases": [
    {
      "Host": "db",
      "Port": 5432,
      "User": "dbuser",
      "Password": "BitespeedTask!",
      "Name": "bitespeed",
      "IdleConnections": 0,
      "OpenConnections": 50,
      "Type": "write",
      "SamplingRateInSeconds": 10
    },
    {
      "Host": "db",
      "Port": 5432,
      "User": "dbuser",
      "Password": "BitespeedTask!",
      "Name": "bitespeed",
      "IdleConnections": 0,
      "OpenConnections": 50,
      "Type": "read",
      "SamplingRateInSeconds": 10
    }
  ]
}
```

Run the following command

```bash
make build
```

This would internally call `docker compose` and run both the app and the db. Keep in mind, this command would not detach the containers hence you would be able to see both application and DB logs. Ideal state of logs should be this:

<details>
  <summary>Ideal State of terminal logs</summary>

```bash
❯ make build
Untagged: bitespeed-backend-task-app:latest
Deleted: sha256:e248469c86219cf1f8f729d521fd69e263e16134630b7113949b31c02c263b35
[+] Building 16.6s (17/17) FINISHED
 => [internal] load build definition from Dockerfile                                                                                                                                                                                                                                 0.0s
 => => transferring dockerfile: 989B                                                                                                                                                                                                                                                 0.0s
 => [internal] load .dockerignore                                                                                                                                                                                                                                                    0.0s
 => => transferring context: 2B                                                                                                                                                                                                                                                      0.0s
 => [internal] load metadata for docker.io/library/alpine:latest                                                                                                                                                                                                                     0.7s
 => [internal] load metadata for docker.io/library/golang:1.20-alpine                                                                                                                                                                                                                0.7s
 => [build 1/6] FROM docker.io/library/golang:1.20-alpine@sha256:fd9d9d7194ec40a9a6ae89fcaef3e47c47de7746dd5848ab5343695dbbd09f8c                                                                                                                                                    0.0s
 => [stage-1 1/5] FROM docker.io/library/alpine:latest@sha256:82d1e9d7ed48a7523bdebc18cf6290bdb97b82302a8a9c27d4fe885949ea94d1                                                                                                                                                       0.0s
 => [internal] load build context                                                                                                                                                                                                                                                    0.0s
 => => transferring context: 14.45kB                                                                                                                                                                                                                                                 0.0s
 => CACHED [build 2/6] RUN apk add --no-cache tzdata                                                                                                                                                                                                                                 0.0s
 => CACHED [build 3/6] WORKDIR /app                                                                                                                                                                                                                                                  0.0s
 => [build 4/6] COPY . /app                                                                                                                                                                                                                                                          0.0s
 => [build 5/6] RUN go mod download                                                                                                                                                                                                                                                  6.0s
 => [build 6/6] RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bite ./cmd/                                                                                                                                                                                           9.7s
 => CACHED [stage-1 2/5] RUN apk add --no-cache tzdata                                                                                                                                                                                                                               0.0s
 => CACHED [stage-1 3/5] WORKDIR /app/                                                                                                                                                                                                                                               0.0s
 => CACHED [stage-1 4/5] COPY --from=build /app/bite /app/bite                                                                                                                                                                                                                       0.0s
 => CACHED [stage-1 5/5] COPY --from=build /app/config /app/config                                                                                                                                                                                                                   0.0s
 => exporting to image                                                                                                                                                                                                                                                               0.0s
 => => exporting layers                                                                                                                                                                                                                                                              0.0s
 => => writing image sha256:e248469c86219cf1f8f729d521fd69e263e16134630b7113949b31c02c263b35                                                                                                                                                                                         0.0s
 => => naming to docker.io/library/bitespeed-backend-task-app                                                                                                                                                                                                                        0.0s
[+] Running 2/2
 ✔ Container bitespeed-backend-task-db-1   Recreated                                                                                                                                                                                                                                 0.0s
 ✔ Container bitespeed-backend-task-app-1  Recreated                                                                                                                                                                                                                                 0.0s
Attaching to bitespeed-backend-task-app-1, bitespeed-backend-task-db-1
bitespeed-backend-task-db-1   |
bitespeed-backend-task-db-1   | PostgreSQL Database directory appears to contain a database; Skipping initialization
bitespeed-backend-task-db-1   |
bitespeed-backend-task-db-1   | 2023-07-04 13:41:24.359 UTC [1] LOG:  starting PostgreSQL 15.3 (Debian 15.3-1.pgdg120+1) on x86_64-pc-linux-gnu, compiled by gcc (Debian 12.2.0-14) 12.2.0, 64-bit
bitespeed-backend-task-db-1   | 2023-07-04 13:41:24.360 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
bitespeed-backend-task-db-1   | 2023-07-04 13:41:24.360 UTC [1] LOG:  listening on IPv6 address "::", port 5432
bitespeed-backend-task-db-1   | 2023-07-04 13:41:24.362 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
bitespeed-backend-task-db-1   | 2023-07-04 13:41:24.365 UTC [29] LOG:  database system was shut down at 2023-07-04 13:41:05 UTC
bitespeed-backend-task-db-1   | 2023-07-04 13:41:24.369 UTC [1] LOG:  database system is ready to accept connections
bitespeed-backend-task-app-1  |
bitespeed-backend-task-app-1  |     _______  ___  _______  _______  _______  _______  _______  _______  ______
bitespeed-backend-task-app-1  |    |  _    ||   ||       ||       ||       ||       ||       ||       ||      |
bitespeed-backend-task-app-1  |    | |_|   ||   ||_     _||    ___||  _____||    _  ||    ___||    ___||  _    |
bitespeed-backend-task-app-1  |    |       ||   |  |   |  |   |___ | |_____ |   |_| ||   |___ |   |___ | | |   |
bitespeed-backend-task-app-1  |    |  _   | |   |  |   |  |    ___||_____  ||    ___||    ___||    ___|| |_|   |
bitespeed-backend-task-app-1  |    | |_|   ||   |  |   |  |   |___  _____| ||   |    |   |___ |   |___ |       |
bitespeed-backend-task-app-1  |    |_______||___|  |___|  |_______||_______||___|    |_______||_______||______|
bitespeed-backend-task-app-1  |
bitespeed-backend-task-app-1  |
bitespeed-backend-task-app-1  | GoVersion: go1.20.5
bitespeed-backend-task-app-1  | GOOS: linux
bitespeed-backend-task-app-1  | GOARCH: amd64
bitespeed-backend-task-app-1  | NumCPU: 12
bitespeed-backend-task-app-1  | GOROOT: /usr/local/go
bitespeed-backend-task-app-1  | Compiler: gc
bitespeed-backend-task-app-1  | Compiler: Tuesday, 4 Jul 2023
bitespeed-backend-task-app-1  | -----------------
bitespeed-backend-task-app-1  | [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.
bitespeed-backend-task-app-1  |
bitespeed-backend-task-app-1  | [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
bitespeed-backend-task-app-1  |  - using env:	export GIN_MODE=release
bitespeed-backend-task-app-1  |  - using code:	gin.SetMode(gin.ReleaseMode)
bitespeed-backend-task-app-1  |
bitespeed-backend-task-app-1  | [GIN-debug] POST   /identify                 --> github.com/AadumKhor/bitespeed-backend-task/src/app/handlers.IdentifyHandler.Handle-fm (4 handlers)
bitespeed-backend-task-app-1  | [GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
bitespeed-backend-task-app-1  | Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
bitespeed-backend-task-app-1  | [GIN-debug] Listening and serving HTTP on :8080

```

</details>

If it is any different, please reach out. It might be an issue that I have missed.



## Miscellaneous
