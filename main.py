import time

from turbo import TCP015Controller

print("Connecting to turbo")
turbo = TCP015Controller("/dev/ttyUSB0")

#turbo.off()

print(f"{'R' if turbo.is_running() else 'Not r'}unning. Firmware {turbo.firmware_version()}, error code {turbo.error_code()}, RPM setpoint {turbo.rpm_setpoint()}")

if not turbo.is_running():
    print("Starting pump station in")
    for i in range(3, 0, -1):
        print(i)
        time.sleep(1)
    turbo.on()

print("Starting main loop.")
while True:
    print(f"{turbo.rpm()}RPM @ {turbo.current_draw()}A")
    time.sleep(1)
