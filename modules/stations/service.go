package stations

import (
	"encoding/json"
	"errors"
	"github.com/marchella2/api-mrtj-schedule/modules/common/client"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	GetAllStations() (response []StationResponse, err error)
	CheckScheduleByStation(id string) (response []ScheduleResponse, err error)
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

func (s *service) CheckScheduleByStation(id string) (response []ScheduleResponse, err error) {
	var url = "https://jakartamrt.co.id/id/val/stasiuns"

	responseByte, err := client.DoRequest(s.client, url)
	if err != nil {
		return nil, err
	}

	var schedules []Schedule
	err = json.Unmarshal(responseByte, &schedules)
	if err != nil {
		return nil, err
	}

	// select schedule by id stasiun
	var scheduleSelected Schedule
	for _, item := range schedules {
		if item.StationId == id {
			scheduleSelected = item
			break
		}
	}

	if scheduleSelected.StationId == "" {
		err = errors.New("Station not found")
		return nil, err
	}

	response, err = ConvertDataToResponse(scheduleSelected)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func ConvertDataToResponse(schedule Schedule) (response []ScheduleResponse, err error) {
	var (
		LebakBulusTripName = "Stasiun Lebak Bulus"
		BundaranHITripName = "Stasiun Bundaran HI"
	)

	scheduleLebakBulus := schedule.ScheduleLebakBulus
	scheduleBundaranHI := schedule.ScheduleBundaranHI

	scheduleLebakBulusParse, err := ConvertScheduleToTimeFormat(scheduleLebakBulus)
	if err != nil {
		return
	}

	scheduleBundaranHIParse, err := ConvertScheduleToTimeFormat(scheduleBundaranHI)
	if err != nil {
		return
	}

	// convert ke response
	for _, item := range scheduleLebakBulusParse {
		if item.Format("15:04") > time.Now().Format("1504") {
			response = append(response, ScheduleResponse{
				StationName: LebakBulusTripName,
				Time:        item.Format("15:04"),
			})
		}
	}

	for _, item := range scheduleBundaranHIParse {
		if item.Format("15:04") > time.Now().Format("1504") {
			response = append(response, ScheduleResponse{
				StationName: BundaranHITripName,
				Time:        item.Format("15:04"),
			})
		}
	}

	return
}

func ConvertScheduleToTimeFormat(schedule string) (response []time.Time, err error) { // function untuk parsing dari schedule string lebih dari satu menjadi time
	var (
		parsedTime time.Time
		schedules  = strings.Split(schedule, ",")
	)

	for _, item := range schedules {
		trimedTime := strings.TrimSpace(item)
		if trimedTime == "" {
			continue
		}

		parsedTime, err = time.Parse("15:04", trimedTime) // untuk parsing string menjadi time
		if err != nil {
			err = errors.New("Invalid time format " + trimedTime)
			return
		}

		response = append(response, parsedTime)
	}

	return
}
