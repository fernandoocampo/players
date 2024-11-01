package application

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fernandoocampo/players/internal/adapters/cryptos"
	"github.com/fernandoocampo/players/internal/adapters/grpc"
	"github.com/fernandoocampo/players/internal/adapters/notifiers"
	"github.com/fernandoocampo/players/internal/adapters/storages"
	"github.com/fernandoocampo/players/internal/adapters/tracers"
	"github.com/fernandoocampo/players/internal/players"
)

// Closer defines behavior to close resources.
type Closer interface {
	Close() error
}

// HealthChecker defines behavior for health checks.
type HealthChecker interface {
	// Health checks resource health and return error if it is unhealthy
	Health() (resourceName string, err error)
}

// Event contains an application event.
type Event struct {
	Message string
	Error   error
}

type Application struct {
	settings          *settings
	dbClient          *sql.DB
	playerRepository  *storages.Storage
	playerService     *players.Service
	playerGRPCServer  *grpc.Server
	playerGRPCHandler *grpc.Handler
	passwordHasher    *cryptos.Bcrypt
	eventNotifier     *notifiers.Notifier
	tracerService     *tracers.TracerService
	logger            *slog.Logger
	resourcesToClose  []Closer
	resourcesHealth   []HealthChecker
	version           string
	buildDate         string
	commitHash        string
}

var (
	version      string
	buildDate    string
	commitHash   string
	errUnhealthy = errors.New("unhealthy")
)

// NewApplication instantiates a new service api application.
func NewApplication() *Application {
	newApplication := Application{
		version:    version,
		buildDate:  buildDate,
		commitHash: commitHash,
	}

	return &newApplication
}

// Run starts this application, loading settings and injecting dependencies.
func (a *Application) Run() error {
	a.printInfo()

	// load configuration
	a.loadConfiguration()
	// initialize logger
	a.initializeLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := a.initializeStorage()
	if err != nil {
		a.logger.Error("initializing player storage", slog.String("error", err.Error()))

		return fmt.Errorf("unable to start application: %w", err)
	}

	err = a.initializeTracerService(ctx)
	if err != nil {
		return fmt.Errorf("unable to start application: %w", err)
	}

	a.initializePasswordHasher()
	a.initializeNotifier()
	a.initializeService()
	a.initializeGRPCTransport()

	eventStream := make(chan Event)
	a.listenToOSSignal(ctx, eventStream)

	a.startEventNotifier(ctx)
	a.startGRPCServer(ctx, eventStream)

	a.addResourceToHealthChecks(a.eventNotifier)
	a.addResourceToHealthChecks(a.playerGRPCServer)
	a.addResourceToHealthChecks(a.playerRepository)

	// for health check only
	a.startWebServer(eventStream)

	a.addResourceToClose(a.playerGRPCServer)
	a.addResourceToClose(a.dbClient)
	a.addResourceToClose(a.tracerService)

	defer a.closeResources()

	appEvent := <-eventStream
	a.logger.Info("ending players api app", "event", appEvent.Message)

	if appEvent.Error != nil {
		a.logger.Error(
			"ending application with error",
			slog.String("error",
				appEvent.Error.Error()),
		)

		return fmt.Errorf("the application has terminated unexpectedly: %w", appEvent.Error)
	}

	return nil
}

func (a *Application) loadConfiguration() {
	slog.Info("loading configuration")

	a.settings = loadSettings()
}

func (a *Application) initializeLogger() {
	slog.Info("initializing logger")

	logLevel := slog.LevelInfo

	if a.settings.logLevel == developmentLog {
		logLevel = slog.LevelDebug
	}

	handlerOptions := &slog.HandlerOptions{
		Level: logLevel,
	}

	loggerHandler := slog.NewJSONHandler(os.Stdout, handlerOptions)
	logger := slog.New(loggerHandler)

	logger.Info("log settings", slog.String("level", handlerOptions.Level.Level().String()))

	slog.SetDefault(logger)

	a.logger = logger
}

func (a *Application) initializeStorage() error {
	a.logger.Info("initializing player storage")

	dbParameters := storages.Parameters{
		Host:     a.settings.repository.host,
		Player:   a.settings.repository.player,
		Password: a.settings.repository.password,
		DBName:   a.settings.repository.dbName,
		Port:     a.settings.repository.port,
	}

	dbClient, err := storages.NewPostgresClient(dbParameters)
	if err != nil {
		return fmt.Errorf("unable to create postgres client: %w", err)
	}

	a.dbClient = dbClient

	storageSetup := storages.StorageSetup{
		DB:     a.dbClient,
		Logger: a.logger,
	}

	a.playerRepository = storages.NewPlayerRepository(storageSetup)

	return nil
}

func (a *Application) initializePasswordHasher() {
	a.logger.Info("initializing password hasher")

	setup := cryptos.BcryptSetup{
		Logger: a.logger,
		Cost:   a.settings.passwordGenerationCost,
	}

	a.passwordHasher = cryptos.NewBcrypt(setup)
}

func (a *Application) initializeService() {
	a.logger.Info("initializing player service")

	setup := players.ServiceSetup{
		Storage:  a.playerRepository,
		Hasher:   a.passwordHasher,
		Notifier: a.eventNotifier,
		Logger:   a.logger,
	}

	a.playerService = players.NewService(&setup)
}

func (a *Application) initializeGRPCTransport() {
	a.logger.Info("initializing grpc transport")

	handlerSetup := grpc.HandlerSetup{
		Service: a.playerService,
		Logger:  a.logger,
	}

	a.playerGRPCHandler = grpc.NewHandler(handlerSetup)

	serverSetup := grpc.ServerSetup{
		Handler:    a.playerGRPCHandler,
		Logger:     a.logger,
		AppVersion: a.version,
		AppName:    serviceName,
		Port:       a.settings.grpcServerPort,
	}

	a.playerGRPCServer = grpc.NewServer(serverSetup)
}

func (a *Application) initializeNotifier() {
	a.logger.Info("initializing player events notifier")

	setup := notifiers.NotifierSetup{
		Logger:           a.logger,
		TimeoutToPublish: a.settings.timeoutToPublishSec,
	}

	a.eventNotifier = notifiers.NewNotifier(setup)
}

func (a *Application) initializeTracerService(ctx context.Context) error {
	a.logger.Info("starting tracer service")

	exporterSetup := tracers.ExporterSetup{
		TraceCollectorUrl: a.settings.tracerServiceURL,
	}

	exporter, err := tracers.CreateOTLPExporter(ctx, &exporterSetup)
	if err != nil {
		return fmt.Errorf("unable to initialize tracer service: %w", err)
	}

	tracerSetup := tracers.TracerSetup{
		Exporter:    exporter,
		ServiceName: serviceName,
	}

	a.tracerService = tracers.NewTracerService(&tracerSetup)

	return nil
}

// startGRPCServer starts the grpc server.
func (a *Application) startGRPCServer(ctx context.Context, eventStream chan<- Event) {
	go func() {
		a.logger.Info("starting grpc server",
			slog.Int("port", a.settings.grpcServerPort))

		err := a.playerGRPCServer.Start()
		if err != nil {
			a.logger.Error("strting grpc server", slog.String("error", err.Error()))
			select {
			case <-ctx.Done():
				return
			case eventStream <- newEvent("grpc server was ended with error", err):
				return
			}
		}

		select {
		case <-ctx.Done():
			return
		case eventStream <- newEvent("grpc server has ended", nil):
			return
		}
	}()
}

// startWebServer starts the web server.
func (a *Application) startWebServer(eventStream chan<- Event) {
	go func() {
		a.logger.Info("starting http server for health check",
			slog.Int("port", a.settings.webServerPort))

		mux := http.NewServeMux()
		mux.HandleFunc("/healthz", func(res http.ResponseWriter, _ *http.Request) {
			report, err := a.checkHealth()
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
			}

			_, err = res.Write(report)
			if err != nil {
				a.logger.Error("writing http response after health check", slog.String("error", err.Error()))
			}
		})

		mux.HandleFunc("/readyz", func(res http.ResponseWriter, _ *http.Request) {
			res.WriteHeader(http.StatusOK)
		})

		srv := http.Server{
			Addr:    fmt.Sprintf(":%d", a.settings.webServerPort),
			Handler: mux,
		}

		err := srv.ListenAndServe()
		if err != nil {
			eventStream <- Event{
				Message: "web server was ended with error",
				Error:   err,
			}

			return
		}

		eventStream <- Event{
			Message: "web server was ended",
		}
	}()
}

func (a *Application) startEventNotifier(ctx context.Context) {
	a.logger.Info("starting event notifier worker")

	a.eventNotifier.Start(ctx)
}

func (a *Application) listenToOSSignal(ctx context.Context, eventStream chan<- Event) {
	signalStream := make(chan os.Signal, 1)
	signal.Notify(signalStream, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		a.logger.Info("starting listener for os signals")

		defer close(signalStream)

		select {
		case <-ctx.Done():
			return
		case osSignal, ok := <-signalStream:
			if !ok {
				return
			}

			event := Event{
				Message: osSignal.String(),
			}

			select {
			case <-ctx.Done():
				return
			case eventStream <- event:
			}
		}
	}()
}

func (a *Application) printInfo() {
	slog.Info(
		"starting service",
		slog.String("version", a.version),
		slog.String("commit", a.commitHash),
		slog.String("build-date", a.buildDate),
	)
}

func (a *Application) Info() string {
	return fmt.Sprintf(
		"{ %q: %q, %q: %q, %q: %q }",
		"version", a.version,
		"commit", a.commitHash,
		"build", a.buildDate,
	)
}

func (a *Application) addResourceToClose(closer Closer) {
	a.resourcesToClose = append(a.resourcesToClose, closer)
}

func (a *Application) addResourceToHealthChecks(healthChecker HealthChecker) {
	a.resourcesHealth = append(a.resourcesHealth, healthChecker)
}

func (a *Application) closeResources() {
	a.logger.Info("closing application resources", slog.Int("count", len(a.resourcesToClose)))

	for _, resource := range a.resourcesToClose {
		err := resource.Close()
		if err != nil {
			a.logger.Error("unable to close resource", slog.String("error", err.Error()))
		}
	}
}

func (a *Application) checkHealth() ([]byte, error) {
	type healthReport struct {
		Info      string   `json:"info"`
		Resources []string `json:"resources"`
	}

	var result healthReport

	result.Info = fmt.Sprintf(
		"%s: %s, %s: %s, %s: %s",
		"version", a.version,
		"commit", a.commitHash,
		"build", a.buildDate,
	)

	allGood := true

	for _, resource := range a.resourcesHealth {
		message := "ok"

		name, healthErr := resource.Health()
		if healthErr != nil {
			allGood = false
			message = healthErr.Error()
		}

		result.Resources = append(result.Resources, name, message)
	}

	var err error
	if !allGood {
		err = errUnhealthy
	}

	report, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		a.logger.Error("marshalling health report", slog.String("error", marshalErr.Error()))

		return nil, marshalErr
	}

	return report, err
}

func newEvent(message string, err error) Event {
	return Event{
		Message: message,
		Error:   err,
	}
}
