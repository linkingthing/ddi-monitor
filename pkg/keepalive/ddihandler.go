package keepalive

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strconv"

	pb "github.com/linkingthing/ddi-monitor/pkg/proto"
)

const (
	DHCPStartCmd     = "keactrl start"
	DHCPStopCmd      = "keactrl stop"
	DNSName          = "named"
	DNSConfName      = "named.conf"
	DNSProxy         = "rndc"
	DNSProxyConfName = "rndc.conf"
	ReloadNginx      = "docker exec -i ddi-nginx nginx -s reload"
)

type DDIHandler struct {
	dnsStartCmd        string
	dnsStopCmd         string
	dnsConfigCmdPrefix string
}

func newDDIHandler(dnsAddr, dnsDir string, dnsProxyPort int) *DDIHandler {
	var buf bytes.Buffer
	buf.WriteString(filepath.Join(dnsDir, "rndc"))
	buf.WriteString(" -c ")
	buf.WriteString(filepath.Join(dnsDir, "rndc.conf"))
	buf.WriteString(" -s ")
	buf.WriteString(dnsAddr)
	buf.WriteString(" -p ")
	buf.WriteString(strconv.Itoa(dnsProxyPort))

	return &DDIHandler{
		dnsStartCmd:        filepath.Join(dnsDir, "named") + " -c " + filepath.Join(dnsDir, "named.conf"),
		dnsStopCmd:         filepath.Join(dnsDir, "rndc") + " stop",
		dnsConfigCmdPrefix: buf.String(),
	}
}

func (h *DDIHandler) startDNS(req *pb.StartDNSRequest) error {
	if isRunning, err := checkDNSIsRunning(); err != nil {
		return err
	} else if isRunning {
		return nil
	}

	return runCommand(h.dnsStartCmd)
}

func (h *DDIHandler) startDHCP(req *pb.StartDHCPRequest) error {
	if isRunning, err := checkDHCPIsRunning(); err != nil {
		return err
	} else if isRunning {
		return nil
	}

	return runCommand(DHCPStartCmd)
}

func runCommand(cmdline string) error {
	cmd := exec.Command("bash", "-c", cmdline)
	return cmd.Run()
}

func (h *DDIHandler) stopDNS(req *pb.StopDNSRequest) error {
	return runCommand(h.dnsStopCmd)
}

func (h *DDIHandler) stopDHCP(req *pb.StopDHCPRequest) error {
	return runCommand(DHCPStopCmd)
}

func (h *DDIHandler) reconfigDNS() error {
	return runCommand(h.genDnsCmd("reconfig"))
}

func (h *DDIHandler) genDnsCmd(params ...string) string {
	var buf bytes.Buffer
	for _, param := range params {
		buf.WriteString(" ")
		buf.WriteString(param)
	}

	return h.dnsConfigCmdPrefix + buf.String()
}

func (h *DDIHandler) reloadDNS() error {
	return runCommand(h.genDnsCmd("reload"))
}

func (h *DDIHandler) addDNSZone(req *pb.AddDNSZoneRequest) error {
	return runCommand(h.genDnsCmd("addzone", req.GetZoneName(), "in", req.GetViewName(), "{ type master; file \""+req.GetZoneFile()+"\";};"))
}

func (h *DDIHandler) updateDNSZone(req *pb.UpdateDNSZoneRequest) error {
	return runCommand(h.genDnsCmd("modzone", req.GetZoneName(), "in", req.GetViewName(), "{ type master; file \""+req.GetZoneFile()+"\";};"))
}

func (h *DDIHandler) deleteDNSZone(req *pb.DeleteDNSZoneRequest) error {
	return runCommand(h.genDnsCmd("delzone -clean", req.GetZoneName(), "in", req.GetViewName()))
}

func (h *DDIHandler) dumpDNSAllZonesConfig() error {
	return runCommand(h.genDnsCmd("sync -clean"))
}

func (h *DDIHandler) dumpDNSZoneConfig(req *pb.DumpDNSZoneConfigRequest) error {
	return runCommand(h.genDnsCmd("sync -clean", req.GetZoneName(), "in", req.GetViewName()))
}

func (handler *DDIHandler) reloadNginxConfig() error {
	return runCommand(ReloadNginx)
}