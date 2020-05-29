package keepalive

import (
	"context"
	"time"

	"github.com/zdnscloud/cement/log"
	"google.golang.org/grpc"

	"github.com/linkingthing/ddi-monitor/config"
	"github.com/linkingthing/ddi-monitor/pkg/metric"
	pb "github.com/linkingthing/ddi-monitor/pkg/proto"
)

var GIsRunning = true

func New(conn *grpc.ClientConn, conf *config.MonitorConfig) error {
	cli := pb.NewMonitorManagerClient(conn)
	for GIsRunning {
		cpuUsage, memUsage, err := metric.New(conf)
		if err != nil {
			log.Errorf("get metric failed: %s", err.Error())
			continue
		}
		var req pb.KeepAliveReq
		req.IP = conf.Server.IP
		req.Roles = conf.Server.Roles
		req.CpuUsage = *cpuUsage
		req.MemUsage = *memUsage
		req.DnsAlive = true
		req.DhcpAlive = true
		req.IsSlave = false
		if Resp, err := cli.KeepAlive(context.Background(), &req); err != nil {
			log.Errorf("grpc client exec KeepAliveReq failed: %s,%s", Resp.Msg, err.Error())
			continue
		}
		time.Sleep(time.Second * time.Duration(conf.Server.ProbeInterval))
	}
	return nil
}
