package register

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/linkingthing/ddi-monitor/config"
	pb "github.com/linkingthing/ddi-monitor/pkg/proto"
	"github.com/zdnscloud/cement/log"
)

func Register(conn *grpc.ClientConn, conf *config.MonitorConfig) {
	cli := pb.NewMonitorManagerClient(conn)
	var err error
	for err != nil {
		if _, err = cli.Register(context.Background(), &pb.RegisterReq{IP: conf.Server.IP, HostName: conf.Server.HostName, Roles: conf.Server.Roles, ControllerIP: strings.Split(conf.ControllerAddr, ":")[0]}); err != nil {
			log.Warnf("grpc client exec Register failed: %s", err.Error())
		}
		time.Sleep(time.Second * time.Duration(conf.Server.ProbeInterval))
	}
}
