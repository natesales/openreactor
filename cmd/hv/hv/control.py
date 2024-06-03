import pyvisa
import serial

ser = serial.Serial("/dev/ttyACM0", 115200)

rm = pyvisa.ResourceManager()
usb = filter(lambda x: 'USB' in x, rm.list_resources())
if len(usb) == 0:
    print("No USB devices found")
    exit(1)

logging.info(f"Connecting to {usb[0]}")
scope = rm.open_resource(usb[0])

ser.write(b's4500')
