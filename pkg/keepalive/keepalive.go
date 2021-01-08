package keepalive

import (
	"fmt"
	"net"
	"strings"

	"github.com/zdnscloud/cement/shell"

	"github.com/linkingthing/ddi-monitor/config"
)

func checkDHCPIsRunning() (bool, error) {
	ret, err := shell.Shell("ps", "-eaf")
	if err != nil {
		return false, fmt.Errorf("exec shell ps -eaf err:%s", err.Error())
	}
	if strings.Index(ret, "kea-dhcp6 -c ") > 0 &&
		strings.Index(ret, "kea-dhcp4 -c ") > 0 && strings.Index(ret, "kea-ctrl-agent -c ") > 0 {
		return true, nil
	}
	return false, nil
}

func checkDNSIsRunning() (bool, error) {
	ret, err := shell.Shell("ps", "-eaf")
	if err != nil {
		return false, fmt.Errorf("exec shell ps -eaf err:%s", err.Error())
	}
	if strings.Index(ret, "named -c") > 0 {
		return true, nil
	}
	return false, nil
}

func IsVIPOnLocal(vip string) (bool, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false, fmt.Errorf("InterfaceAddrs err:%s", err.Error())
	}
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok == false {
			continue
		}
		if ipnet.IP.To4() != nil && ipnet.IP.To4().String() == vip {
			return true, nil
		} else if ipnet.IP.To4() == nil && ipnet.IP.To16() != nil && ipnet.IP.To16().String() == vip {
			return true, nil
		}
	}

	return false, nil
}

func formatRoles(monitorRoles []config.ServiceRole) []string {
	var roles []string
	for _, pbRole := range monitorRoles {
		switch pbRole {
		case config.ServiceRoleController:
			roles = append(roles, string(config.ServiceRoleController))
		case config.ServiceRoleDNS:
			roles = append(roles, string(config.ServiceRoleDNS))
		case config.ServiceRoleDHCP:
			roles = append(roles, string(config.ServiceRoleDHCP))
		case config.ServiceRoleDataCenter:
			roles = append(roles, string(config.ServiceRoleDataCenter))
		}
	}
	return roles
}
