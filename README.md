# meshmgr

A Meshtastic node manager written in Go.

## Overview

meshmgr is a command-line tool for managing Meshtastic nodes. It provides functionality to configure, monitor, and control Meshtastic devices.

## Project Structure

```
.
├── cmd/
│   └── meshmgr/        # Main application entry point
├── internal/           # Private application code
├── pkg/                # Public library code
├── go.mod              # Go module definition
└── README.md           # This file
```

## Installation

```bash
go build -o meshmgr ./cmd/meshmgr
```

## Usage

```bash
./meshmgr <command>
```

## Development Status

This project is in early development. Commands and features will be added progressively.

## License

TBD
