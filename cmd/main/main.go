package main

import (
	"flag"

	"github.com/zdnscloud/cement/log"
	"google.golang.org/grpc"

	"github.com/linkingthing/ddi-monitor/config"
	"github.com/linkingthing/ddi-monitor/pkg/grpcserver"
	"github.com/linkingthing/ddi-monitor/pkg/keepalive"
	"github.com/linkingthing/ddi-monitor/pkg/metric/exporter"
	"github.com/linkingthing/ddi-monitor/pkg/register"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "c", "/usr/local/etc/ddi-monitor.conf", "configure file path")
	flag.Parse()

	log.InitLogger(log.Debug)
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("load config file failed: %s", err.Error())
	}

	conn, err := grpc.Dial(conf.ControllerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial controller grpc server failed: %s", err.Error())
	}
	defer conn.Close()

	register.Register(conn, conf)
	grcpServer, err := grpcserver.New(conf)
	if err != nil {
		log.Fatalf("new grpc server failed: %s", err.Error())
	}

	go exporter.NodeExporter(conf)
	go keepalive.Run(conn, conf)
	grcpServer.Run()
}
