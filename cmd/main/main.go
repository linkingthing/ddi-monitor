package main

import (
	"flag"
	"time"

	"github.com/zdnscloud/cement/log"

	"github.com/linkingthing/ddi-monitor/config"
	"github.com/linkingthing/ddi-monitor/pkg/grpcserver"
	"github.com/linkingthing/ddi-monitor/pkg/ha"
	"github.com/linkingthing/ddi-monitor/pkg/keepalive"
	"github.com/linkingthing/ddi-monitor/pkg/metric/exporter"
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

	err = keepalive.NewMonitorNode(conf)
	for err != nil {
		log.Errorf("register node failed: %s", err.Error())
		time.Sleep(3 * time.Second)
		err = keepalive.NewMonitorNode(conf)
	}

	grcpServer, err := grpcserver.New(conf)
	if err != nil {
		log.Fatalf("new grpc server failed: %s", err.Error())
	}

	go exporter.NodeExporter(conf)
	go ha.Server(conf)

	grcpServer.Run()
}
