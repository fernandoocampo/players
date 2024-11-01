package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	pb "github.com/fernandoocampo/players/pkg/pb/players"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServerSetup struct {
	Handler    *Handler
	Logger     *slog.Logger
	AppVersion string
	AppName    string
	Port       int
}

type Server struct {
	grpcServer *grpc.Server
	logger     *slog.Logger
	handler    *Handler
	appVersion string
	appName    string
	port       int
}

const resourceName = "grpc-server"

var errUnhealthy = errors.New("unhealthy")

func NewServer(setup ServerSetup) *Server {
	tracerInterceptor := makeTracerUnaryInterceptor(setup.AppVersion, setup.AppName)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(tracerInterceptor))
	newServer := Server{
		handler:    setup.Handler,
		logger:     setup.Logger,
		grpcServer: grpcServer,
		appVersion: setup.AppVersion,
		appName:    setup.AppName,
		port:       setup.Port,
	}

	return &newServer
}

// Start starts grpc server.
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.logger.Error("creating tcp listener for grpc server",
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("unable to start grpc server: %w", err)
	}

	pb.RegisterPlayerHandlerServer(s.grpcServer, s.handler)

	s.logger.Info("grpc server listening ready", slog.String("at", lis.Addr().String()))

	if err := s.grpcServer.Serve(lis); err != nil {
		s.logger.Error("serving with grpc server",
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("unable to serve with grpc server: %w", err)
	}

	return nil
}

func (s *Server) Close() error {
	s.grpcServer.GracefulStop()

	return nil
}

func (s *Server) Health() (string, error) {
	address := fmt.Sprintf(":%d", s.port)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.logger.Error("checking grpc server health", slog.String("error", err.Error()))

		return resourceName, fmt.Errorf("unable to check grpc server health: %w", err)
	}

	client := pb.NewPlayerHandlerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = client.SearchPlayers(ctx, &pb.SearchPlayersRequest{})
	if err != nil {
		return resourceName, errUnhealthy
	}

	return resourceName, nil
}
