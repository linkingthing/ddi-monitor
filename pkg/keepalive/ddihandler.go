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

func newDDIHandler(dnsIp, dnsDir string, dnsProxyPort int) *DDIHandler {
	return &DDIHandler{
		dnsStartCmd: joinStringWithSpace(filepath.Join(dnsDir, DNSName), "-c", filepath.Join(dnsDir, DNSConfName)),
		dnsStopCmd:  joinStringWithSpace(filepath.Join(dnsDir, DNSProxy), "stop"),
		dnsConfigCmdPrefix: joinStringWithSpace(filepath.Join(dnsDir, DNSProxy), "-c", filepath.Join(dnsDir, DNSProxyConfName),
			"-s", dnsIp, "-p", strconv.Itoa(dnsProxyPort)),
	}
}

func joinStringWithSpace(params ...string) string {
	var buf bytes.Buffer
	for _, param := range params {
		buf.WriteString(" ")
		buf.WriteString(param)
	}
	return buf.String()
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
	return h.dnsConfigCmdPrefix + joinStringWithSpace(params...)
}

func (h *DDIHandler) reloadDNS() error {
	return runCommand(h.genDnsCmd("reload"))
}

func (h *DDIHandler) addDNSZone(req *pb.AddDNSZoneRequest) error {
	return runCommand(h.genDnsCmd("addzone", req.GetZone().GetZoneName(), "in", req.GetZone().GetViewName(),
		"'{ type "+req.GetZone().ZoneRole+
			"; file \""+req.GetZone().GetZoneFile()+
			"\"; allow-transfer {key key"+req.GetZone().GetViewName()+
			";}; also-notify {"+req.GetZone().ZoneSlaves+
			"}; masters {"+req.GetZone().ZoneMasters+"};};'"))
}

func (h *DDIHandler) updateDNSZone(req *pb.UpdateDNSZoneRequest) error {
	if err := runCommand(h.genDnsCmd("freeze",
		req.GetZone().GetZoneName(), "in", req.GetZone().GetViewName())); err != nil {
		return err
	}

	if err := runCommand(h.genDnsCmd("modzone", req.GetZone().GetZoneName(), "in", req.GetZone().GetViewName(),
		"'{ type "+req.GetZone().ZoneRole+
			"; file \""+req.GetZone().GetZoneFile()+
			"\"; allow-transfer {key key"+req.GetZone().GetViewName()+
			";}; also-notify {"+req.GetZone().ZoneSlaves+
			"}; masters {"+req.GetZone().ZoneMasters+"};};'")); err != nil {
		return err
	}

	return runCommand(h.genDnsCmd("thaw", req.GetZone().GetZoneName(), "in", req.GetZone().GetViewName()))
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

func (h *DDIHandler) reloadNginxConfig() error {
	return runCommand(ReloadNginx)
}
