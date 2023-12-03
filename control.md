# Control System

## Analog inputs
- Pirani
- Low vac
- Peliter voltage (10:1 voltage divider)
- Grid voltage (5000:1 voltage divider)
- Grid current
- Heater current
- Total current (heater + HV primary + control electronics)
- Diff pump thermocouple
- Outside chamber thermocouple
- Ambient temp/humidity
- Hydrogen gas sensor

## Other data sources
- Cameras (viewport, external)
- Optical spectrometer (can I just use a camera?)
- Neutron counter
- xrays

## Warnings

- Reactor on and pressure too high
- Total current draw too high (breaker trip warning)
- D2 gas leak outside chamber
- xray radiation
- Temps (all thermocouples)

## Datalogging

- Something that can export to grafana, probably not prometheus since it isn't meant for duration-bounded time series

## Control box

- Raspberry pi
- Monitor to display reactor state and all measurements
- Big hardware kill switch for HV supply
- Lead shielding to protect controller
- USB webcams

### Outputs

- Mass flow controller (TODO: How do I calibrate it for deuterium? Does a hydrogen calibration work? Or does another element heavier than hydrgen work? like He?, or do I even need a mass flow controller if I can get away with a manually set needle valve and pulse the output with a solenoid?)
- Maybe HV supply variac?
- Plasma shaping coils (wrapped around sides of reactor)
- Coolant pump
- Relay bank for all AC power (Backing pump, HV supply, diffusion pump heater)
- Speaker for audible alerts

## Startup Procedure

Start...
- Datalogging
- Coolant pump
- Backing pump
- Wait until low vaccuum is reached, then start diff pump heater
- Wait for high vacuum
- Start HV supply
- Open mass flow controller

## Safeties
- Audible warning before HV supply starts from buzzer/speaker in control box
- Manual valve on deuterium tank
- Manual shutoff for HV supply
- Everything grounded to mains earth
- HV supply won't start unless there's a high vacuum?
- Lead shielding around viewport (3d printed enclosure for webcam with lead inserts)
- Control electronics on UPS with automatic shutdown sequence

## UIs

- Reactor monitor: interior camera, instantaneous sensor values, alerts
- Control dashboard
- Audience
    - SpaceX style "mission progress" bar at the bottom
        - Cooling and rough vacuum
        - High vacuum startup
        - High vacuum
        - HV Start
        - Deuterium flow control

## Control Interface

- Rugged laptop
- Port-specific VLAN for manual reactor control
- Fiber to reactor with media converters on either end for electrical isolation

## Public UI for running the reactor

- Startup/shutdown
- Power monitoring
- Limited mass flow control (of hydrogen instead of deuterium?)
- Voltage control

## Automatic Fusion Startup

Reads pressure, current draw, and optical spectrometry

Controls D2 flow, grid voltage, plasma shaping coils?

Multiple control loops:
- Constant chamber pressure (D2 flow)
- Overcurrent protection

Generate current/optical spectrum curves

## Processing

- Raspberry pi for web UI/datalogging
- Pi pico for closed loop control