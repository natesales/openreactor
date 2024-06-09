# OpenReactor

Open-Source Nuclear Fusion Reactor

## Architecture

Each hardware component of the reactor has an associated microservice to control and monitor that hardware subsystem.

- Prewriting
    - OpenReactor runs as a collection of microservices that interface with the hardware components of a reactor.
      serial, and either report their 
    - How components interface with eachother
        - Each hardware component has either a RS232 or USB TTY serial interface that
    - How humans interface with the system
        - `fusionctl`
    - How computers interface with the system
        - API

### High Voltage Supply Controller

### Mass Flow Controller

### Vacuum System

### Data Acquisition

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
