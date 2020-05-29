package register

import (
	"context"

	"google.golang.org/grpc"

	"github.com/linkingthing/ddi-monitor/config"
	pb "github.com/linkingthing/ddi-monitor/pkg/proto"
	"github.com/zdnscloud/cement/log"
)

func New(conn *grpc.ClientConn, conf *config.MonitorConfig) error {
	cli := pb.NewMonitorManagerClient(conn)
	var target pb.RegisterReq
	target.IP = conf.Server.IP
	target.HostName = conf.Server.HostName
	for _, v := range conf.Server.Roles {
		target.Roles = append(target.Roles, v)
	}
	if _, err := cli.Register(context.Background(), &target); err != nil {
		log.Errorf("grpc client exec Register failed: %s", err.Error())
		return err
	}
	return nil
}
