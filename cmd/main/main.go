package main

import (
	"flag"
	"fmt"

	"github.com/zdnscloud/cement/log"
	"google.golang.org/grpc"

	"github.com/linkingthing/ddi-monitor/config"
	"github.com/linkingthing/ddi-monitor/pkg/keepalive"
	"github.com/linkingthing/ddi-monitor/pkg/register"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "c", "ddi-monitor.conf", "configure file path")
	flag.Parse()

	log.InitLogger(log.Debug)
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("load config file failed: %s", err.Error())
	}

	fmt.Println("1111", conf)
	conn, err := grpc.Dial(conf.ControllerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial grpc server failed: %s", err.Error())
	}
	defer conn.Close()

	register.New(conn, conf)
	keepalive.New(conn, conf)
}
