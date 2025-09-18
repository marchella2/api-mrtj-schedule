package stations

import "C"
import (
	"github.com/gin-gonic/gin"
	"github.com/marchella2/api-mrtj-schedule/modules/common/response"
	"net/http"
)

func Initiate(router *gin.RouterGroup) {
	stationService := NewService()

	station := router.Group("/stations")
	station.GET("", func(c *gin.Context) {
		// redirect to service
		GetAllStations(c, stationService)
	})
}

func GetAllStations(c *gin.Context, service Service) {
	datas, err := service.GetAllStations()

	// jika gagal mengambil data ke service getAllStations
	if err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Success get all stations",
		Data:    datas,
	})

}
