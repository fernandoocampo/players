# Players

Players service is an API for all players in the solution. The service allows to manage players.

This service is implemented in Go.

## Design

see my design philosophy for this service in [document](docs/DESIGN.md)

## Assumptions

see [assumptions](docs/ASSUMPTIONS.md) document.

## Project Roadmap

see [roadmap](docs/ROADMAP.md) document.

## Future

see future next steps for this application [here](docs/FUTURE.md).

## How to test?

* unit tests
```sh
make test
```

* run e2e for storage (make sure to start database and migrations first)
```sh
make e2e-test-storage
```

* other individual e2e tests for storage

```sh
# test postgresql connection
make e2e-test-pg-connection
# test save player in db
make e2e-test-save-db
# test update player in db
make e2e-test-update-db
# test delete player from db
make e2e-test-delete-db
# test get player from db
make e2e-test-get-db
# test listing players with given email or nickname no resutls
make e2e-test-email-nickname-db
# test listing players with given email or nickname with results
make e2e-test-email-nickname-exist-db
# test listing players with given email or nickname ignoring player id
make e2e-test-email-nickname-ignore-db
# test listing players with search criteria.
make e2e-test-search-db
```

* other individual e2e tests for grpc server (make sure players api service is up)

```sh
# create player
make e2e-test-grpc-create
# update player
make e2e-test-grpc-update
# delete player
make e2e-test-grpc-delete
# search players
make e2e-test-grpc-search
```

## How to build?

* binary file

```sh
# linux
make build-linux
make build-linux-amd64
make build-linux-arm64
```

* image

```sh
make build-image
```

## How to start and stop local database?

* start database
```sh
make api-db-up
```

* stop database
```sh
make api-db-down
```

## How to run application?

### Run locally
make sure to create `.env` file based on `.env.example` file and set following environment variables:

```sh
CONTAINERTOOL=docker

PLAYERS_POSTGRES_PLAYER=playersdb
PLAYERS_POSTGRES_PASSWORD=playerspwd
PLAYERS_POSTGRES_DB=playersdb
PLAYERS_POSTGRES_PORT=5432
PLAYERS_POSTGRES_HOST=localhost

PLAYERS_LOG_LEVEL=development
PLAYERS_WEB_SERVER_PORT=8080
PLAYERS_GRPC_SERVER_PORT=50051
PLAYERS_PASSWORD_GENERATION_COST=4
```

`PLAYERS_LOG_LEVEL` could have 2 values: `development` or `production`

verify migrations are in place
```sh
make migration-up
```

then run

```sh
make run
```

### Run from containers

```sh
make api-up
```

once you finish, run 

```sh
make api-down
```

## How to run migrations?

* run migrations up
```sh
make migration-up
```

* run migrations down
```sh
make migration-down
```

## How to generate protobuffers?

```sh
make compile-proto
```