package config

import (
	"github.com/zdnscloud/cement/configure"
)

type ServiceRole string

const (
	ServiceRoleDHCP       ServiceRole = "dhcp"
	ServiceRoleDNS        ServiceRole = "dns"
	ServiceRoleController ServiceRole = "controller"

	ActionStartHa    string = "start_ha"
	ActionMasterUp   string = "master_up"
	ActionMasterDown string = "master_down"
	ActionRegister          = "register"
	ActionKeepalive         = "keepalive"
)

const (
	AuthKey  = "authorization"
	Username = "systemapi"
	Password = "systemapi"
)

type MonitorConfig struct {
	Path       string         `yaml:"-"`
	Server     ServerConf     `yaml:"server"`
	Controller ControllerConf `yaml:"controller"`
	Prometheus PrometheusConf `yaml:"prometheus"`
	Master     string         `yaml:"master"`
	VIP        string         `yaml:"vip"`
	DNS        DNSConf        `yaml:"dns"`
	PgHaCliDir string         `yaml:"pgha_cli_dir"`
}

type ControllerConf struct {
	Ip    string `yaml:"ip"`
	ApiIp string `yaml:"api_ip"`
	Port  int    `yaml:"port"`
}

type ServerConf struct {
	IP            string        `yaml:"ip"`
	HostName      string        `yaml:"hostname"`
	Roles         []ServiceRole `yaml:"roles"`
	ProbeInterval int           `yaml:"probe_interval"`
	ExporterPort  string        `yaml:"exporter_port"`
	GrpcAddr      string        `yaml:"grpc_addr"`
	HaHttpPort    int           `yaml:"ha_http_port"`
}

type PrometheusConf struct {
	Addr string `yaml:"addr"`
}

type DNSConf struct {
	Ip        string `yaml:"ip"`
	ConfigDir string `yaml:"config_dir"`
	ProxyPort int    `yaml:"proxy_port"`
}

func LoadConfig(path string) (*MonitorConfig, error) {
	var conf MonitorConfig
	conf.Path = path
	if err := conf.Reload(); err != nil {
		return nil, err
	}

	return &conf, nil
}

func (c *MonitorConfig) Reload() error {
	var newConf MonitorConfig
	if err := configure.Load(&newConf, c.Path); err != nil {
		return err
	}

	newConf.Path = c.Path
	*c = newConf
	return nil
}
