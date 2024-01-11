import time

from influxdb_client.client.write.point import Point

import db
from turbo import TCP015Controller

print("Connecting to turbo")
turbo = TCP015Controller("/dev/ttyUSB0")

# turbo.off()

print(
    f"{'R' if turbo.is_running() else 'Not r'}unning. Firmware {turbo.firmware_version()}, error code {turbo.error_code()}, RPM setpoint {turbo.rpm_setpoint()}")

if not turbo.is_running():
    print("Starting pump station in")
    for i in range(3, 0, -1):
        print(i)
        time.sleep(1)
    turbo.on()

print("Starting data collection.")
while True:
    print(".", end="", flush=True)
    time.sleep(1)
    # db.write(Point("turbo_current_draw").field("amps", turbo.current_draw()))
    # db.write(Point("turbo_speed").field("hz", turbo.hz()))
    # time.sleep(0.25)
