import math
import numpy
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


# List in range from 2 to 10 inclusive, in steps of 0.2
x = [
    2.00,
    2.50,
    3.00,
    3.20,
    3.40,
    3.60,
    3.80,
    4.00,
    4.20,
    4.40,
    4.60,
    4.80,
    5.00,
    5.20,
    5.40,
    5.60,
    5.80,
    6.00,
    6.20,
    6.40,
    6.60,
    6.80,
    7.00,
    7.20,
    7.40,
    7.60,
    7.80,
    8.00,
    8.20,
    8.40,
    8.60,
    8.80,
    9.00,
    9.20,
    9.40,
    9.60,
    9.80,
    9.90,
    10.00,
]

y = [
    "7.5E-9",
    "1.8E-8",
    "4.4E-8",
    "6.1E-8",
    "8.3E-8",
    "1.1E-7",
    "1.6E-7",
    "2.2E-7",
    "3.0E-7",
    "4.1E-7",
    "5.5E-7",
    "7.4E-7",
    "9.8E-7",
    "1.3E-6",
    "1.7E-6",
    "2.1E-6",
    "2.7E-6",
    "3.4E-6",
    "4.2E-6",
    "5.2E-6",
    "6.3E-6",
    "7.5E-6",
    "9.0E-6",
    "1.1E-5",
    "1.3E-5",
    "1.5E-5",
    "1.8E-5",
    "2.2E-5",
    "2.6E-5",
    "3.2E-5",
    "4.3E-5",
    "5.9E-5",
    "9.0E-5",
    "1.4E-4",
    "2.5E-4",
    "5.0E-4",
    "1.3E-3",
    "2.7E-3",
    "7.5E-3"
]
y = [float(i) for i in y]

# Logarithmic best fit


plt.scatter(x, y, label="Torr (from datasheet)")
plt.legend()
plt.show()
