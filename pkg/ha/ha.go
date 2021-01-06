package ha

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linkingthing/ddi-monitor/config"
	"github.com/linkingthing/ddi-monitor/pkg/ha/handler"
)

func Server(conf *config.MonitorConfig) {
	router := gin.Default()
	handler.NewHandler(conf)

	router.POST("start_ha", handler.StartHa)
	router.POST("master_up", handler.MasterUp)
	router.POST("master_down", handler.MasterDown)
	log.Fatal(router.Run(":" + strconv.Itoa(conf.Server.HaHttpPort)))
}
