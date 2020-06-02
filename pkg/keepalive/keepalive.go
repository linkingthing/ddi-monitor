package keepalive

import (
	"context"
	"time"

	"github.com/zdnscloud/cement/log"
	"google.golang.org/grpc"

	"github.com/linkingthing/ddi-monitor/config"
	"github.com/linkingthing/ddi-monitor/pkg/metric/importer"
	pb "github.com/linkingthing/ddi-monitor/pkg/proto"
)

var GIsRunning = true

func New(conn *grpc.ClientConn, conf *config.MonitorConfig) error {
	cli := pb.NewMonitorManagerClient(conn)
	for GIsRunning {
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
		if Resp, err := cli.KeepAlive(context.Background(), &req); err != nil {
			log.Errorf("grpc client exec KeepAliveReq failed: %s,%s", Resp.Msg, err.Error())
			continue
		}
		time.Sleep(time.Second * time.Duration(conf.Server.ProbeInterval))
	}
	return nil
}
