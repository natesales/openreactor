import numpy as np
import matplotlib.pyplot as plt

plt.axis([0, 10, 0, 1])

xs = []
ys = []

for i in range(10):
    y = np.random.random()
   
    xs.append(i)
    ys.append(y)
    plt.clear()
    plt.plot(xs, ys)
    plt.pause(0.05)

plt.show()