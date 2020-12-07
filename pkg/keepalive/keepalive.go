package keepalive

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/zdnscloud/cement/log"
	"github.com/zdnscloud/cement/shell"
	"google.golang.org/grpc"

	"github.com/linkingthing/ddi-monitor/config"
	pb "github.com/linkingthing/ddi-monitor/pkg/proto"
	"github.com/linkingthing/ddi-monitor/pkg/util"
)

func Run(conn *grpc.ClientConn, conf *config.MonitorConfig) {
	cli := pb.NewMonitorManagerClient(conn)
	ticker := time.NewTicker(time.Duration(conf.Server.ProbeInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var err error
			req := pb.KeepAliveReq{
				IP:           conf.Server.IP,
				Master:       conf.Master,
				Roles:        util.GetPbRoles(conf.Server.Roles),
				ControllerIP: strings.Split(conf.ControllerAddr, ":")[0],
			}

			if req.DnsAlive, err = checkDNSIsRunning(); err != nil {
				log.Warnf("check dns running failed:%s", err.Error())
				continue
			}

			if req.DhcpAlive, err = checkDHCPIsRunning(); err != nil {
				log.Warnf("check dhcp running failed:%s", err.Error())
				continue
			}

			if isLocalVip, err := isVIPOnLocal(conf.VIP); err != nil {
				log.Warnf("isVIPOnLocal err:%s", err.Error())
				req.Vip = ""
			} else if isLocalVip {
				req.Vip = conf.VIP
			}

			if _, err := cli.KeepAlive(context.Background(), &req); err != nil {
				log.Warnf("grpc client exec KeepAliveReq failed: %s", err.Error())
				if conn_, err := grpc.Dial(conf.ControllerAddr, grpc.WithInsecure()); err != nil {
					log.Warnf("reDial controller grpc server failed: %s", err.Error())
					continue
				} else {
					conn.Close()
					conn = conn_
					cli = pb.NewMonitorManagerClient(conn)
				}
			}
		}
	}
}

func checkDHCPIsRunning() (bool, error) {
	ret, err := shell.Shell("ps", "-eaf")
	if err != nil {
		return false, fmt.Errorf("exec shell ps -eaf err:%s", err.Error())
	}
	if strings.Index(ret, "kea-dhcp6 -c ") > 0 && strings.Index(ret, "kea-dhcp4 -c ") > 0 && strings.Index(ret, "kea-ctrl-agent -c ") > 0 {
		return true, nil
	}
	return false, nil
}

func checkDNSIsRunning() (bool, error) {
	ret, err := shell.Shell("ps", "-eaf")
	if err != nil {
		return false, fmt.Errorf("exec shell ps -eaf err:%s", err.Error())
	}
	if strings.Index(ret, "named -c") > 0 {
		return true, nil
	}
	return false, nil
}

func isVIPOnLocal(vip string) (bool, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false, fmt.Errorf("InterfaceAddrs err:%s", err.Error())
	}
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok == false {
			continue
		}
		if ipnet.IP.To4() != nil && ipnet.IP.To4().String() == vip {
			return true, nil
		} else if ipnet.IP.To4() == nil && ipnet.IP.To16() != nil && ipnet.IP.To16().String() == vip {
			return true, nil
		}
	}

	return false, nil
}
