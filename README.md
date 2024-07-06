# OpenReactor

Open-Source IEC nuclear fusion reactor control, monitoring, and data logging system



### Overview

OpenReactor is an open source reference design and control system for a small scale neutron-producing IEC fusor. The control system integrates with Pfeiffer high vacuum [turbo pumps](#high-vacuum-system), MKS and Edwards [high vacuum gauges](#vacuum-gauges), MKS and Sierra [mass flow controllers](#gas-regulation), [high voltage power supplies](#high-voltage-supply-controller), and [proportional neutron counters](#neutron-emission-detection) and NIM instrumentation.



### Features

- Repeatable "one-click" fusion - OpenReactor runs reusable YAML configuration profiles and monitors power, vacuum, and gas delivery conditions to achieve fusion with no user interaction

- Mobile Control - a mobile web app communicates with the reactor to tweak settings remotely in real time over the network
- `kubectl`-inspired CLI for monitoring and control
- Realtime data logging and monitoring with InfluxDB and Grafana
- Modular microservice architecture
- 100% open source



|               ![plasma](docs/img/photos/plasma.avif)               |                      ![plasma](docs/img/photos/lowenergy.avif)                      |                     ![full](docs/img/photos/full.avif)                     |
|:-----------------------------------------------------------------:|:----------------------------------------------------------------------------------:|:-------------------------------------------------------------------------:|
| **Deuterium plasma** <br> *(During neutron-producing fusion run)* | **Low-energy plasma test**    <br> *(Poor insulation on high voltage feedthrough)* | **Reactor hardware**   <br> *(Chamber behind lead shield and foil cover)* |



### Contents

- [Architecture](#Architecture)
  - [Hardware Subsystems](#Hardware-Subsystems)
  - [Deployment](#Deployment)
  - [State Machine](#State-Machine)
  - [Remote Control](#Remote-Control)
    - [Fusion Profiles](#Fusion-Profiles)
    - [`fusionctl`](#fusionctl)
  - [Data Logging](#Data-Logging)
- [Hardware](#Hardware)
  - [High Voltage Power Supply](#High-Voltage-Supply-Controller)
  - [Gas Conversion and Delivery System (GCDS)](#gas-conversion-and-delivery-system-gcds)
    - [D<sub>2</sub>O to D<sub>2</sub> Conversion](#d2o-to-d2-conversion)
    - [Gas Regulation](#Gas-Regulation)
  - [High Vacuum System](#high-vacuum-system)
    - [Turbo Pump Controller](#turbo-pump-controller)
    - [Vacuum Gauges](#vacuum-gauges)
      - [MKS](#MKS)
      - [Edwards](#Edwards)
  - [Neutron Emission Detection](#neutron-emission-detection)
    - [Adding a pulse output to scalar](#adding-a-pulse-output-to-the-ludlum-2000)
  - [Future Work](#future-work)



### Architecture

OpenReactor runs as a collection of microservices that interface with the hardware components of a reactor over RS232 or USB. Each runs in a container and exposes a REST API that's shared between similar hardware, such as different vacuum guages and mass flow controllers. Hardware subsystem microservices share a common [service](https://github.com/natesales/openreactor/tree/main/pkg/service) interface that manages configuration, polling, and data logging for each service.

#### Hardware Subsystems

- [`hv`](https://github.com/natesales/openreactor/tree/main/cmd/hv) - High voltage power supply
- [`pfturbo`](https://github.com/natesales/openreactor/tree/main/cmd/pfturbo) - Pfeiffer turbo molecular pump controllers
- [`edwgauge`](https://github.com/natesales/openreactor/tree/main/cmd/edwgauge) - Edwards high vacuum gauges
- [`mksgauge`](https://github.com/natesales/openreactor/tree/main/cmd/mksgauge) - MKS high vacuum gauges
- [`mksmfc`](https://github.com/natesales/openreactor/tree/main/cmd/mksmfc) - MKS mass flow controllers
- [`sierramfc`](https://github.com/natesales/openreactor/tree/main/cmd/sierramfc) - Sierra mass Flow controllers
- [`counter`](https://github.com/natesales/openreactor/tree/main/cmd/counter) - Neutron counter
- [`radiacode`](https://github.com/natesales/openreactor/tree/main/cmd/radiacode) - Radiacode scintillation counter


#### Deployment

The hardware subsystem services, state machine, and Grafana and Caddy, all run in Docker Compose on the reactor control computer. InfluxDB runs on a larger compute cluster to offload query processing.



#### State Machine

The central state machine and operator controls run in a service called ["`maestro`"](https://github.com/natesales/openreactor/tree/main/cmd/maestro) . When running a fusion profile, `maestro` monitors and adjusts reactor parameters such as voltage, vacuum level, and flow rate, to achieve repeatable fusion. It also serves the API and WebSocket server for the web UI and `fusionctl` control program.

![fsm](docs/img/diagrams/fsm.avif)

*FSM Flowchart*

`maestro` detects long term over-current and low vacuum errors and shuts the high voltage supply down to limit arcing. The high voltage supply controller also has an internal current limit that reacts much faster than maestro would be able to.



#### Remote Control

Reactor operators interact with OpenReactor through a web UI and `kubectl`-inspired CLI. The web UI communicates with `maestro` over WebSockets for reactor state management and manual control over the high voltage, gas delivery, and vacuum systems.

| ![ui](docs/img/ui.webp) | ![full](docs/img/fusionctl.webp) |
|:-----------------------:|:--------------------------------:|
|         Web UI          |          `fusionctl` CLI         |



#### Fusion Profiles

Fusion profiles describe the reactor's operational parameters for a fusion run. They're stored as YAML documents and managed with the  `fusionctl` control program.

#### `fusionctl`

For example, `fusionctl profile generate` creates a new default profile:

```yaml
name: Default # User specified name
revision: "0" # User specified revision identifier
auto:
    startOnApply: false # Start profile on apply
    startGas: false     # Skip to GasRegulating immediately on CathodeVoltageReached
vacuum:
    turboRotorSpeed: 90000
    turboRotorStartupHold: 15s
    lowVacuum: 0.01
    targetVacuum: 1e-05
    targetVacuumHold: 15s
cathode:
    tripCurrent: 8
    rampCurve: "0" # Voltage ramp curve in slope-intercept form in terms of x (in milliseconds)
    voltageCutoff: 40
gas:
    flowRate: 10
    flowSlop: 0.1
    runtime: 1m
```

To apply the profile, run `fusionctl profile apply -f 20240701001.yaml`. OpenReactor will start the profile at `TurboSpinup` if `auto.startOnApply` is set, otherwise it'll wait in the `ProfileReady` state for an operator to start the profile from the web UI or with `fusionctl fsm next`. Running `fusionctl fsm reset` or clicking the FSM reset button on the web UI sets the FSM back to `ProfileWait` and stops the current profile.




#### Data Logging

Each hardware subsystem reports metrics to a central InfluxDB server for visualization in Grafana.

![grafana](docs/img/grafana.webp)



## Hardware


### High Voltage Supply Controller

The [hv](https://github.com/natesales/openreactor/tree/main/cmd/hv) service controls and monitors a Spellman PTV power supply and communicates communicates with a [RP2040 over serial](https://github.com/natesales/openreactor/blob/main/cmd/hv/hv.ino). The microcontroller features an internal over-current shutoff, read from a 250W ballast resistor to sink momentary arc faults.

![hv](docs/img/diagrams/hv.avif)

*HV supply overview*

The power supply case is grounded to the chamber and mains earth through the AC plug, and a RG8 coax cable supplies the high voltage output to the cathode feedthrough.

![hv](docs/img/photos/hv.avif)

*High voltage power supply and ballast resistor*

### Gas Conversion and Delivery System (GCDS)

The gas delivery conversion and delivery system manages Deuterium production and regulation via electrolysis and flow restriction.

The system has two tasks:

1. Convert Deuterium Oxide (D<sub>2</sub>O) to Deuterium (D<sub>2</sub>) gas via electrolysis
2. Regulate gas flow into the central vacuum chamber

There are two independent gas supply lines connected to the vacuum chamber. The primary supply line handles gas conversion and closed-loop flow control, while a secondary manual syringe and needle valve allows the chamber to be purged with inert gas or short fusion runs when supplied with D<sub>2</sub> gas.

![gas](docs/img/diagrams/gas-delivery.avif)
*Gas Conversion and Delivery System Overview*

#### D<sub>2</sub>O to D<sub>2</sub> Conversion

D<sub>2</sub>O is manually injected into a PEM cell mounted under the mass flow controller as needed prior to reactor operation. A 3-way luer lock valve on the PEM cell input supply tube permits a syringe to push any remaining liquid into the PEM cell to conserve D<sub>2</sub>O. The PEM cell output is connected to a second 3-way luer lock "divert" valve fitted with a reservoir syringe to store D<sub>2</sub> gas during or prior to fusion operation. The divert valve enables the gas supply to come directly from the reservoir syringe or both the reservoir syringe and PEM cell for continuous operation.

#### Gas Regulation

The divert valve feeds the mass flow controller, which regulates gas flow from the reservoir syringe and, depending on the divert valve position, the PEM cell too.

OpenReactor supports [MKS](https://github.com/natesales/openreactor/tree/main/cmd/mksmfc) and [Sierra](https://github.com/natesales/openreactor/tree/main/cmd/sierramfc) mass flow controllers. Each shares an identical internal REST API and communicates with the MFC over RS232 (Sierra) or USB to a [RP2040-based control board](https://github.com/natesales/openreactor/blob/main/cmd/mksmfc/mksmfc.ino).

|           ![gas](docs/img/photos/gas.avif)           | ![Deuterium Oxide](docs/img/photos/deuterium.avif) |
|:---------------------------------------------------:|:-------------------------------------------------:|
| GDCS with MFC, PEM cell, and adapter column visible |           Deuterium Oxide (Heavy Water)           |

All high vacuum fittings are 1/4" VCR, with a series of reducers to adapt down to a luer lock syringe connector.



### High Vacuum System

The high vacuum system consists of a series of pumps and gauges to pull and monitor the high vacuum environment in the chamber.

![Vacuum System](docs/img/diagrams/vacuum.avif)

*High vacuum system overview*

#### Turbo Pump Controller

OpenReactor supports the [Pfeiffer Vacuum Protocol](https://mmrc.caltech.edu/Vacuum/Pfeiffer%20Turbo/Pfeiffer%20Interface%20RS@32.pdf) to control an monitor a wide range of Pfeiffer turbo pump controllers and drive units.

##### Adding a RS232 port to a Pfeiffer turbo pump controller

Many turbo pump controllers have a panel mount RS232 port, but some expose RS232 over their X5 port which is blocked by the turbo pump control cable. In this case, it's trivial to break out the RS232 TX, RX, and ground lines to a small panel mount jack between the X1 and X2 ports.

| ![pump](docs/img/photos/turbo-serial.avif) | ![pump](docs/img/photos/vacuum.avif) |
|:-----------------------------------------:|:-----------------------------------:|
| RS232 line added to the TCP control board | RS232 port via back panel connector |


#### Vacuum Gauges

OpenReactor supports [MKS](#MKS) and [Edwards](#Edwards) vacuum gauges for vacuum measurement from atmospheric to 1e-9 Torr.

##### MKS

MKS gauges connect to the [mksgauge](https://github.com/natesales/openreactor/tree/main/cmd/mksgauge) service over RS232 and report all available pressure reading slots (PR1/PR2 for 901p pirani and piezo gauges).

##### Edwards

Edwards gauges connect to a RP2040-based gauge controller which controls the gauge enable state and reports the analog gauge output over USB. The [edwgauge](https://github.com/natesales/openreactor/tree/main/cmd/edwgauge) service converts the analog gauge signal to a vacuum level according to an interpolation profile (currently supports Edwards AIM-S and APG-L).

Adding other gauges and gas-dependent curves is as easy as adding a new set of interpolation points to the [YAML config](https://github.com/natesales/openreactor/blob/main/cmd/edwgauge/gauge-lut.yml).



### Neutron Emission Detection

We detect neutron emissions using a proportional neutron counter tube, an amplifier, and a counter running on a RP2040. The [counter](https://github.com/natesales/openreactor/tree/main/cmd/counter) service logs the count rate over serial and supports any falling-edge signal from a NIM rack or scalar. A RadiaCode scintillation counter monitors lower energy background radiation over USB.

![neutron](docs/img/diagrams/neutron-detection.avif)
*Neutron detection overview*

#### Adding a pulse output to the Ludlum 2000

Older Ludlum scalars don't have a RS232 interface like the new ones do, so instead of wiring up a microcontroller to read and control the internal counter's time base and reset state, we can simply expose the pulse signal and trigger an interrupt on a microcontroller. We can wire the counter pulse trigger pin through a voltage divider to get a 3.3V falling-edge trigger signal, and pass it through a panel mount BNC jack to a RP2040 digital input pin.

|        ![pump](docs/img/photos/ludlum-tap.avif)         |    ![pump](docs/img/photos/ludlum-bnc.avif)     |
|:------------------------------------------------------:|:----------------------------------------------:|
| Tap on the count trigger pin on the 2000 control board | Panel mount BNC connector with voltage divider |



#### Future Work
- Plasma conditions
  - Langmuir probe
  - Optical spectrometer
- Neutron activation decay monitoring with RadiaCode




#### License

OpenReactor is licensed under GPL GNU v3.0.
