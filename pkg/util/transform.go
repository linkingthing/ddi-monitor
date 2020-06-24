package util

import (
	"github.com/linkingthing/ddi-monitor/config"
	pb "github.com/linkingthing/ddi-monitor/pkg/proto"
)

func GetPbRoles(monitorRoles []config.ServiceRole) []pb.ServiceRole {
	roles := []pb.ServiceRole{}
	for _, r := range monitorRoles {
		if r == config.ServiceRoleController {
			roles = append(roles, pb.ServiceRole_ServiceRoleController)
		} else if r == config.ServiceRoleDNS {
			roles = append(roles, pb.ServiceRole_ServiceRoleDNS)
		} else if r == config.ServiceRoleDHCP {
			roles = append(roles, pb.ServiceRole_ServiceRoleDHCP)
		}
	}
	return roles
}

func GetMonitorRoles(pbRoles []pb.ServiceRole) []config.ServiceRole {
	roles := []config.ServiceRole{}
	for _, r := range pbRoles {
		if r == pb.ServiceRole_ServiceRoleController {
			roles = append(roles, config.ServiceRoleController)
		} else if r == pb.ServiceRole_ServiceRoleDNS {
			roles = append(roles, config.ServiceRoleDNS)
		} else if r == pb.ServiceRole_ServiceRoleDHCP {
			roles = append(roles, config.ServiceRoleDHCP)
		}
	}
	return roles
}
