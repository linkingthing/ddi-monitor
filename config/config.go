package config

import (
	"github.com/zdnscloud/cement/configure"
)

type MonitorConfig struct {
	Path           string         `yaml:"-"`
	Server         ServerConf     `yaml:"server"`
	ControllerAddr string         `yaml:"controller_addr"`
	Prometheus     PrometheusConf `yaml:"prometheus"`
}

type ServerConf struct {
	IP            string   `yaml:"ip"`
	HostName      string   `yaml:"hostname"`
	Roles         []string `yaml:"roles"`
	ProbeInterval uint     `yaml:"probe_interval"`
	ExporterPort  string   `yaml:"exporter_port"`
}

type PrometheusConf struct {
	Addr string `yaml:"addr"`
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
