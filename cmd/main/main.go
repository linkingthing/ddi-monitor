package main

import (
	"flag"

	"github.com/zdnscloud/cement/log"

	"github.com/linkingthing/ddi-controller/pkg/metric"
	"github.com/linkingthing/ddi-monitor/config"
	"github.com/linkingthing/ddi-monitor/pkg/db"
	"github.com/linkingthing/ddi-monitor/pkg/grpcserver"
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

	db.RegisterResources(metric.PersistentResources()...)
	if err := db.Init(conf); err != nil {
		log.Fatalf("init db failed: %s", err.Error())
	}

	if err := keepalive.NewMonitorNode(conf); err != nil {
		log.Fatalf("register node failed: %s", err.Error())
	}

	grcpServer, err := grpcserver.New(conf)
	if err != nil {
		log.Fatalf("new grpc server failed: %s", err.Error())
	}

	go exporter.NodeExporter(conf)
	grcpServer.Run()
}
