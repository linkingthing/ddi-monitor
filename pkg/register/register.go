package register

import (
	"context"
	"strings"
	"time"

	"github.com/zdnscloud/cement/log"
	"google.golang.org/grpc"

	"github.com/linkingthing/ddi-monitor/config"
	pb "github.com/linkingthing/ddi-monitor/pkg/proto"
	"github.com/linkingthing/ddi-monitor/pkg/util"
)

func Register(conn *grpc.ClientConn, conf *config.MonitorConfig) {
	cli := pb.NewMonitorManagerClient(conn)
	s := strings.Split(conf.ControllerAddr, ":")
	if len(s) < 1 {
		log.Warnf("can not get the ip of the controller")
		return
	}
	_, err := cli.Register(context.Background(), &pb.RegisterReq{IP: conf.Server.IP, HostName: conf.Server.HostName, Roles: util.GetPbRoles(conf.Server.Roles), ControllerIP: s[0]})
	for err != nil {
		if _, err = cli.Register(context.Background(), &pb.RegisterReq{IP: conf.Server.IP, HostName: conf.Server.HostName, Roles: util.GetPbRoles(conf.Server.Roles), ControllerIP: s[0]}); err != nil {
			log.Warnf("grpc client exec Register failed: %s", err.Error())
		}
		time.Sleep(time.Second * time.Duration(conf.Server.ProbeInterval))
	}
}
