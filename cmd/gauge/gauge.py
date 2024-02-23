import time
import board
import busio
import supervisor

from digitalio import DigitalInOut, Direction

import adafruit_ads1x15.ads1015 as ADS
from adafruit_ads1x15.analog_in import AnalogIn

# Seconds to wait before checking gauge
GAUGE_CHECK_SEC = 5

i2c = busio.I2C(board.GP3, board.GP2)
ads = ADS.ADS1015(i2c)

gauge = DigitalInOut(board.GP6)
gauge.direction = Direction.OUTPUT
gauge.value = False

adc_c0 = AnalogIn(ads, ADS.P0)

def channel_voltage(chan):
    r1 = 39.2
    r2 = 15.02
    return round(chan.voltage / (r2 / (r1 + r2)), 4)

last_gauge_check = time.time()
while True:
    gauge_voltage = channel_voltage(adc_c0)
    print(f"{gauge_voltage};")

    if time.time() - last_gauge_check > GAUGE_CHECK_SEC:
        last_gauge_check = time.time()
        if not gauge.value:
            gauge.value = True
            time.sleep(1)
            gauge.value = channel_voltage(adc_c0) > 4 and channel_voltage(adc_c0) < 10

    time.sleep(0.1)
