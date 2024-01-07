import serial


def zero_pad(s, count: int) -> str:
    """
    Pad a string with zeroes to a given length
    """
    s = str(s)
    while len(s) != count:
        s = "0" + s
    return s


def cksum(data: str) -> str:
    """
    Calculate the checksum for a given string, returning the command with appended checksum
    """
    data = data.replace(" ", "")

    accum = 0
    for c in data:
        accum += ord(c)

    ck = str(accum % 256)
    while len(ck) != 3:
        ck = "0" + ck
    return ck


def cksum_valid(data: str) -> bool:
    """
    Check if a message's checksum is valid
    """
    body = data[:-3]
    ck = data[len(data) - 3:]
    # print(f"checking {data} body: {body} ck: {ck}")
    return cksum(body) == ck


class Message:
    def __init__(self, raw: str):
        self.raw = raw

        self.addr = int(self.raw[0:3])
        self.action = int(self.raw[3:5])
        self.param = int(self.raw[5:8])
        self.data_len = int(self.raw[8:10])
        self.payload = self.raw[10:10 + self.data_len]
        self.ck = int(self.raw[10 + self.data_len:])

        self.ck_valid: bool = cksum_valid(raw)


class PfeifferController:
    def __init__(self, port: str, addr: int = 1):
        self.ser = serial.Serial(
            port,
            9600,
            parity=serial.PARITY_NONE,
            stopbits=serial.STOPBITS_ONE,
        )
        self.addr = addr

        # while not self.ser.isOpen():
        #     print(f"Waiting for serial at {self.ser}...")
        #     time.sleep(1)

    def __del__(self):
        self.ser.close()

    def send(self, command: str):
        """
        Send a message, appending the checksum
        """
        self.ser.write(str(command + cksum(command) + '\r').encode())

    def read_register(self, reg: int) -> str:
        """
        Read a register by index
        """
        self.send(f"{zero_pad(self.addr, 3)}00{zero_pad(reg, 3)}02=?")
        msg = Message(self.ser.read_until(b'\r').decode().strip('\r'))
        if not msg.ck_valid:
            raise Exception(f"Invalid checksum while reading register {reg}")
        return msg.payload

    def write_register(self, reg: int, payload: str):
        """
        Write a register by index
        """
        body = zero_pad(self.addr, 3)
        body += "10"
        body += zero_pad(reg, 3)
        body += zero_pad(len(payload), 2)
        body += payload
        self.send(body)

    def set_register(self, reg: int, state: bool):
        """
        Set a boolean register's state
        """
        payload = "1" * 6
        if not state:
            payload = "0" * 6
        self.write_register(reg, payload)


class TCP015Controller(PfeifferController):
    def off(self):
        """
        Stop the station
        """
        self.set_register(10, True)

    def on(self):
        """
        Start the station
        """
        self.set_register(10, False)

    def is_running(self) -> bool:
        """
        Check if the station is running
        """
        return self.read_register(10) == "0" * 6

    def rpm(self) -> int:
        """
        Get the current pump RPM
        """
        # TODO: check format
        return int(self.read_register(309))

    def current_draw(self) -> int:
        """
        Get the pump current draw
        """
        # TODO: check format
        return int(self.read_register(310))

    def firmware_version(self) -> str:
        """
        Get the firmware version
        """
        return self.read_register(312).replace("  ", " ")

    def error_code(self) -> str:
        """
        Get the current error code
        """
        # TODO: Check this
        return self.read_register(303)

    def rpm_setpoint(self) -> int:
        """
        Get the current RPM setpoint
        """
        return int(self.read_register(308))
