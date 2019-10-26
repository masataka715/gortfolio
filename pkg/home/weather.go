package home

import (
	"encoding/json"
	"fmt"
	"gortfolio/config"
	"net/http"
	"net/url"
)

const (
	// OpenWeatherApiKey = "42ec1c264bf771dfe9e06897261c2b47"
	EndPoint = "https://api.openweathermap.org/data/2.5/weather"
)

type Coord struct {
	Lon float32
	Lat float32
}

type Weather struct {
	Id          int
	Main        string
	Description string
	Icon        string
}

type Main struct {
	Temp     float32
	Pressure int
	Humidity int
	TempMin  float32
	TempMax  float32
}

type Wind struct {
	Speed float32
	Deg   int
}

type Clouds struct {
	All int
}

type Sys struct {
	Type    int
	Id      int
	Message float32
	Country string
	Sunrise int
	Sunset  int
}

type Response struct {
	Coord      Coord
	Weather    []Weather
	Base       string
	Main       Main
	Visibility int
	Wind       Wind
	Clouds     Clouds
	Dt         int
	Sys        Sys
	Id         int
	Name       string
	Cod        int
}

type ViewData struct {
	Weather   string
	Temp      string
	WindSpeed string
}

func GetWeather() *ViewData {
	values := url.Values{}
	values.Add("q", "Tokyo")
	values.Add("APPID", config.Config.OpenWeatherApiKey)

	req, err := http.NewRequest("GET", EndPoint, nil)
	if err != nil {
		panic(err)
	}

	req.URL.RawQuery = values.Encode()
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	response := Response{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	weatherEn := response.Weather[0].Main
	var weatherJa string
	switch weatherEn {
	case "Clear":
		weatherJa = "快晴"
	case "Clouds":
		weatherJa = "くもり"
		if response.Weather[0].Description == "few clouds" {
			weatherJa = "晴れ"
		}
	case "Rain":
		weatherJa = "雨"
	case "Snow":
		weatherJa = "雪"
	case "Drizzle":
		weatherJa = "小雨"
	case "Thunderstorm":
		weatherJa = "雷雨"
	default:
		weatherJa = "その他"
	}

	tempFloat := response.Main.Temp - 273.15
	temp := fmt.Sprintf("%.1f", tempFloat) + "℃"

	// windSpeed := strconv.FormatFloat(float64(response.Wind.Speed), 'f', -1, 64)
	windSpeed := fmt.Sprintf("%.1f", response.Wind.Speed) + "メートル / 秒"

	viewData := &ViewData{
		Weather:   weatherJa,
		Temp:      temp,
		WindSpeed: windSpeed,
	}

	return viewData
}
