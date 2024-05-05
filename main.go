package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// MainInfo представляет основные погодные данные.
type MainInfo struct {
	Temp     float64 `json:"temp"`
	Pressure int     `json:"pressure"`
	Humidity int     `json:"humidity"`
}

// WeatherCondition представляет информацию о текущем состоянии погоды.
type WeatherCondition struct {
	Main        string `json:"main"`
	Description string `json:"description"`
}

// WeatherData представляет данные о погоде, полученные от API.
type WeatherData struct {
	Main    MainInfo           `json:"main"`
	Weather []WeatherCondition `json:"weather"`
	Name    string             `json:"name"`
}

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

// Базовая ссылка, к которой мы будем добовлять желаемый город и наш ApiKey.
const baseURL = "https://api.openweathermap.org/data/2.5/weather"

// Функция, которая читает наш файл с ApiKey и предоставляет нам его в удобном формате.
func loadApiConfig(fileName string) (apiConfigData, error) {
	// Читаем файл fileName и загружаем его содержимое в bytes(fileName - это json файл в нашем случае).
	bytes, err := os.ReadFile(fileName)

	// Обрабатываем ошибку на случай, если она будет.
	if err != nil {
		return apiConfigData{}, err
	}

	// Создаем переменную "apiKey" типа apiConfigData.
	var apiKey apiConfigData

	// Используем функцию json.Unmarshal, чтобы преобразовать данные из json файла в необходимый нам тип apiConfigData.
	err = json.Unmarshal(bytes, &apiKey)

	// Обрабатываем ошибку на случай, если она будет.
	if err != nil {
		return apiConfigData{}, err
	}

	// В случае прохождения всех этапов без ошибок выводим наш объект apiKey типа apiConfigData и nil.
	return apiKey, nil
}

// Функция, которая принимает на вход имя города чью погоду мы хотим узнать, на выходе предоставляет нам информацию в виде объекта типа MainWeatherData.
func FetchWeather(cityName string) (*WeatherData, error) {
	apiKey, err := loadApiConfig(".apiConfig")

	if err != nil {
		return nil, err
	}

	requestURL := fmt.Sprintf("%s?q=%s&appid=%s", baseURL, cityName, apiKey.OpenWeatherMapApiKey)
	fmt.Println(requestURL)
	response, err := http.Get(requestURL)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка api: статус %v", response.Status)
	}

	var weatherData WeatherData

	if err := json.NewDecoder(response.Body).Decode(&weatherData); err != nil {
		return nil, err
	}

	return &weatherData, nil
}

func HandleWeather(w http.ResponseWriter, r *http.Request) {
	cityName := r.URL.Query().Get("city")

	if cityName == "" {
		http.Error(w, "Город не указан", http.StatusBadRequest)
		return
	}

	data, err := FetchWeather(cityName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/weather/", HandleWeather)
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
