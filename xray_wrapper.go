package libXray

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/xtls/libxray/nodep"
	"github.com/xtls/libxray/xray"
)

var tunFile *os.File

// Структуры запросов (обязательны для десериализации JSON)
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

func manageFd(fd int) {
	if tunFile != nil {
		tunFile.Close()
		tunFile = nil
	}
	if fd > 0 {
		tunFile = os.NewFile(uintptr(fd), "tun")
	}
}

func RunXray(fd int, base64Text string) string {
	var response nodep.CallResponse[string]
	manageFd(fd)
	req, err := base64.StdEncoding.DecodeString(base64Text)
	if err != nil {
		return response.EncodeToBase64("", err)
	}
	var request RunXrayRequest
	if err := json.Unmarshal(req, &request); err != nil {
		return response.EncodeToBase64("", err)
	}
	err = xray.RunXray(request.DatDir, request.MphCachePath, request.ConfigPath)
	return response.EncodeToBase64("", err)
}

func RunXrayFromJSON(fd int, base64Text string) string {
	var response nodep.CallResponse[string]
	manageFd(fd)
	req, err := base64.StdEncoding.DecodeString(base64Text)
	if err != nil {
		return response.EncodeToBase64("", err)
	}
	var request RunXrayFromJSONRequest
	if err := json.Unmarshal(req, &request); err != nil {
		return response.EncodeToBase64("", err)
	}
	err = xray.RunXrayFromJSON(request.DatDir, request.MphCachePath, request.ConfigJSON)
	return response.EncodeToBase64("", err)
}

func BuildMphCache(fd int, base64Text string) string {
	var response nodep.CallResponse[string]
	manageFd(fd)
	req, err := base64.StdEncoding.DecodeString(base64Text)
	if err != nil {
		return response.EncodeToBase64("", err)
	}
	var request RunXrayRequest
	if err := json.Unmarshal(req, &request); err != nil {
		return response.EncodeToBase64("", err)
	}
	err = xray.BuildMphCache(request.DatDir, request.MphCachePath, request.ConfigPath)
	return response.EncodeToBase64("", err)
}

func StopXray() string {
	var response nodep.CallResponse[string]
	manageFd(0)
	err := xray.StopXray()
	return response.EncodeToBase64("", err)
}

func GetXrayState() bool {
	return xray.GetXrayState()
}

func XrayVersion() string {
	var response nodep.CallResponse[string]
	return response.EncodeToBase64(xray.XrayVersion(), nil)
}