package keepalive

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/zdnscloud/cement/log"

	"github.com/linkingthing/ddi-monitor/config"
	"github.com/linkingthing/ddi-monitor/pkg/util"
)

type MonitorNode struct {
	ControllerAddr string       `json:"-"`
	Client         *http.Client `json:"-"`
	ID             string       `json:"id"`
	Roles          []string     `json:"roles"`
	HostName       string       `json:"hostName"`
	Master         string       `json:"master"`
	ControllerIp   string       `json:"controllerIp"`
	StartTime      time.Time    `json:"startTime"`
	DnsAlive       bool         `json:"dnsAlive"`
	DhcpAlive      bool         `json:"dhcpAlive"`
	Vip            string       `json:"vip"`
}

func NewMonitorNode(conf *config.MonitorConfig) error {
	monitorNode := &MonitorNode{
		ControllerAddr: conf.Controller.ApiIp + ":" + strconv.Itoa(conf.Controller.Port),
		Client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		ID:           conf.Server.IP,
		Roles:        formatRoles(conf.Server.Roles),
		HostName:     conf.Server.HostName,
		Master:       conf.Master,
		ControllerIp: conf.Controller.Ip,
		StartTime:    time.Now(),
		DnsAlive:     false,
		DhcpAlive:    false,
		Vip:          conf.VIP,
	}

	if err := monitorNode.registerNode(); err != nil {
		return err
	}
	go monitorNode.RunKeepalive(conf)

	return nil
}

func (monitorNode *MonitorNode) RunKeepalive(conf *config.MonitorConfig) {
	ticker := time.NewTicker(time.Duration(conf.Server.ProbeInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var err error
			if monitorNode.DnsAlive, err = checkDNSIsRunning(); err != nil {
				log.Warnf("check dns running failed:%s", err.Error())
				continue
			}

			if monitorNode.DhcpAlive, err = checkDHCPIsRunning(); err != nil {
				log.Warnf("check dhcp running failed:%s", err.Error())
				continue
			}

			if isLocalVip, err := IsVIPOnLocal(conf.VIP); err != nil {
				log.Warnf("IsVIPOnLocal err:%s", err.Error())
				monitorNode.Vip = ""
			} else if isLocalVip {
				monitorNode.Vip = conf.VIP
				monitorNode.Master = ""
			} else {
				monitorNode.Vip = ""
			}

			if err := monitorNode.keepaliveNode(); err != nil {
				log.Warnf("save keepAlive to db failed:%s", err.Error())
			}
		}
	}
}

func (monitorNode *MonitorNode) keepaliveNode() error {
	return monitorNode.notifyController(config.ActionKeepalive)
}

func (monitorNode *MonitorNode) registerNode() error {
	return monitorNode.notifyController(config.ActionRegister)
}

func (monitorNode *MonitorNode) notifyController(action string) error {
	if monitorNode.ControllerAddr == "" {
		return fmt.Errorf("controller addr is empty")
	}

	token, err := util.GetToken(monitorNode.Client, monitorNode.ControllerAddr)
	if err != nil {
		return err
	}

	return util.HttpRequest(monitorNode.Client, http.MethodPost,
		util.GenControllerRequestUrl(monitorNode.ControllerAddr, action, monitorNode.ID), &token, &monitorNode)
}
