import binascii

import serial


gases = {
    1: "Air",
    2: "Argon",
    3: "CO2",
    4: "CO",
    5: "He",
    6: "H",
    7: "CH4",
    8: "N",
    9: "NO",
    10: "O",
}

def calc_crc(cmd: str):
    # cmd is a byte array containing the command ASCII string.
    # Example: cmnd="Sinv2.000"
    # an unsigned 32-bit integer is returned to the calling program
    # only the lower 16 bits contain the crc

    crc = 0xffff  # initialize crc to hex value 0xffff

    for character in cmd:  # this for loop starts with ASCII 'S' and loops through to the last ASCII '0'
        hex_char = int(ord(character))
        # hex_char = character
        crc = crc ^ (hex_char * 0x0100)  # the ASCII value is times by 0x0100 first then XORED to the current crc value
        # for(j=0; j<8; j++) # the crc is hashed 8 times with this for loop
        for j in range(0, 8):
            # if the 15th bit is set (tested by ANDING with hex 0x8000 and testing for 0x8000 result)
            # then crc is shifted left one bit (same as times 2) XORED with hex 0x1021 and ANDED to
            # hex 0xffff to limit the crc to lower 16 bits. If the 15th bit is not set then the crc
            # is shifted left one bit and ANDED with hex 0xffff to limit the crc to lower 16 bits.
            if (crc & 0x8000) == 0x8000:
                crc = ((crc << 1) ^ 0x1021) & 0xffff
            else:
                crc = (crc << 1) & 0xffff
            # end of j loop
        # end of i loop
    # There are some crc values that are not allowed, 0x00 and 0x0d

    # These are byte values so the high byte and the low byte of the crc must be checked and incremented if
    # the bytes are either 0x00 0r 0x0d
    if (crc & 0xff00) == 0x0d00:
        crc += 0x0100
    if (crc & 0x00ff) == 0x000d:
        crc += 0x0001
    if (crc & 0xff00) == 0x0000:
        crc += 0x0100
    if (crc & 0x00ff) == 0x0000:
        crc += 0x0001

    crc_hex_string = str(hex(crc))
    if len(crc_hex_string) < 6:
        crc_hex_string_final = crc_hex_string[:2] + '0' + crc_hex_string[2:]
    else:
        crc_hex_string_final = crc_hex_string
    first_byte = crc_hex_string_final[2:4]
    second_byte = crc_hex_string_final[4:6]
    final = binascii.unhexlify(first_byte + second_byte)

    return final

    # If the string Sinv2.000 is sent through this routine the crc = 0x8f55
    # The complete command "Sinv2.000" will look like this in hex:
    # 0x53 0x69 0x6E 0x76 0x32 0x2e 0x30 0x30 0x30 0x8f 0x55 0x0d


class MFC:
    def __init__(self, port):
        self.ser = serial.Serial(
            port,
            9600,
            parity=serial.PARITY_NONE,
            stopbits=serial.STOPBITS_ONE,
            timeout=3,
        )

    def is_connected(self):
        return self.cmd_controller("?Srnm").startswith("Srnm")

    def read_streaming_state(self):
        self.cmd_controller('?Strm')

    def set_streaming_state(self, mode):
        self.cmd_controller('!Strm' + mode)

    def read_gas(self):
        resp = self.cmd_controller('?Gasi')
        gas_id = int(resp[-1])
        if gas_id not in gases:
            return f"unknown ({gas_id})"
        return gases[gas_id]

    def set_gas(self, gas_index):
        if gas_index not in gases:
            raise ValueError("Invalid gas index")
        rsp = "Gasi" + str(gas_index)
        rsp = rsp + calc_crc(rsp) + '\x0d'
        return self.cmd_controller("!Gasi" + str(gas_index)) == rsp

    def get_setpoint(self):
        return float(self.cmd_controller("?Sinv").removeprefix("Sinv"))

    def set_setpoint(self, setpoint):
        rsp = "Sinv" + ('%.3f' % setpoint)
        rsp = rsp.encode() + calc_crc(rsp) + '\x0d'.encode()
        return self.cmd_controller("!Sinv" + ('%.3f' % setpoint)) == rsp

    def read_flow(self):
        return float(self.cmd_controller("?Flow").removeprefix("Flow"))

    def cmd_controller(self, cmd):
        print(f"Sending command: {cmd}")
        self.ser.write(cmd.encode() + calc_crc(cmd) + '\x0d'.encode())
        ser_rsp = self.ser.read_until(b'\r')
        response = ser_rsp[:-3].decode()
        crc = ser_rsp[-3:-1]
        if crc != calc_crc(response):
            print("CRC error")
        print(response)
        return response

mfc = MFC("/dev/ttyS1")
if mfc.is_connected():
    print("Connected to MFC")
# run various commands to test MFC
# mc_1.set_gas(8)

print(f"{mfc.read_flow()} of {mfc.read_gas()} setpoint {mfc.get_setpoint()}")

# mfc.set_setpoint(1.337)
