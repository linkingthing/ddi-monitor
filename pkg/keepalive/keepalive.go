package keepalive

import (
	"context"
	"fmt"
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
		var target pb.KeepAliveReq
		target.IP = conf.Server.IP
		target.Roles = conf.Server.Roles
		fmt.Println("cpu,mem:", cpuUsage, memUsage)
		target.CpuUsage = *cpuUsage
		target.MemUsage = *memUsage
		if Resp, err := cli.KeepAlive(context.Background(), &target); err != nil {
			log.Errorf("grpc client exec KeepAliveReq failed: %s,%s", Resp.Msg, err.Error())
			continue
		}
		time.Sleep(time.Second * time.Duration(conf.Server.ProbeInterval))
	}
	return nil
}
