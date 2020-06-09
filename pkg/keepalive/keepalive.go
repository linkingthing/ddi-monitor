package keepalive

import (
	"context"
	"strings"
	"time"

	"github.com/zdnscloud/cement/log"
	"google.golang.org/grpc"

	"github.com/linkingthing/ddi-monitor/config"
	"github.com/linkingthing/ddi-monitor/pkg/metric/importer"
	pb "github.com/linkingthing/ddi-monitor/pkg/proto"
)

func New(conn *grpc.ClientConn, conf *config.MonitorConfig) error {
	cli := pb.NewMonitorManagerClient(conn)
	for {
		var req pb.KeepAliveReq
		cpuUsage, memUsage, err := importer.GetMetric(conf)
		if err != nil {
			log.Errorf("get metric from importer failed: %s", err.Error())
			req.CpuUsage = "0"
			req.MemUsage = "0"
		} else {
			req.CpuUsage = cpuUsage
			req.MemUsage = memUsage
		}
		req.IP = conf.Server.IP
		req.Roles = conf.Server.Roles
		req.DnsAlive = true
		req.DhcpAlive = true
		req.ControllerIP = strings.Split(conf.ControllerAddr, ":")[0]
		if _, err := cli.KeepAlive(context.Background(), &req); err != nil {
			log.Warnf("grpc client exec KeepAliveReq failed: %s", err.Error())
		}
		time.Sleep(time.Second * time.Duration(conf.Server.ProbeInterval))
	}
	return nil
}
