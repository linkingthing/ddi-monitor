package keepalive

import (
	"fmt"
	"time"

	"github.com/zdnscloud/cement/log"
	restdb "github.com/zdnscloud/gorest/db"

	metrichandler "github.com/linkingthing/ddi-controller/pkg/metric/handler"
	"github.com/linkingthing/ddi-controller/pkg/metric/resource"
	"github.com/linkingthing/ddi-monitor/config"
	"github.com/linkingthing/ddi-monitor/pkg/db"
)

type MonitorNode struct {
	*resource.Node
}

func NewMonitorNode(conf *config.MonitorConfig) error {
	monitorNode := &MonitorNode{Node: &resource.Node{}}
	monitorNode.StartTime = time.Now()
	monitorNode.Ip = conf.Server.IP
	monitorNode.HostName = conf.Server.HostName
	monitorNode.Roles = formatRoles(conf.Server.Roles)
	monitorNode.ControllerIp = conf.ControllerIp
	monitorNode.Master = conf.Master
	monitorNode.SetID(conf.Server.IP)

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
			if monitorNode.DnsIsAlive, err = checkDNSIsRunning(); err != nil {
				log.Warnf("check dns running failed:%s", err.Error())
				continue
			}

			if monitorNode.DhcpIsAlive, err = checkDHCPIsRunning(); err != nil {
				log.Warnf("check dhcp running failed:%s", err.Error())
				continue
			}

			if isLocalVip, err := isVIPOnLocal(conf.VIP); err != nil {
				log.Warnf("isVIPOnLocal err:%s", err.Error())
				monitorNode.Vip = ""
			} else if isLocalVip {
				monitorNode.Vip = conf.VIP
			}

			if err := monitorNode.updateKeepAliveToDB(); err != nil {
				log.Warnf("save keepAlive to db failed:%s", err.Error())
			}
		}
	}
}

func (monitorNode *MonitorNode) updateKeepAliveToDB() error {
	return restdb.WithTx(db.GetDB(), func(tx restdb.Transaction) error {
		_, err := tx.Update(metrichandler.TableNode, map[string]interface{}{
			"roles":         monitorNode.Roles,
			"node_is_alive": true,
			"dns_is_alive":  monitorNode.DnsIsAlive,
			"dhcp_is_alive": monitorNode.DhcpIsAlive,
			"master":        monitorNode.Master,
			"controller_ip": monitorNode.ControllerIp,
			"start_time":    monitorNode.StartTime,
			"vip":           monitorNode.Vip,
		}, map[string]interface{}{restdb.IDField: monitorNode.ID})
		return err
	})
}

func (monitorNode *MonitorNode) registerNode() error {
	if err := restdb.WithTx(db.GetDB(), func(tx restdb.Transaction) error {
		if exists, err := tx.Exists(metrichandler.TableNode, map[string]interface{}{restdb.IDField: monitorNode.ID}); err != nil {
			return err
		} else if exists {
			_, err := tx.Update(metrichandler.TableNode, map[string]interface{}{
				"roles":         monitorNode.Roles,
				"host_name":     monitorNode.HostName,
				"master":        monitorNode.Master,
				"controller_ip": monitorNode.ControllerIp,
				"start_time":    monitorNode.StartTime},
				map[string]interface{}{restdb.IDField: monitorNode.ID})
			return err
		} else {
			_, err := tx.Insert(monitorNode.Node)
			return err
		}
	}); err != nil {
		return fmt.Errorf("register node %s failed: %s", monitorNode.HostName, err.Error())
	}

	return nil
}
