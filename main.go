package main

import (
	"github.com/gin-gonic/gin"
	"github.com/marchella2/api-mrtj-schedule/modules/stations"
)

func main() {
	Init()
}

func Init() {
	var (
		router = gin.Default()
		api    = router.Group("/v1/api")
	)

	stations.Initiate(api)

	router.Run(":9091")
}
