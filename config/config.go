package config

import (
	"github.com/zdnscloud/cement/configure"
)

type ServiceRole string

const (
	ServiceRoleDHCP       ServiceRole = "dhcp"
	ServiceRoleDNS        ServiceRole = "dns"
	ServiceRoleController ServiceRole = "controller"
)

type MonitorConfig struct {
	Path           string         `yaml:"-"`
	Server         ServerConf     `yaml:"server"`
	ControllerAddr string         `yaml:"controller_addr"`
	Prometheus     PrometheusConf `yaml:"prometheus"`
	Master         string         `yaml:"master"`
	VIP            string         `yaml:"vip"`
	DNS            DNSConf        `yaml:"dns"`
}

type ServerConf struct {
	IP            string        `yaml:"ip"`
	HostName      string        `yaml:"hostname"`
	Roles         []ServiceRole `yaml:"roles"`
	ProbeInterval int           `yaml:"probe_interval"`
	ExporterPort  string        `yaml:"exporter_port"`
	GrpcAddr      string        `yaml:"grpc_addr"`
}

type PrometheusConf struct {
	Addr string `yaml:"addr"`
}

type DNSConf struct {
	ConfigDir string `yaml:"config_dir"`
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
