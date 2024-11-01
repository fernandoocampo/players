# Design

I explain here the design philosophy that I am applying to players-api service to govern how I would extend and maintain this application.

Please notice that this is an evolutionary architecture design and it aims to evolve over time as new requirements, customer demands, tech debt and new technologies emerge.

## Application layout

I want to adopt and adapt a package oriented design as explained in the post [Package Oriented Design](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html) together with a typical Go layout (cmd/internal/pkg), my intent is to follow the language practices of Go (aka idiomatic Go) and show first the domain this application deals with rather than the technical layers (controllers/service/storage) that are common in other technologies.

You will find the following structure at the file system level as shown below.

```sh
.
├── Makefile
├── README.md
├── cmd
│   └── playersd
│       └── main.go
├── db-data
├── deploy
│   ├── Dockerfile
│   └── docker-compose.yaml
├── docs
│   ├── ASSUMPTIONS.md
│   ├── DESIGN.md
│   ├── FUTURE.md
│   └── ROADMAP.md
├── go.mod
├── go.sum
├── internal
│   ├── adapters
│   │   ├── cryptos
│   │   ├── grpc
│   │   ├── notifiers
│   │   └── storages
│   ├── appkit
│   │   ├── e2etests
│   │   └── unittests
│   ├── application
│   └── players
│
├── migrations
│   ├── 000001_create_schema.down.sql
│   └── 000001_create_schema.up.sql
└── pkg
    └── pb
        └── players
```

### cmd/

The folders under `cmd/` are always named for each program that will be built.

### cmd/playersd/

This is the folder that contains the main package and it is the one I am going to call to run the players-api service. It is a daemon that will run a gRPC server.

### internal/

Here I add packages that are internal to this project only, it means I don't want to share them with the outside world. This is because this service is not a library that other applications could import, so it does not make sense that other project outside `players-api` imports the `players` package. Go guarantees this at the compiler level.

### internal/application/

The `application` package provides the logic to start the player-api service.

To start this service you need to instantiate all the components and inject its required dependencies, this is a responsability of the `application` package. This package is the only that knows how to do that and gets the required parameters from env vars to inject into the other packages.

Since this is a small project, this package is also responsible for providing the logic to load all the configurations that the application needs to start. So it's important to make sure that no other package load settings or env vars by its own, if another package needs a parameter to work, `application` package should load such value from a specific source at startup time and inject the value into the package that requires the value.

Here you can provide logic to load values from environment variables or from configuration files.

### internal/players/

This package represents the `players` domain, so here I provide handlers (endpoints) and business logic related to handle players.

### internal/adapters/

This directory contains all the packages that provide logic for communicating with external resources a.k.a infrastructure. In other words these packages are a `translation` layer between the domain and a specific external technology. e.g. postgres database, redis cache, sqs queue, sns topic, other microservices, etc.

Here you can find two types of adapters: `ingoing` (aka. driving) and `outgoing` (aka. driven) adapters.

* `ingoing adapters` converts requests from a specific technology to requests the domain can understand.
* `outgoing adapters` converts calls from the domain into a specific technology. These adapters are the ones you are using mostly in the project. For instance: storage.

each adapter here must provide a factory method to instantiate an object that represents the connection to the external technology.

### internal/adapters/storages

Provides access to external storage mechanisms such as relational databases and cache systems (outgoing). Here you should provide methods to connect to these repositories and execute actions like create/update/delete. To achieve that, the package provides a function that allows us to create a client that connects to postgres.

### internal/adapters/grpc

Its responsibility is to provide logic to create grpc server and grpc handler (inbound adapter). Here I added logic to define the grpc handler which communicates with endpoints or services to expose player business logic to RPC clients.

### internal/adapters/notifiers

Its responsibility is to provide logic to create capabilities related to publish events into eventbus platforms.

### internal/adapters/appkit

provides utilities for tests only.

### pkg/pb/players

provide protobuffer artifacts for players service. Since these resources are under the `pkg` folder, they could be imported for other projects.