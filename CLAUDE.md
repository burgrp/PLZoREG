# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**plzoreg** is a power regulator for photovoltaic (PV) domestic water heaters. The device sits between a PV inverter and an electric boiler, regulating power delivery to maintain the inverter's output voltage above a minimum threshold. This prevents the inverter from shutting down or operating inefficiently under heavy load.

The project combines:
- **Firmware** (TinyGo-based embedded code in `fw/`)
- **Hardware** (KiCad PCB design files in `board/`)

**Current Status**: The firmware implements a complete power regulation system with voltage monitoring, adaptive PWM control, temperature sensing, and a multi-page user interface for monitoring and configuration.

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

**Module Name**: `plzoreg` (as defined in `go.mod`)

#### Core Components

1. **`main.go`**: Application entry point and control logic
   - **System State Management**: Maintains voltage, temperature, duty cycle, and UI state
   - **Control Loop (500ms)**: Reads ADC values, implements regulation algorithm, updates PWM and display
   - **Regulation Algorithm**: Compares VSense with VTarget and adjusts duty cycle (±1% per cycle) to maintain voltage
   - **User Interface**: Multi-page display system with 5 information pages:
     - Page 1: Voltage Sense (no decimal point, e.g., `220`)
     - Page 2: Voltage Target (trailing decimal, e.g., `220.`) - adjustable
     - Page 3: Duty Cycle (decimal format, e.g., `0.75` = 75%)
     - Page 4: Triac heatsink temperature (° symbol, e.g., `50°`)
     - Page 5: MCU temperature (.° symbols, e.g., `50.°`)
   - **Settings Persistence**: Saves/loads VTarget from flash memory
   - **Safety Features**: Error detection for no sync (E01) and overtemperature >90°C (E02)

   **UI Navigation**:
   - Two buttons: Up (PA11) and Down (PA12)
   - Single press: cycle through pages
   - Hold both buttons 1s: enter/exit setting mode on VTarget and Duty pages
   - In setting mode: Up/Down buttons adjust values (VTarget: ±1V, Duty: ±5%)

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
   - Each key has an ID field (KeyUp=1, KeyDown=-1) used for navigation and value adjustment

4. **`flash/` package**: Flash memory operations
   - Provides page erase and program functions for persistent settings storage
   - VTarget value stored in last page of main flash with version key (FlashKeyV1)

5. **Hardware peripheral packages**:
   - `adc.go`: ADC configuration for voltage and temperature sensing
   - `pwm.go`: Timer-based PWM for phase-angle power control
   - `zcd.go`: Zero-cross detection using EXTI for AC synchronization

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
