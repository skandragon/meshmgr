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
	"fmt"
	"os"
)

func main() {
	fmt.Println("Meshtastic Node Manager")
	fmt.Println("Version: 0.1.0")

	if len(os.Args) > 1 {
		fmt.Printf("Command: %s\n", os.Args[1])
		fmt.Println("(Command handling to be implemented)")
	} else {
		fmt.Println("Usage: meshmgr <command>")
		fmt.Println("Commands will be implemented soon.")
	}
}
