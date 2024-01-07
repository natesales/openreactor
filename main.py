import time

from turbo import TCP015Controller

turbo = TCP015Controller("/dev/ttyUSB0")
# turbo.off()

print(f"Connected to turbo pump at {turbo.ser}, firmware version {turbo.firmware_version()}, error code {turbo.error_code()}, RPM setpoint {turbo.rpm_setpoint()}")

if not turbo.is_running():
    print("Starting pump station in")
    for i in range(3, 0, -1):
        print(i)
        time.sleep(1)
    turbo.on()

while True:
    print(f"RPM: {turbo.rpm()}, current draw: {turbo.current_draw()}")
    time.sleep(1)
