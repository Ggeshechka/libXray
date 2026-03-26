package libXray

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/xtls/libxray/geo"
	"github.com/xtls/libxray/nodep"
	"github.com/xtls/libxray/xray"
)

var tunFile *os.File

// --- Оригинальные структуры (не удалять для компиляции) ---

type RunXrayRequest struct {
	DatDir       string `json:"datDir,omitempty"`
	MphCachePath string `json:"mphCachePath,omitempty"`
	ConfigPath   string `json:"configPath,omitempty"`
}

type RunXrayFromJSONRequest struct {
	DatDir       string `json:"datDir,omitempty"`
	MphCachePath string `json:"mphCachePath,omitempty"`
	ConfigJSON   string `json:"configJSON,omitempty"`
}

type CountGeoDataRequest struct {
	DatDir  string `json:"datDir,omitempty"`
	Name    string `json:"name,omitempty"`
	GeoType string `json:"geoType,omitempty"`
}

// --- Внутренняя автоматика ---

func manageFd(fd int) {
	if tunFile != nil {
		tunFile.Close()
		tunFile = nil
	}
	if fd > 0 {
		tunFile = os.NewFile(uintptr(fd), "tun")
	}
}

// --- Основные функции с поддержкой FD ---

func RunXray(fd int, base64Text string) string {
	var response nodep.CallResponse[string]
	manageFd(fd)
	req, _ := base64.StdEncoding.DecodeString(base64Text)
	var request RunXrayRequest
	json.Unmarshal(req, &request)
	err := xray.RunXray(request.DatDir, request.MphCachePath, request.ConfigPath)
	return response.EncodeToBase64("", err)
}

func RunXrayFromJSON(fd int, base64Text string) string {
	var response nodep.CallResponse[string]
	manageFd(fd)
	req, _ := base64.StdEncoding.DecodeString(base64Text)
	var request RunXrayFromJSONRequest
	json.Unmarshal(req, &request)
	err := xray.RunXrayFromJSON(request.DatDir, request.MphCachePath, request.ConfigJSON)
	return response.EncodeToBase64("", err)
}

func BuildMphCache(fd int, base64Text string) string {
	var response nodep.CallResponse[string]
	manageFd(fd)
	req, _ := base64.StdEncoding.DecodeString(base64Text)
	var request RunXrayRequest
	json.Unmarshal(req, &request)
	err := xray.BuildMphCache(request.DatDir, request.MphCachePath, request.ConfigPath)
	return response.EncodeToBase64("", err)
}

func StopXray() string {
	var response nodep.CallResponse[string]
	manageFd(0) // Закрываем FD
	err := xray.StopXray()
	return response.EncodeToBase64("", err)
}

// --- Остальные оригинальные функции ---

func CountGeoData(base64Text string) string {
	var response nodep.CallResponse[string]
	req, _ := base64.StdEncoding.DecodeString(base64Text)
	var request CountGeoDataRequest
	json.Unmarshal(req, &request)
	err := geo.CountGeoData(request.DatDir, request.Name, request.GeoType)
	return response.EncodeToBase64("", err)
}

func GetXrayState() bool {
	return xray.GetXrayState()
}

func XrayVersion() string {
	var response nodep.CallResponse[string]
	return response.EncodeToBase64(xray.XrayVersion(), nil)
}

// Оставил заглушки для совместимости
func TestXray(base64Text string) string { return "" }
func QueryStats(server string) string   { return "" }
func ReadGeoFiles(base64Text string) string { return "" }