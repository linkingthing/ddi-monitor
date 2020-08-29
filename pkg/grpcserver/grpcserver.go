package grpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/linkingthing/ddi-monitor/config"
	"github.com/linkingthing/ddi-monitor/pkg/keepalive"
	pb "github.com/linkingthing/ddi-monitor/pkg/proto"
)

type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
}

func New(conf *config.MonitorConfig) (*GRPCServer, error) {
	listener, err := net.Listen("tcp", conf.Server.GrpcAddr)
	if err != nil {
		return nil, fmt.Errorf("create listener with addr %s failed: %s", conf.Server.GrpcAddr, err.Error())
	}

	grpcServer := &GRPCServer{
		server:   grpc.NewServer(),
		listener: listener,
	}

	pb.RegisterDDIMonitorServer(grpcServer.server, keepalive.NewDDIService(conf))
	return grpcServer, nil
}

func (s *GRPCServer) Run() error {
	defer s.Stop()
	return s.server.Serve(s.listener)
}

func (s *GRPCServer) Stop() error {
	s.server.GracefulStop()
	return nil
}
