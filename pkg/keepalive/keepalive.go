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

func New(conn *grpc.ClientConn, conf *config.MonitorConfig) error {
	cli := pb.NewMonitorManagerClient(conn)
	for {
		var err error
		var req pb.KeepAliveReq
		req.IP = conf.Server.IP
		req.Master = conf.Master
		req.Roles = util.GetPbRoles(conf.Server.Roles)
		if req.DnsAlive, err = checkDNSProcess(); err != nil {
			return fmt.Errorf("execute checkDNSProcess fail:%s", err.Error())
		}
		if req.DhcpAlive, err = checkDHCPProcess(); err != nil {
			return fmt.Errorf("execute checkDHCPProcess fail:%s", err.Error())
		}
		req.ControllerIP = strings.Split(conf.ControllerAddr, ":")[0]
		isLocalVip, err := isVIPOnLocal(conf.VIP)
		if err != nil {
			log.Warnf("isVIPOnLocal err:%s", err.Error())
			req.Vip = ""
		}
		if isLocalVip {
			req.Vip = conf.VIP
		}
		if _, err := cli.KeepAlive(context.Background(), &req); err != nil {
			log.Warnf("grpc client exec KeepAliveReq failed: %s", err.Error())
		}
		time.Sleep(time.Second * time.Duration(conf.Server.ProbeInterval))
	}
}

func checkDHCPProcess() (bool, error) {
	ret, err := shell.Shell("ps", "-eaf")
	if err != nil {
		return false, fmt.Errorf("exec shell ps -eaf err:%s", err.Error())
	}
	if strings.Index(ret, "kea-dhcp6 -c ") > 0 && strings.Index(ret, "kea-dhcp4 -c ") > 0 && strings.Index(ret, "kea-ctrl-agent -c ") > 0 {
		return true, nil
	}
	return false, nil
}

func checkDNSProcess() (bool, error) {
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
		} else {
			continue
		}
	}

	return false, nil
}
