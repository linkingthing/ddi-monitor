package register

import (
	"context"

	"google.golang.org/grpc"

	"github.com/linkingthing/ddi-monitor/config"
	pb "github.com/linkingthing/ddi-monitor/pkg/proto"
	"github.com/zdnscloud/cement/log"
)

func Register(conn *grpc.ClientConn, conf *config.MonitorConfig) error {
	cli := pb.NewMonitorManagerClient(conn)
	if _, err := cli.Register(context.Background(), &pb.RegisterReq{IP: conf.Server.IP, HostName: conf.Server.HostName, Roles: conf.Server.Roles}); err != nil {
		log.Errorf("grpc client exec Register failed: %s", err.Error())
		return err
	}
	return nil
}
