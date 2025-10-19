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
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	meshtastic_cli "github.com/skandragon/meshmgr/meshtastic-cli/proto"
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
	WANT_CONFIG_ID = 64 // Request config
)

type MeshtasticPacket struct {
	Magic1    byte
	Magic2    byte
	PacketLen uint16
	Data      []byte
}

type DeviceConfig struct {
	NodeNum      uint32 `json:"node_num"`
	HardwareID   string `json:"hardware_id"`
	LongName     string `json:"long_name"`
	ShortName    string `json:"short_name"`
	HasGPS       bool   `json:"has_gps"`
	Firmware     string `json:"firmware_version"`
	ConfigComplete bool `json:"config_complete"`
}

func main() {
	// Command line flags
	port := flag.String("port", "/dev/tty.usbmodem101", "Serial port device")
	baud := flag.Int("baud", 115200, "Baud rate")
	verbose := flag.Bool("v", false, "Verbose output")
	jsonOutput := flag.Bool("json", false, "Output as JSON")
	flag.Parse()

	if !*jsonOutput {
		fmt.Printf("Meshtastic Device Config Reader\n")
		fmt.Printf("================================\n")
		fmt.Printf("Connecting to %s at %d baud...\n", *port, *baud)
	}

	// Configure serial port
	config := &serial.Config{
		Name:        *port,
		Baud:        *baud,
		ReadTimeout: time.Millisecond * 100, // Short timeout for byte-by-byte reading
	}

	// Open serial port
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

	// Wake device and reset state machine (like Python does)
	if err := wakeDevice(s); err != nil {
		log.Fatalf("Failed to wake device: %v", err)
	}

	// Start packet reader in background
	packets := make(chan []byte, 100)
	debugLines := make(chan string, 100)
	go readPackets(s, packets, debugLines, *verbose)

	// Give device time to wake up
	time.Sleep(time.Millisecond * 100)

	// Request device info
	if !*jsonOutput {
		fmt.Println("\nRequesting device configuration...")
	}
	if err := requestConfig(s); err != nil {
		log.Fatalf("Failed to request config: %v", err)
	}

	// Collect device config
	deviceConfig := &DeviceConfig{}
	timeout := time.After(10 * time.Second) // Increased timeout

	for {
		select {
		case packet := <-packets:
			// Parse packet and update config
			if parseConfigPacket(packet, deviceConfig, *verbose && !*jsonOutput) {
				// Got config complete
				if *jsonOutput {
					jsonData, _ := json.Marshal(deviceConfig)
					fmt.Println(string(jsonData))
				} else {
					fmt.Println("\n✅ Device Configuration:")
					fmt.Printf("  Node Number: %d\n", deviceConfig.NodeNum)
					fmt.Printf("  Hardware ID: %s\n", deviceConfig.HardwareID)
					fmt.Printf("  Long Name: %s\n", deviceConfig.LongName)
					fmt.Printf("  Short Name: %s\n", deviceConfig.ShortName)
					fmt.Printf("  Has GPS: %v\n", deviceConfig.HasGPS)
					fmt.Printf("  Firmware: %s\n", deviceConfig.Firmware)
				}
				return
			}

		case line := <-debugLines:
			// Show debug output if verbose
			if *verbose && !*jsonOutput {
				fmt.Printf("  [DEBUG] %s\n", line)
			}

		case <-timeout:
			if !*jsonOutput {
				fmt.Println("\n⚠️  Timeout waiting for device configuration")
			}
			return
		}
	}
}

func wakeDevice(s *serial.Port) error {
	// Send START2 bytes to wake device and reset state machine
	// This mirrors what the Python library does
	wakeBytes := make([]byte, 32)
	for i := range wakeBytes {
		wakeBytes[i] = MAGIC_END // 0xC3
	}
	if _, err := s.Write(wakeBytes); err != nil {
		return fmt.Errorf("failed to wake device: %w", err)
	}
	return nil
}

func readPackets(s *serial.Port, packets chan []byte, debugLines chan string, verbose bool) {
	var rxBuf []byte
	var debugLine string

	for {
		// Read one byte at a time (like Python implementation)
		b := make([]byte, 1)
		n, err := s.Read(b)
		if err != nil || n != 1 {
			continue
		}

		c := b[0]
		ptr := len(rxBuf)

		// State machine for packet detection
		if ptr == 0 {
			// Looking for START1 (0x94)
			if c == MAGIC_START {
				rxBuf = append(rxBuf, c)
			} else {
				// Not a packet start - it's debug output
				handleDebugByte(c, &debugLine, debugLines)
			}
		} else if ptr == 1 {
			// Looking for START2 (0xC3)
			if c == MAGIC_END {
				rxBuf = append(rxBuf, c)
			} else {
				// Failed to match START2, treat previous START1 and this byte as debug
				handleDebugByte(rxBuf[0], &debugLine, debugLines)
				handleDebugByte(c, &debugLine, debugLines)
				rxBuf = nil
			}
		} else if ptr == 2 || ptr == 3 {
			// Reading length bytes
			rxBuf = append(rxBuf, c)
			if ptr == 3 {
				// Check packet length
				msb := uint16(rxBuf[2])
				lsb := uint16(rxBuf[3])
				packetLen := (msb << 8) | lsb

				if packetLen > MAX_TO_FROM_RADIO_SIZE {
					// Invalid length - treat as debug and reset
					for _, b := range rxBuf {
						handleDebugByte(b, &debugLine, debugLines)
					}
					rxBuf = nil
				}
			}
		} else {
			// Reading packet data
			rxBuf = append(rxBuf, c)

			// Check if packet is complete
			if len(rxBuf) >= 4 {
				msb := uint16(rxBuf[2])
				lsb := uint16(rxBuf[3])
				packetLen := (msb << 8) | lsb
				totalLen := 4 + int(packetLen)

				if len(rxBuf) >= totalLen {
					// Complete packet - send just the protobuf data (without header)
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
	// Convert byte to string (if possible)
	if c == '\r' {
		return // Ignore carriage returns
	} else if c == '\n' {
		// End of line - send it
		if len(*debugLine) > 0 {
			select {
			case debugLines <- *debugLine:
			default:
				// Don't block if channel is full
			}
			*debugLine = ""
		}
	} else if c >= 32 && c <= 126 {
		// Printable ASCII character
		*debugLine += string(c)
	}
}

func requestConfig(s *serial.Port) error {
	// Create ToRadio message with want_config_id
	toRadio := &meshtastic_cli.ToRadio{
		Payload: &meshtastic_cli.ToRadio_WantConfigId{
			WantConfigId: WANT_CONFIG_ID,
		},
	}

	// Serialize to protobuf
	data, err := protobuf.Marshal(toRadio)
	if err != nil {
		return fmt.Errorf("failed to marshal ToRadio: %w", err)
	}

	// Build packet with header
	packetLen := uint16(len(data))
	packet := make([]byte, 4+len(data))
	packet[0] = MAGIC_START
	packet[1] = MAGIC_END
	packet[2] = uint8(packetLen >> 8) // MSB
	packet[3] = uint8(packetLen & 0xFF) // LSB
	copy(packet[4:], data)

	if _, err = s.Write(packet); err != nil {
		return fmt.Errorf("failed to send config request: %w", err)
	}
	return nil
}

func parseConfigPacket(packet []byte, config *DeviceConfig, verbose bool) bool {
	// packet is now just the protobuf data without header
	// Try to unmarshal as FromRadio message
	fromRadio := &meshtastic_cli.FromRadio{}
	err := protobuf.Unmarshal(packet, fromRadio)
	if err != nil {
		return false
	}

	// Check payload type
	switch p := fromRadio.Payload.(type) {
	case *meshtastic_cli.FromRadio_MyInfo:
		if p.MyInfo != nil {
			config.NodeNum = p.MyInfo.MyNodeNum
			config.HasGPS = p.MyInfo.HasGps
			config.Firmware = p.MyInfo.FirmwareVersion
			if verbose {
				fmt.Println("  ✓ Got MyNodeInfo")
			}
		}

	case *meshtastic_cli.FromRadio_NodeInfo:
		// Only process our own node info (matches NodeNum from MyInfo)
		if p.NodeInfo != nil && p.NodeInfo.Num == config.NodeNum {
			if p.NodeInfo.User != nil {
				config.HardwareID = p.NodeInfo.User.Id
				config.LongName = p.NodeInfo.User.LongName
				config.ShortName = p.NodeInfo.User.ShortName
				config.ConfigComplete = true // Mark complete once we have the device node info
				if verbose {
					fmt.Println("  ✓ Got device NodeInfo")
				}
				// Check if we have all essential fields
				if config.NodeNum != 0 && config.HardwareID != "" {
					return true // Signal config complete
				}
			}
		}

	case *meshtastic_cli.FromRadio_ConfigCompleteId:
		if p.ConfigCompleteId == WANT_CONFIG_ID {
			config.ConfigComplete = true
			if verbose {
				fmt.Println("  ✓ Config complete signal received")
			}
			return true // Signal config complete
		}
	}

	return false
}

