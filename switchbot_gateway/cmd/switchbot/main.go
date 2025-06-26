package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"gopkg.in/ini.v1"
	"tinygo.org/x/bluetooth"
)

type switchbotIdentifier string

const (
	switchbotIdentifierLivingRoom switchbotIdentifier = "f4:ae:1a:28:fe:48"
	switchbotIdentifierBedroom    switchbotIdentifier = "f5:bf:47:e6:11:79"
	switchbotIdentifierHallway    switchbotIdentifier = "da:5e:04:87:5c:ef"
)

type Gateway struct {
	apiUrl       string
	apiKey       string
	lastMessages map[string]time.Time
}

type SwitchbotData struct {
	Mac      string    `json:"mac"`
	Temp     float32   `json:"temp"`
	Humidity byte      `json:"humidity"`
	Battery  byte      `json:"battery"`
	Date     time.Time `json:"date"`
}

func (gateway *Gateway) onScan(_ *bluetooth.Adapter, device bluetooth.ScanResult) {
	slog.Debug(fmt.Sprintf("found device: %s, address: %s\n", device.LocalName(), device.Address))
	for _, service := range device.AdvertisementPayload.ServiceData() {
		mac := device.Address.String()
		serviceData := service.Data

		if len(serviceData) < 6 {
			continue
		}

		deviceType := serviceData[0] & 0b01111111
		if mac == "F4:AE:1A:28:FE:48" {
			slog.Info(fmt.Sprintf("Found living room: 0x%x\n", deviceType))
		}
		if deviceType != 0x69 {
			continue
		}
		tempInteger := serviceData[4] & 0b01111111
		tempDecimal := serviceData[3] & 0b00001111
		tempFloat := float32(tempInteger) + float32(tempDecimal)/10.0
		tempFlag := (serviceData[4] & 0b10000000) >> 7
		humidity := serviceData[5] & 0b01111111
		battery := serviceData[2] & 0b01111111
		if battery == 0 {
			continue
		}

		slog.Info(fmt.Sprintf("[%v] deviceType: %x temp: %d.%d, humidity: %d, battery: %d, tempFlag: %d\n", mac, deviceType, tempInteger, tempDecimal, humidity, battery, tempFlag))
		go gateway.sendData(SwitchbotData{strings.ToLower(mac), tempFloat, humidity, battery, time.Now()})
	}
}

func (gateway *Gateway) sendData(data SwitchbotData) {
	if lastMessage, exists := gateway.lastMessages[data.Mac]; exists {
		if time.Since(lastMessage) < 60*time.Second {
			return
		}
	}
	jsonBody, err := json.Marshal(data)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to marshal data: %v\n", err))
		return
	}
	request, err := http.NewRequest("POST", gateway.apiUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create request: %v\n", err))
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", gateway.apiKey)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		slog.Error(fmt.Sprintf("failed to send data: %v\n", err))
		return
	}
	if response.StatusCode >= 400 {
		slog.Error(fmt.Sprintf("failed to send data, status code: %d\n", response.StatusCode))
		return
	}
	gateway.lastMessages[data.Mac] = time.Now()
}

func run() (err error) {
	slog.Info("starting bluetooth scan...")
	cfg, err := ini.Load("config.ini")
	if err != nil {
		err = fmt.Errorf("failed to read config.ini: %v", err)
		return
	}

	apiConfig := cfg.Section("API")

	dataEndpoint := apiConfig.Key("DATA_ENDPOINT").String()
	if dataEndpoint == "" {
		err = fmt.Errorf("data_endpoint is not set in config.ini")
		return
	}
	apiKey := apiConfig.Key("API_KEY").String()
	if apiKey == "" {
		err = fmt.Errorf("api_key is not set in config.ini")
		return
	}

	gateway := &Gateway{
		apiUrl:       apiConfig.Key("DATA_ENDPOINT").String(),
		apiKey:       apiConfig.Key("API_KEY").String(),
		lastMessages: make(map[string]time.Time),
	}

	adapter := bluetooth.DefaultAdapter
	if err = adapter.Enable(); err != nil {
		err = fmt.Errorf("failed to enable bluetooth adapter: %w", err)
		return
	}

	if err = adapter.Scan(gateway.onScan); err != nil {
		err = fmt.Errorf("failed to start scanning: %w", err)
		return
	}

	return
}

func main() {
	err := run()
	if err != nil {
		slog.Error("Error: %v\n", err)
		return
	}
}
