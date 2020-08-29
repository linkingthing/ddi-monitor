package keepalive

import (
	"context"
	"path/filepath"

	"github.com/zdnscloud/cement/shell"

	"github.com/linkingthing/ddi-monitor/config"
	pb "github.com/linkingthing/ddi-monitor/pkg/proto"
)

const DNSConfName = "named.conf"

type DDIService struct {
	dnsConfDir string
}

func NewDDIService(conf *config.MonitorConfig) *DDIService {
	return &DDIService{dnsConfDir: conf.DNS.ConfigDir}
}

func (s *DDIService) StartDNS(ctx context.Context, req *pb.StartDNSRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.startDNS(req); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}

func (s *DDIService) startDNS(req *pb.StartDNSRequest) error {
	if isRunning, err := checkDNSIsRunning(); err != nil {
		return err
	} else if isRunning {
		return nil
	}

	if _, err := shell.Shell(filepath.Join(s.dnsConfDir, "named"), "-c", filepath.Join(s.dnsConfDir, DNSConfName)); err != nil {
		return err
	}

	return nil
}

func (s *DDIService) StartDHCP(ctx context.Context, req *pb.StartDHCPRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.startDHCP(req); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}

func (s *DDIService) startDHCP(req *pb.StartDHCPRequest) error {
	_, err := shell.Shell("keactrl", "start")
	return err
}

func (s *DDIService) StopDNS(ctx context.Context, req *pb.StopDNSRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.stopDNS(req); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}

func (s *DDIService) stopDNS(req *pb.StopDNSRequest) error {
	if _, err := shell.Shell(filepath.Join(s.dnsConfDir, "rndc"), "stop"); err != nil {
		return err
	}

	return nil
}

func (s *DDIService) StopDHCP(ctx context.Context, req *pb.StopDHCPRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.stopDHCP(req); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}

func (s *DDIService) stopDHCP(req *pb.StopDHCPRequest) error {
	if _, err := shell.Shell("keactrl", "stop"); err != nil {
		return err
	}

	return nil
}

func (s *DDIService) GetDNSState(context.Context, *pb.GetDNSStateRequest) (*pb.DDIStateResponse, error) {
	if isRunning, err := checkDNSIsRunning(); err != nil {
		return &pb.DDIStateResponse{}, err
	} else {
		return &pb.DDIStateResponse{IsRunning: isRunning}, nil
	}
}

func (s *DDIService) GetDHCPState(context.Context, *pb.GetDHCPStateRequest) (*pb.DDIStateResponse, error) {
	if isRunning, err := checkDHCPIsRunning(); err != nil {
		return &pb.DDIStateResponse{}, err
	} else {
		return &pb.DDIStateResponse{IsRunning: isRunning}, nil
	}
}
