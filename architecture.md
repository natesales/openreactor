Each hardware subsystem (HV grid power supply, vacuum system, gas delivery, etc.) has an associated microservice for control and monitoring. They communicate with hardware over serial (RS232 or USB) and with the central state machine and over REST calls. A Svelte web UI and `kubectl`-inspired CLI communicate with the reactor control system over WebSockets and REST RPCs to allow for remote control in real time.

Hardware subsystem microservices store their metrics in a central InfluxDB database and are displayed in Grafana.

![Grafana](docs/img/grafana.png)

|           ![UI](docs/img/ui.png)           |              foo              |
|:---------------------------------------------------:|:-----------------------------:|
| GDCS with MFC, PEM cell, and adapter column visible |      <pre>
$ fusionctl status
HV:
Setpoint: 0.00 v
Voltage: 0.00 kV
Vacuum:
Level: 1.00e-5 Torr
Turbo:
Running: false
Rotor Speed: 1502 Hz
Rotor Current: 1.81 A
Neutrons:
Count: 0 c/s
Gas:
Setpoint: 0 sccm
Flow: 0 sccm
</pre>       |





# TODO: Architecture diagram

- Prewriting

    - OpenReactor runs as a collection of microservices that interface with the hardware components of a reactor.

      Each hardware subsystem communicates with it's microservice for control and logging over RS232/USB TTY serial.

    - How humans interface with the system
        - `fusionctl`

    - How computers interface with the system
        - API

The `service` package manages configuration, polling, and control of each service.


## Data Acquisition

Subsystem microservices write metrics to a central InfluxDB database.

## FSM

The reactor

![fsm](docs/img/diagrams/fsm.jpg)


### Web UI


### `fusionctl` CLI

The `fusionctl` CLI provides a `kubectl` -inspired interface to control and monitor the reactor.


