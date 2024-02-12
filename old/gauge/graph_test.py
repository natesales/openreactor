import matplotlib.pyplot as plt
import numpy as np

from gauge.edwards import AIM_S, APG_L

for g in [AIM_S, APG_L]:
    plt.plot(g.x, g.y, 'o', label=g.name)
    x_interp = np.linspace(g.x[0], g.x[-1], np.size(g.x) * 10)
    plt.plot(x_interp, g.value_at(x_interp), label=g.name + " interp")

plt.title("Vacuum Gauge Calibration")
plt.xlabel("Voltage (V)")
plt.ylabel("Vacuum (Torr)")
plt.legend()
plt.show()
