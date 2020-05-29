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
	var req pb.RegisterReq
	req.IP = conf.Server.IP
	req.HostName = conf.Server.HostName
	req.IsSlave = false
	for _, v := range conf.Server.Roles {
		req.Roles = append(req.Roles, v)
	}
	if _, err := cli.Register(context.Background(), &req); err != nil {
		log.Errorf("grpc client exec Register failed: %s", err.Error())
		return err
	}
	return nil
}
