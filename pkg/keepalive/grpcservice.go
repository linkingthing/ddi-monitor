package keepalive

import (
	"context"
	"net"
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
	if isRunning, err := checkDHCPIsRunning(); err != nil {
		return err
	} else if isRunning {
		return nil
	}
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

func (s *DDIService) GetInterfaces(ctx context.Context, req *pb.GetInterfacesRequest) (*pb.GetInterfacesResponse, error) {
	if interfaces4, interfaces6, err := getInterfaces(); err != nil {
		return &pb.GetInterfacesResponse{}, err
	} else {
		return &pb.GetInterfacesResponse{Interfaces4: interfaces4, Interfaces6: interfaces6}, nil
	}
}

func getInterfaces() ([]string, []string, error) {
	interfaces4 := []string{"*"}
	interfaces6 := []string{"*"}
	its, err := net.Interfaces()
	if err != nil {
		return interfaces4, interfaces6, nil
	}

	for _, it := range its {
		addrs, err := it.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if ok == false {
				continue
			}

			ip := ipnet.IP
			if ip.To4() != nil {
				if ip.IsGlobalUnicast() {
					interfaces4 = append(interfaces4, it.Name+"/"+ip.String())
				}
			} else {
				if ip.IsGlobalUnicast() {
					interfaces6 = append(interfaces6, it.Name+"/"+ip.String())
				}
			}
		}
	}

	return interfaces4, interfaces6, nil
}
