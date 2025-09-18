package stations

import (
	"encoding/json"
	"github.com/marchella2/api-mrtj-schedule/modules/common/client"
	"net/http"
	"time"
)

type Service interface {
	GetAllStations() (response []StationResponse, err error)
}

type service struct {
	client *http.Client // berfungsi untuk hit API MRT Jakarta
}

func NewService() Service { // memanggil interface yang dibuat
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *service) GetAllStations() (response []StationResponse, err error) {
	var url = "https://jakartamrt.co.id/id/val/stasiuns"

	// hit url
	responseByte, err := client.DoRequest(s.client, url)
	if err != nil {
		return nil, err
	}

	var stations []Station
	err = json.Unmarshal(responseByte, &stations)

	// mapping response
	for _, item := range stations {
		response = append(response, StationResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}
	return
}
