package metric

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/linkingthing/ddi-monitor/config"
)

var (
	schema             = "http://"
	physicalExportPort = "9100"
	STATUS_SUCCCESS    = "success"
)

type Metric struct {
	Name     string `json:"__name__"`
	Instance string
	Job      string
}

type ValueIntf [2]interface {
}
type ValueIntfOne interface {
}

type Result struct {
	Metric Metric
	Value  []ValueIntfOne
	Values []ValueIntf
}
type Data struct {
	ResultType string
	Result     []Result
}
type Response struct {
	Status string
	Data   Data
}

func New(conf *config.MonitorConfig) (*float64, *float64, error) {
	memPQL := "(node_memory_MemFree_bytes{instance=\"" + conf.Server.IP + ":" + physicalExportPort + "\"}+node_memory_Cached_bytes{instance=\"" +
		conf.Server.IP + ":" + physicalExportPort + "\"}+node_memory_Buffers_bytes{instance=\"" + conf.Server.IP + ":" + physicalExportPort + "\"}) / node_memory_MemTotal_bytes * 100"
	cpuPQL := "100 - (avg(irate(node_cpu_seconds_total{instance=\"" + conf.Server.IP + ":" + physicalExportPort + "\", mode=\"idle\"}[5m])) by (instance) * 100)"
	memResp, err := getCurrentMetric(conf, memPQL)
	if err != nil {
		return nil, nil, err
	}
	mem, err := getUsage(memResp)
	if err != nil {
		return nil, nil, err
	}
	cpuResp, err := getCurrentMetric(conf, cpuPQL)
	if err != nil {
		return nil, nil, err
	}
	cpu, err := getUsage(cpuResp)
	if err != nil {
		return nil, nil, err
	}
	return cpu, mem, nil
}

func getCurrentMetric(conf *config.MonitorConfig, pql string) ([]byte, error) {
	param := url.Values{}
	param.Add("query", pql)
	param.Add("start", strconv.FormatInt(time.Now().Unix(), 10))
	param.Add("end", strconv.FormatInt(time.Now().Unix(), 10))
	param.Add("step", "20")
	url := schema + conf.ControllerAddr + "?" + param.Encode()
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func getUsage(body []byte) (*float64, error) {
	var rsp Response
	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	err := d.Decode(&rsp)
	if rsp.Status != STATUS_SUCCCESS {
		return nil, err
	}
	for _, v := range rsp.Data.Result {
		if len(v.Values) > 1 {
			tmp := v.Values[0][1].(*float64)
			return tmp, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
