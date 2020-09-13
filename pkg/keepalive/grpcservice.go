package keepalive

import (
	"context"
	"net"

	"github.com/linkingthing/ddi-monitor/config"
	pb "github.com/linkingthing/ddi-monitor/pkg/proto"
)

type DDIService struct {
	handler *DDIHandler
}

func NewDDIService(conf *config.MonitorConfig) *DDIService {
	return &DDIService{handler: newDDIHandler(conf.DNS.Addr, conf.DNS.ConfigDir, conf.DNS.ProxyPort)}
}

func (s *DDIService) StartDNS(ctx context.Context, req *pb.StartDNSRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.handler.startDNS(req); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}

func (s *DDIService) StartDHCP(ctx context.Context, req *pb.StartDHCPRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.handler.startDHCP(req); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}

func (s *DDIService) StopDNS(ctx context.Context, req *pb.StopDNSRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.handler.stopDNS(req); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}

func (s *DDIService) StopDHCP(ctx context.Context, req *pb.StopDHCPRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.handler.stopDHCP(req); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
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

func (s *DDIService) ReconfigDNS(ctx context.Context, req *pb.ReconfigDNSRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.handler.reconfigDNS(); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}

func (s *DDIService) ReloadDNSConfig(ctx context.Context, req *pb.ReloadDNSConfigRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.handler.reloadDNS(); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}

func (s *DDIService) AddDNSZone(ctx context.Context, req *pb.AddDNSZoneRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.handler.addDNSZone(req); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}

func (s *DDIService) UpdateDNSZone(ctx context.Context, req *pb.UpdateDNSZoneRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.handler.updateDNSZone(req); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}

func (s *DDIService) DeleteDNSZone(ctx context.Context, req *pb.DeleteDNSZoneRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.handler.deleteDNSZone(req); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}

func (s *DDIService) DumpDNSAllZonesConfig(ctx context.Context, req *pb.DumpDNSAllZonesConfigRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.handler.dumpDNSAllZonesConfig(); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}

func (s *DDIService) DumpDNSZoneConfig(ctx context.Context, req *pb.DumpDNSZoneConfigRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.handler.dumpDNSZoneConfig(req); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}

func (s *DDIService) ReloadNginxConfig(ctx context.Context, req *pb.ReloadNginxConfigRequest) (*pb.DDIMonitorResponse, error) {
	if err := s.handler.reloadNginxConfig(); err != nil {
		return &pb.DDIMonitorResponse{Succeed: false}, err
	} else {
		return &pb.DDIMonitorResponse{Succeed: true}, nil
	}
}
