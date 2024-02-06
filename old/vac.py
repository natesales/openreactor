import numpy
from gauge.highvac.aim_s import VOLTAGES, VACUUMS
from matplotlib import pyplot as plt


def voltage_divider(vout: float) -> float:
    # R2 is the resistor in parallel with the ADC input
    r1 = 10
    r2 = 10
    return (vout * (r1 + r2)) / r2


def vacuum() -> float:
    return v_to_torr(
        voltage_divider(
            adc()
        )
    )


def sci(f: float) -> str:
    return "{:.3e}".format(f)
