# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**plzoreg** is a power regulator for photovoltaic (PV) domestic water heaters. The device sits between a PV inverter and an electric boiler, regulating power delivery to maintain the inverter's output voltage above a minimum threshold. This prevents the inverter from shutting down or operating inefficiently under heavy load.

The project combines:
- **Firmware** (TinyGo-based embedded code in `fw/`)
- **Hardware** (KiCad PCB design files in `board/`)

**Current Development Stage**: The firmware is in early development. The current implementation provides basic UI components (7-segment display, button input handling) that will be used for the power regulation control interface. The counter demo demonstrates the display and input system that will eventually show power levels and allow user control.

## Build System

All build commands must be run from the `fw/` directory.

### Building and Flashing
```bash
cd fw
make build    # Build firmware with TinyGo
make flash    # Build and flash to device via pyOCD
```

### Debugging
```bash
make gdb      # Start GDB server with semihosting
make rtt      # Connect to RTT debug output
```

### Analysis
```bash
make disassembly  # Generate disassembly.txt from built image
```

### Target Configuration
The Makefile defines the target board configuration:
- `TARGET_TINYGO`: TinyGo target (currently `embedfire-py32f030`)
- `TARGET_PYOCD`: pyOCD target (currently `py32f030x8`)
- `CMSIS_PACK`: CMSIS pack identifier (currently `PY32F030`)

## Architecture

### Firmware Structure (`fw/`)

The firmware is built with TinyGo using these compiler flags:
- `--scheduler tasks`: Task-based scheduler
- `--gc leaking`: No garbage collection (embedded context)
- `--serial rtt`: RTT for debug output

**Module Name**: `template` (as defined in `go.mod`)

#### Core Components

1. **`main.go`**: Application entry point
   - Initializes display with 3 digits at brightness 200
   - Creates keyboard with two buttons (PA11=increment, PA12=decrement)
   - Main loop: displays counter, handles button events
   - Holds both buttons for 3+ seconds to reset counter to 100

2. **`display/` package**: 7-segment LED display driver
   - `generic.go`: Platform-independent display API
     - `Init(digits, intensity)`: Initialize display
     - `NumberAt(number, leadingZero, dotPosition, basePosition)`: Display multi-digit number
     - `DigitAt(digit, dp, position)`: Display single digit with optional decimal point
     - `GlyphAt(glyph, position)`: Display raw 7-segment glyph
   - `py32.go`: PY32-specific peripheral configuration (build tag: `py32`)
     - Configures LED peripheral via memory-mapped registers
     - Direct hardware access using `device/py32` package
   - `py32f030.go`: GPIO pin mapping for PY32F030 (build tag: `py32f030`)
     - Maps 11 GPIO pins to LED segments/commons with alternate functions

3. **`keyboard/` package**: Button input handler
   - Polls GPIO pins configured as input with pullup
   - Generates three event types: `KeyDown`, `KeyPress` (initial + repeat), `KeyUp`
   - Key repeat: sends additional `KeyPress` events if held >500ms
   - Runs in goroutine, sends events via channel
   - Each key has an ID field used to identify increment vs decrement

#### Build Tags
The display driver uses Go build tags for hardware abstraction:
- `py32`: Common PY32 family code
- `py32f030`: Specific to PY32F030 variant

When modifying display code, ensure changes respect the build tag architecture:
- Generic display logic → `generic.go`
- PY32 peripheral setup → `py32.go` (tagged `//go:build py32`)
- Pin mappings → `py32f030.go` (tagged `//go:build py32f030`)

### Hardware (`board/`)

KiCad 8 project files for PCB design. The board interfaces with the PY32F030 microcontroller.

## Development Environment

The `.vscode/settings.json` configures the Go language server for TinyGo cross-compilation:
- Target: ARM Cortex-M (baremetal)
- GOROOT points to TinyGo's custom Go root
- Build tags include: `cortexm`, `baremetal`, `py32`, `py32f030`, and TinyGo-specific tags

VSCode default build task is `make flash`.
