package importer

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

type ValueIntf []interface {
}

type ResultData struct {
	Metric Metric
	Values []ValueIntf
}
type Data struct {
	ResultType string
	Result     []ResultData
}
type Response struct {
	Status string
	Data   Data
}

func GetMetric(conf *config.MonitorConfig) (string, string, error) {
	cpuPQL := "100 - (avg(irate(node_cpu_seconds_total{instance=\"" + conf.Server.IP + ":" + physicalExportPort + "\", mode=\"idle\"}[5m])) by (instance) * 100)"
	memPQL := "(node_memory_MemFree_bytes{instance=\"" + conf.Server.IP + ":" + physicalExportPort + "\"}+node_memory_Cached_bytes{instance=\"" +
		conf.Server.IP + ":" + physicalExportPort + "\"}+node_memory_Buffers_bytes{instance=\"" + conf.Server.IP + ":" + physicalExportPort + "\"}) / node_memory_MemTotal_bytes * 100"
	cpuResp, err := getCurrentMetric(conf, cpuPQL)
	if err != nil {
		return "", "", fmt.Errorf("cpu getCurrentMetric fail", err.Error())
	}
	cpu, err := getUsage(cpuResp)
	if err != nil {
		return "", "", fmt.Errorf("cpu getUsage fail", err.Error())
	}
	memResp, err := getCurrentMetric(conf, memPQL)
	if err != nil {
		return "", "", fmt.Errorf("memory getCurrentMetric fail", err.Error())
	}
	mem, err := getUsage(memResp)
	if err != nil {
		return "", "", fmt.Errorf("memory getUsage fail", err.Error())
	}
	return cpu, mem, nil
}

func getCurrentMetric(conf *config.MonitorConfig, pql string) ([]byte, error) {
	param := url.Values{}
	param.Add("query", pql)
	param.Add("start", strconv.FormatInt(time.Now().Unix(), 10))
	param.Add("end", strconv.FormatInt(time.Now().Unix(), 10))
	param.Add("step", "20")
	path := schema + conf.Prometheus.Addr + "/api/v1/query_range?" + param.Encode()
	resp, err := http.Get(path)
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

func getUsage(body []byte) (string, error) {
	var rsp Response
	err := json.Unmarshal(body, &rsp)
	if err != nil {
		return "", err
	}
	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	if err := d.Decode(&rsp); err != nil {
		return "", err
	}

	if rsp.Status != STATUS_SUCCCESS {
		return "", err
	}
	for _, r := range rsp.Data.Result {
		for _, vs := range r.Values {
			for _, v := range vs {
				if t, ok := v.(string); ok {
					if f, err := strconv.ParseFloat(t, 64); err == nil {
						return fmt.Sprintf("%.2f", f), nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("not found")
}
