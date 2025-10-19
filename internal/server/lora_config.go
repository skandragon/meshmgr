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

package server

import (
	"encoding/json"
	"net/http"
)

// LoRaConfig holds the frequency slot ranges for all region/preset combinations
type LoRaConfig struct {
	Regions []RegionInfo          `json:"regions"`
	Presets []PresetInfo          `json:"presets"`
	Slots   map[string]PresetSlots `json:"slots"` // key: region name
}

type RegionInfo struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type PresetInfo struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type PresetSlots struct {
	LongFast   int `json:"LongFast"`
	LongSlow   int `json:"LongSlow"`
	LongMod    int `json:"LongMod"`
	MediumFast int `json:"MediumFast"`
	MediumSlow int `json:"MediumSlow"`
	ShortFast  int `json:"ShortFast"`
	ShortSlow  int `json:"ShortSlow"`
	ShortTurbo int `json:"ShortTurbo"`
}

var loraConfig = LoRaConfig{
	Regions: []RegionInfo{
		{Code: "UNSET", Name: "Unset"},
		{Code: "US", Name: "US"},
		{Code: "EU_868", Name: "EU 868 MHz"},
		{Code: "EU_433", Name: "EU 433 MHz"},
		{Code: "CN", Name: "China"},
		{Code: "JP", Name: "Japan"},
		{Code: "ANZ", Name: "Australia/NZ"},
		{Code: "ANZ_433", Name: "Australia/NZ 433 MHz"},
		{Code: "KR", Name: "Korea"},
		{Code: "TW", Name: "Taiwan"},
		{Code: "RU", Name: "Russia"},
		{Code: "IN", Name: "India"},
		{Code: "NZ_865", Name: "New Zealand 865 MHz"},
		{Code: "TH", Name: "Thailand"},
		{Code: "UA_433", Name: "Ukraine 433 MHz"},
		{Code: "UA_868", Name: "Ukraine 868 MHz"},
		{Code: "MY_433", Name: "Malaysia 433 MHz"},
		{Code: "MY_919", Name: "Malaysia 919 MHz"},
		{Code: "SG_923", Name: "Singapore 923 MHz"},
		{Code: "KZ_433", Name: "Kazakhstan 433 MHz"},
		{Code: "KZ_863", Name: "Kazakhstan 863 MHz"},
		{Code: "BR_902", Name: "Brazil 902 MHz"},
		{Code: "PH_433", Name: "Philippines 433 MHz"},
		{Code: "PH_868", Name: "Philippines 868 MHz"},
		{Code: "PH_915", Name: "Philippines 915 MHz"},
		{Code: "NP_865", Name: "Nepal 865 MHz"},
		{Code: "LORA_24", Name: "LoRa 2.4 GHz"},
	},
	Presets: []PresetInfo{
		{Code: "ShortTurbo", Name: "Short Range / Turbo"},
		{Code: "ShortFast", Name: "Short Range / Fast"},
		{Code: "ShortSlow", Name: "Short Range / Slow"},
		{Code: "MediumFast", Name: "Medium Range / Fast"},
		{Code: "MediumSlow", Name: "Medium Range / Slow"},
		{Code: "LongFast", Name: "Long Range / Fast"},
		{Code: "LongMod", Name: "Long Range / Moderate"},
		{Code: "LongSlow", Name: "Long Range / Slow"},
	},
	Slots: map[string]PresetSlots{
		"UNSET":   {LongFast: 0, LongSlow: 0, LongMod: 0, MediumFast: 0, MediumSlow: 0, ShortFast: 0, ShortSlow: 0, ShortTurbo: 0},
		"ANZ":     {LongFast: 51, LongSlow: 103, LongMod: 103, MediumFast: 51, MediumSlow: 51, ShortFast: 51, ShortSlow: 51, ShortTurbo: 25},
		"ANZ_433": {LongFast: 5, LongSlow: 12, LongMod: 12, MediumFast: 5, MediumSlow: 5, ShortFast: 5, ShortSlow: 5, ShortTurbo: 2},
		"BR_902":  {LongFast: 21, LongSlow: 43, LongMod: 43, MediumFast: 21, MediumSlow: 21, ShortFast: 21, ShortSlow: 21, ShortTurbo: 10},
		"CN":      {LongFast: 159, LongSlow: 319, LongMod: 319, MediumFast: 159, MediumSlow: 159, ShortFast: 159, ShortSlow: 159, ShortTurbo: 79},
		"EU_433":  {LongFast: 3, LongSlow: 7, LongMod: 7, MediumFast: 3, MediumSlow: 3, ShortFast: 3, ShortSlow: 3, ShortTurbo: 1},
		"EU_868":  {LongFast: 0, LongSlow: 1, LongMod: 1, MediumFast: 0, MediumSlow: 0, ShortFast: 0, ShortSlow: 0, ShortTurbo: 0},
		"IN":      {LongFast: 7, LongSlow: 15, LongMod: 15, MediumFast: 7, MediumSlow: 7, ShortFast: 7, ShortSlow: 7, ShortTurbo: 3},
		"JP":      {LongFast: 11, LongSlow: 23, LongMod: 23, MediumFast: 11, MediumSlow: 11, ShortFast: 11, ShortSlow: 11, ShortTurbo: 5},
		"KR":      {LongFast: 11, LongSlow: 23, LongMod: 23, MediumFast: 11, MediumSlow: 11, ShortFast: 11, ShortSlow: 11, ShortTurbo: 5},
		"KZ_433":  {LongFast: 5, LongSlow: 12, LongMod: 12, MediumFast: 5, MediumSlow: 5, ShortFast: 5, ShortSlow: 5, ShortTurbo: 2},
		"KZ_863":  {LongFast: 19, LongSlow: 39, LongMod: 39, MediumFast: 19, MediumSlow: 19, ShortFast: 19, ShortSlow: 19, ShortTurbo: 9},
		"LORA_24": {LongFast: 101, LongSlow: 204, LongMod: 204, MediumFast: 101, MediumSlow: 101, ShortFast: 101, ShortSlow: 101, ShortTurbo: 50},
		"MY_433":  {LongFast: 7, LongSlow: 15, LongMod: 15, MediumFast: 7, MediumSlow: 7, ShortFast: 7, ShortSlow: 7, ShortTurbo: 3},
		"MY_919":  {LongFast: 19, LongSlow: 39, LongMod: 39, MediumFast: 19, MediumSlow: 19, ShortFast: 19, ShortSlow: 19, ShortTurbo: 9},
		"NP_865":  {LongFast: 11, LongSlow: 23, LongMod: 23, MediumFast: 11, MediumSlow: 11, ShortFast: 11, ShortSlow: 11, ShortTurbo: 5},
		"NZ_865":  {LongFast: 15, LongSlow: 31, LongMod: 31, MediumFast: 15, MediumSlow: 15, ShortFast: 15, ShortSlow: 15, ShortTurbo: 7},
		"PH_433":  {LongFast: 5, LongSlow: 12, LongMod: 12, MediumFast: 5, MediumSlow: 5, ShortFast: 5, ShortSlow: 5, ShortTurbo: 2},
		"PH_868":  {LongFast: 4, LongSlow: 10, LongMod: 10, MediumFast: 4, MediumSlow: 4, ShortFast: 4, ShortSlow: 4, ShortTurbo: 1},
		"PH_915":  {LongFast: 11, LongSlow: 23, LongMod: 23, MediumFast: 11, MediumSlow: 11, ShortFast: 11, ShortSlow: 11, ShortTurbo: 5},
		"RU":      {LongFast: 1, LongSlow: 3, LongMod: 3, MediumFast: 1, MediumSlow: 1, ShortFast: 1, ShortSlow: 1, ShortTurbo: 0},
		"SG_923":  {LongFast: 31, LongSlow: 63, LongMod: 63, MediumFast: 31, MediumSlow: 31, ShortFast: 31, ShortSlow: 31, ShortTurbo: 15},
		"TH":      {LongFast: 19, LongSlow: 39, LongMod: 39, MediumFast: 19, MediumSlow: 19, ShortFast: 19, ShortSlow: 19, ShortTurbo: 9},
		"TW":      {LongFast: 19, LongSlow: 39, LongMod: 39, MediumFast: 19, MediumSlow: 19, ShortFast: 19, ShortSlow: 19, ShortTurbo: 9},
		"UA_433":  {LongFast: 5, LongSlow: 12, LongMod: 12, MediumFast: 5, MediumSlow: 5, ShortFast: 5, ShortSlow: 5, ShortTurbo: 2},
		"UA_868":  {LongFast: 1, LongSlow: 3, LongMod: 3, MediumFast: 1, MediumSlow: 1, ShortFast: 1, ShortSlow: 1, ShortTurbo: 0},
		"US":      {LongFast: 103, LongSlow: 207, LongMod: 207, MediumFast: 103, MediumSlow: 103, ShortFast: 103, ShortSlow: 103, ShortTurbo: 51},
	},
}

// GetMaxSlot returns the maximum slot for a given region and preset
func GetMaxSlot(region, preset string) int {
	slots, ok := loraConfig.Slots[region]
	if !ok {
		return 319 // Default to max if unknown region
	}

	switch preset {
	case "LongFast":
		return slots.LongFast
	case "LongSlow":
		return slots.LongSlow
	case "LongMod":
		return slots.LongMod
	case "MediumFast":
		return slots.MediumFast
	case "MediumSlow":
		return slots.MediumSlow
	case "ShortFast":
		return slots.ShortFast
	case "ShortSlow":
		return slots.ShortSlow
	case "ShortTurbo":
		return slots.ShortTurbo
	default:
		return 319 // Default to max if unknown preset
	}
}

// handleGetLoRaConfig returns the LoRa configuration metadata
func (s *Server) handleGetLoRaConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(loraConfig); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
