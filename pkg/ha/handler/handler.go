package handler

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/linkingthing/ddi-monitor/config"
	"github.com/linkingthing/ddi-monitor/pkg/util"
)

type HaHandler struct {
	ControllerAddr string
	Client         *http.Client
	DhcpHa         bool
	DnsHa          bool
	ControllerHa   bool
	PgHaCliDir     string
	Vip            string
}

type HaRequest struct {
	MasterIP string             `json:"masterIP"`
	Role     config.ServiceRole `json:"role"`
	Vip      string             `json:"vip"`
	SlaveIP  string             `json:"slaveIP"`
}

type HaResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *HaResponse) Error(msg string) *HaResponse {
	r.Code = 0
	r.Message = msg
	return r
}

func (r *HaResponse) Success(msg string) *HaResponse {
	r.Code = 1
	r.Message = msg
	return r
}

var (
	PgHaCmd = "%s -c %s -m %s -s %s"
)

func NewHandler(conf *config.MonitorConfig) *HaHandler {
	handler := &HaHandler{
		ControllerAddr: conf.Controller.ApiIp + ":" + strconv.Itoa(conf.Controller.Port),
		PgHaCliDir:     conf.PgHaCliDir,
		Vip:            conf.VIP,
		Client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	for _, role := range conf.Server.Roles {
		if role == config.ServiceRoleDHCP && conf.Master != "" {
			handler.DhcpHa = true
		}
		if role == config.ServiceRoleDNS && conf.Master != "" {
			handler.DnsHa = true
		}
		if role == config.ServiceRoleController && conf.Master != "" {
			handler.ControllerHa = true
		}
	}

	return handler
}

func (handler *HaHandler) StartHa(ctx *gin.Context) {
	var haRequest *HaRequest
	haResponse := &HaResponse{}
	if err := ctx.ShouldBindJSON(&haRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, haResponse.Error(err.Error()))
		return
	}

	if handler.DhcpHa {
		if err := runCommand(fmt.Sprintf(PgHaCmd, handler.PgHaCliDir, config.ActionStartHa,
			haRequest.MasterIP, haRequest.SlaveIP)); err != nil {
			ctx.JSON(http.StatusBadRequest, haResponse.Error(err.Error()))
			return
		}
	}

	ctx.Status(http.StatusOK)
}

func (handler *HaHandler) MasterUp(ctx *gin.Context) {
	var haRequest *HaRequest
	haResponse := &HaResponse{}
	if err := ctx.ShouldBindJSON(&haRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, haResponse.Error(err.Error()))
		return
	}

	switch {
	case handler.DhcpHa:
		haRequest.Role = config.ServiceRoleDHCP
		if err := runCommand(fmt.Sprintf(PgHaCmd, handler.PgHaCliDir, config.ActionMasterUp,
			haRequest.MasterIP, haRequest.SlaveIP)); err != nil {
			ctx.JSON(http.StatusBadRequest, haResponse.Error(err.Error()))
			return
		}
	case handler.DnsHa:
		haRequest.Role = config.ServiceRoleDNS
	case handler.ControllerHa:
		haRequest.Role = config.ServiceRoleController
	}

	haRequest.Vip = handler.Vip
	if err := handler.notifyController(config.ActionMasterUp, haRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, haResponse.Error(err.Error()))
		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *HaHandler) MasterDown(ctx *gin.Context) {
	var haRequest *HaRequest
	haResponse := &HaResponse{}
	if err := ctx.ShouldBindJSON(&haRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, haResponse.Error(err.Error()))
		return
	}
	switch {
	case handler.DhcpHa:
		haRequest.Role = config.ServiceRoleDHCP
		if err := runCommand(fmt.Sprintf(PgHaCmd, handler.PgHaCliDir, config.ActionMasterDown,
			haRequest.MasterIP, haRequest.SlaveIP)); err != nil {
			ctx.JSON(http.StatusBadRequest, haResponse.Error(err.Error()))
			return
		}
	case handler.DnsHa:
		haRequest.Role = config.ServiceRoleDNS
	case handler.ControllerHa:
		haRequest.Role = config.ServiceRoleController
	}

	haRequest.Vip = handler.Vip
	if err := handler.notifyController(config.ActionMasterDown, haRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, haResponse.Error(err.Error()))
		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *HaHandler) notifyController(action string, haRequest *HaRequest) error {
	if handler.ControllerAddr == "" {
		return nil
	}

	token, err := util.GetToken(handler.Client, handler.ControllerAddr)
	retryTime := 5
	for err != nil || retryTime == 0 {
		time.Sleep(time.Second * 2)
		retryTime--
		token, err = util.GetToken(handler.Client, handler.ControllerAddr)
	}

	_, err = util.HttpRequest(handler.Client, http.MethodPost,
		util.GenControllerRequestUrl(handler.ControllerAddr, action, haRequest.MasterIP),
		token, haRequest)
	return err
}

func runCommand(cmdline string) error {
	cmd := exec.Command("bash", "-c", cmdline)
	return cmd.Run()
}
