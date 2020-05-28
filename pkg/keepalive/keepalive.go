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
	time.Sleep(time.Second * time.Duration(conf.Server.ProbeInterval))
	cli := pb.NewMonitorManagerClient(conn)
	for GIsRunning {
		cpuUsage, memUsate, err := metric.New(conf)
		if err != nil {
			log.Errorf("get metric failed: %s", err.Error())
		}
		var target pb.KeepAliveReq
		target.IP = conf.Server.IP
		target.Role = conf.Server.Role
		target.CpuUsage = float32(*cpuUsage)
		target.MemUsage = float32(*memUsate)
		if Resp, err := cli.KeepAlive(context.Background(), &target); err != nil {
			log.Errorf("grpc client exec KeepAliveReq failed: %s,%s", Resp.Msg, err.Error())
		}
	}
	return nil
}
