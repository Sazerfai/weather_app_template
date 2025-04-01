package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Структуры для геокодинга и погоды
type GeocodingResponse struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type WeatherResponse struct {
	Name    string `json:"name"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

type RequestData struct {
	City string `json:"city"`
}

var apiKey = "YOUR_API_KEY"

func getCoordinates(city string) (float64, float64, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var geocoding []GeocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&geocoding); err != nil || len(geocoding) == 0 {
		return 0, 0, fmt.Errorf("город не найден")
	}

	return geocoding[0].Lat, geocoding[0].Lon, nil
}

func getWeather(lat, lon float64) (*WeatherResponse, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric&lang=ru", lat, lon, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, err
	}

	return &weather, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData RequestData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	lat, lon, err := getCoordinates(requestData.City)
	if err != nil {
		http.Error(w, "Город не найден", http.StatusNotFound)
		return
	}

	weather, err := getWeather(lat, lon)
	if err != nil {
		http.Error(w, "Ошибка получения погоды", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}

func main() {
	//TODO: Вписать свой ключ API от OpenWeather здесь или указать в env
	apiKey = ""
	//apiKey = os.Getenv("OPENWEATHER_API_KEY")

	http.HandleFunc("/weather", weatherHandler)

	log.Println("Сервер запущен на :5000")
	log.Println("Локальный доступ - http://localhost:81")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}
}
