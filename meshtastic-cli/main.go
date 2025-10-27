// Copyright (C) 2025 Michael Graff
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/skandragon/meshmgr/meshtastic-cli/proto/meshtastic"
	"github.com/tarm/serial"
	protobuf "google.golang.org/protobuf/proto"
)

const (
	// Meshtastic serial protocol constants
	MAGIC_START = 0x94
	MAGIC_END   = 0xc3

	// Max packet size
	MAX_TO_FROM_RADIO_SIZE = 512

	// Config request ID
	WANT_CONFIG_ID = 64
)

type DeviceConfig struct {
	// Basic device info
	NodeNum    uint32 `json:"node_num"`
	DeviceID   []byte `json:"device_id,omitempty"`   // MAC address / unique device ID
	HardwareID string `json:"hardware_id"`
	LongName   string `json:"long_name"`
	ShortName  string `json:"short_name"`

	// Device metadata
	Metadata *pb.DeviceMetadata `json:"metadata,omitempty"`

	// Device configuration (merged from multiple Config messages)
	LocalConfig *pb.LocalConfig `json:"config,omitempty"`

	// Module configuration (merged from multiple ModuleConfig messages)
	LocalModuleConfig *pb.LocalModuleConfig `json:"module_config,omitempty"`

	// Channels (up to 8)
	Channels []*pb.Channel `json:"channels,omitempty"`

	// Status
	ConfigComplete bool `json:"config_complete"`
}

func main() {
	port := flag.String("port", "/dev/tty.usbmodem101", "Serial port device")
	baud := flag.Int("baud", 115200, "Baud rate")
	jsonOutput := flag.Bool("json", false, "Output as JSON")
	adminURL := flag.String("admin-url", "https://meshmanager.svc.rpi.flame.org", "Admin server URL")
	apiKey := flag.String("api-key", os.Getenv("MESHMANAGER_API_KEY"), "API key for authentication")
	meshID := flag.String("mesh-id", "", "Mesh ID to upload config to (required for upload)")
	flag.Parse()

	if !*jsonOutput {
		fmt.Printf("Meshtastic Device Config Reader\n")
		fmt.Printf("================================\n")
		fmt.Printf("Connecting to %s at %d baud...\n", *port, *baud)
	}

	config := &serial.Config{
		Name:        *port,
		Baud:        *baud,
		ReadTimeout: time.Millisecond * 100,
	}

	s, err := serial.OpenPort(config)
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}
	defer func() {
		_ = s.Close()
	}()

	if !*jsonOutput {
		fmt.Println("Connected successfully!")
	}

	if err := wakeDevice(s); err != nil {
		log.Fatalf("Failed to wake device: %v", err)
	}

	packets := make(chan []byte, 100)
	debugLines := make(chan string, 100)
	go readPackets(s, packets, debugLines)

	time.Sleep(time.Millisecond * 100)

	if !*jsonOutput {
		fmt.Println("\nRequesting device configuration...")
	}
	if err := requestConfig(s); err != nil {
		log.Fatalf("Failed to request config: %v", err)
	}

	deviceConfig := &DeviceConfig{
		Channels:          make([]*pb.Channel, 0, 8),
		LocalConfig:       &pb.LocalConfig{},
		LocalModuleConfig: &pb.LocalModuleConfig{},
	}
	timeout := time.After(15 * time.Second)

	for {
		select {
		case packet := <-packets:
			if parseConfigPacket(packet, deviceConfig) {
				// Got config complete
				outputResult(deviceConfig, *jsonOutput)
				uploadConfig(deviceConfig, *adminURL, *apiKey, *meshID, *jsonOutput)
				return
			}

		case <-debugLines:
			// Discard debug output

		case <-timeout:
			if !*jsonOutput {
				fmt.Println("\n⚠️  Timeout waiting for device configuration")
			}
			// Output what we have
			outputResult(deviceConfig, *jsonOutput)
			uploadConfig(deviceConfig, *adminURL, *apiKey, *meshID, *jsonOutput)
			return
		}
	}
}

func outputResult(config *DeviceConfig, jsonOutput bool) {
	if jsonOutput {
		jsonData, _ := json.MarshalIndent(config, "", "  ")
		fmt.Println(string(jsonData))
	} else {
		fmt.Println("\n✅ Device Configuration:")
		fmt.Printf("  Node Number: %d (0x%08x)\n", config.NodeNum, config.NodeNum)
		fmt.Printf("  Hardware ID: %s\n", config.HardwareID)
		fmt.Printf("  Long Name: %s\n", config.LongName)
		fmt.Printf("  Short Name: %s\n", config.ShortName)
		if config.Metadata != nil {
			fmt.Printf("  Firmware: %s\n", config.Metadata.FirmwareVersion)
			fmt.Printf("  Hardware Model: %s\n", config.Metadata.HwModel)
		}
		fmt.Printf("  Channels: %d\n", len(config.Channels))
		fmt.Printf("  Config Complete: %v\n", config.ConfigComplete)
	}
}

func wakeDevice(s *serial.Port) error {
	wakeBytes := make([]byte, 32)
	for i := range wakeBytes {
		wakeBytes[i] = MAGIC_END
	}
	if _, err := s.Write(wakeBytes); err != nil {
		return fmt.Errorf("failed to wake device: %w", err)
	}
	return nil
}

func readPackets(s *serial.Port, packets chan []byte, debugLines chan string) {
	var rxBuf []byte
	var debugLine string

	for {
		b := make([]byte, 1)
		n, err := s.Read(b)
		if err != nil || n != 1 {
			continue
		}

		c := b[0]
		ptr := len(rxBuf)

		if ptr == 0 {
			if c == MAGIC_START {
				rxBuf = append(rxBuf, c)
			} else {
				handleDebugByte(c, &debugLine, debugLines)
			}
		} else if ptr == 1 {
			if c == MAGIC_END {
				rxBuf = append(rxBuf, c)
			} else {
				handleDebugByte(rxBuf[0], &debugLine, debugLines)
				handleDebugByte(c, &debugLine, debugLines)
				rxBuf = nil
			}
		} else if ptr == 2 || ptr == 3 {
			rxBuf = append(rxBuf, c)
			if ptr == 3 {
				msb := uint16(rxBuf[2])
				lsb := uint16(rxBuf[3])
				packetLen := (msb << 8) | lsb

				if packetLen > MAX_TO_FROM_RADIO_SIZE {
					for _, b := range rxBuf {
						handleDebugByte(b, &debugLine, debugLines)
					}
					rxBuf = nil
				}
			}
		} else {
			rxBuf = append(rxBuf, c)

			if len(rxBuf) >= 4 {
				msb := uint16(rxBuf[2])
				lsb := uint16(rxBuf[3])
				packetLen := (msb << 8) | lsb
				totalLen := 4 + int(packetLen)

				if len(rxBuf) >= totalLen {
					packetData := make([]byte, packetLen)
					copy(packetData, rxBuf[4:totalLen])
					packets <- packetData
					rxBuf = nil
				}
			}
		}
	}
}

func handleDebugByte(c byte, debugLine *string, debugLines chan string) {
	if c == '\r' {
		return
	} else if c == '\n' {
		if len(*debugLine) > 0 {
			select {
			case debugLines <- *debugLine:
			default:
			}
			*debugLine = ""
		}
	} else if c >= 32 && c <= 126 {
		*debugLine += string(c)
	}
}

func requestConfig(s *serial.Port) error {
	toRadio := &pb.ToRadio{
		PayloadVariant: &pb.ToRadio_WantConfigId{
			WantConfigId: WANT_CONFIG_ID,
		},
	}

	data, err := protobuf.Marshal(toRadio)
	if err != nil {
		return fmt.Errorf("failed to marshal ToRadio: %w", err)
	}

	packetLen := uint16(len(data))
	packet := make([]byte, 4+len(data))
	packet[0] = MAGIC_START
	packet[1] = MAGIC_END
	packet[2] = uint8(packetLen >> 8)
	packet[3] = uint8(packetLen & 0xFF)
	copy(packet[4:], data)

	if _, err = s.Write(packet); err != nil {
		return fmt.Errorf("failed to send config request: %w", err)
	}
	return nil
}

func parseConfigPacket(packet []byte, config *DeviceConfig) bool {
	fromRadio := &pb.FromRadio{}
	err := protobuf.Unmarshal(packet, fromRadio)
	if err != nil {
		return false
	}

	switch p := fromRadio.PayloadVariant.(type) {
	case *pb.FromRadio_MyInfo:
		if p.MyInfo != nil {
			config.NodeNum = p.MyInfo.MyNodeNum
			config.DeviceID = p.MyInfo.DeviceId
		}

	case *pb.FromRadio_NodeInfo:
		// Only process our own node info
		if p.NodeInfo != nil && p.NodeInfo.Num == config.NodeNum {
			if p.NodeInfo.User != nil {
				config.HardwareID = p.NodeInfo.User.Id
				config.LongName = p.NodeInfo.User.LongName
				config.ShortName = p.NodeInfo.User.ShortName
				// If MyNodeInfo didn't provide device_id, fallback to User.macaddr (deprecated but still sent by some firmware)
				if len(config.DeviceID) == 0 && len(p.NodeInfo.User.Macaddr) > 0 {
					config.DeviceID = p.NodeInfo.User.Macaddr
				}
			}
		}

	case *pb.FromRadio_Metadata:
		if p.Metadata != nil {
			config.Metadata = p.Metadata
		}

	case *pb.FromRadio_Config:
		if p.Config != nil {
			mergeConfig(config.LocalConfig, p.Config)
		}

	case *pb.FromRadio_ModuleConfig:
		if p.ModuleConfig != nil {
			mergeModuleConfig(config.LocalModuleConfig, p.ModuleConfig)
		}

	case *pb.FromRadio_Channel:
		if p.Channel != nil {
			config.Channels = append(config.Channels, p.Channel)
		}

	case *pb.FromRadio_ConfigCompleteId:
		if p.ConfigCompleteId == WANT_CONFIG_ID {
			config.ConfigComplete = true
			// Check if we have all essential fields
			if config.NodeNum != 0 && config.HardwareID != "" {
				return true
			}
		}
	}

	return false
}

func mergeConfig(local *pb.LocalConfig, incoming *pb.Config) {
	switch v := incoming.PayloadVariant.(type) {
	case *pb.Config_Device:
		local.Device = v.Device
	case *pb.Config_Position:
		local.Position = v.Position
	case *pb.Config_Power:
		local.Power = v.Power
	case *pb.Config_Network:
		local.Network = v.Network
	case *pb.Config_Display:
		local.Display = v.Display
	case *pb.Config_Lora:
		local.Lora = v.Lora
	case *pb.Config_Bluetooth:
		local.Bluetooth = v.Bluetooth
	case *pb.Config_Security:
		local.Security = v.Security
	}
}

func mergeModuleConfig(local *pb.LocalModuleConfig, incoming *pb.ModuleConfig) {
	switch v := incoming.PayloadVariant.(type) {
	case *pb.ModuleConfig_Mqtt:
		local.Mqtt = v.Mqtt
	case *pb.ModuleConfig_Serial:
		local.Serial = v.Serial
	case *pb.ModuleConfig_ExternalNotification:
		local.ExternalNotification = v.ExternalNotification
	case *pb.ModuleConfig_StoreForward:
		local.StoreForward = v.StoreForward
	case *pb.ModuleConfig_RangeTest:
		local.RangeTest = v.RangeTest
	case *pb.ModuleConfig_Telemetry:
		local.Telemetry = v.Telemetry
	case *pb.ModuleConfig_CannedMessage:
		local.CannedMessage = v.CannedMessage
	case *pb.ModuleConfig_Audio:
		local.Audio = v.Audio
	case *pb.ModuleConfig_RemoteHardware:
		local.RemoteHardware = v.RemoteHardware
	case *pb.ModuleConfig_NeighborInfo:
		local.NeighborInfo = v.NeighborInfo
	case *pb.ModuleConfig_AmbientLighting:
		local.AmbientLighting = v.AmbientLighting
	case *pb.ModuleConfig_DetectionSensor:
		local.DetectionSensor = v.DetectionSensor
	case *pb.ModuleConfig_Paxcounter:
		local.Paxcounter = v.Paxcounter
	}
}

func uploadConfig(config *DeviceConfig, adminURL, apiKey, meshID string, jsonOutput bool) {
	// Skip upload if API key or mesh ID not provided
	if apiKey == "" || meshID == "" {
		return
	}

	if !jsonOutput {
		fmt.Printf("\nUploading configuration to %s...\n", adminURL)
	}

	// Marshal config to JSON
	jsonData, err := json.Marshal(config)
	if err != nil {
		if !jsonOutput {
			fmt.Printf("❌ Failed to marshal config: %v\n", err)
		}
		return
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/api/meshes/%s/nodes/import", adminURL, meshID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		if !jsonOutput {
			fmt.Printf("❌ Failed to create request: %v\n", err)
		}
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		if !jsonOutput {
			fmt.Printf("❌ Failed to upload: %v\n", err)
		}
		return
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		if !jsonOutput {
			fmt.Printf("❌ Upload failed with status: %d\n", resp.StatusCode)
		}
		return
	}

	if !jsonOutput {
		fmt.Printf("✅ Configuration uploaded successfully\n")
	}
}
