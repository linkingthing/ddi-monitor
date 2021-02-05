package ha

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linkingthing/ddi-monitor/config"
	"github.com/linkingthing/ddi-monitor/pkg/ha/handler"
)

func Server(conf *config.MonitorConfig) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	haHandler := handler.NewHandler(conf)

	router.POST("start_ha", haHandler.StartHa)
	router.POST("master_up", haHandler.MasterUp)
	router.POST("master_down", haHandler.MasterDown)
	log.Fatal(router.Run(":" + strconv.Itoa(conf.Server.HaHttpPort)))
}
