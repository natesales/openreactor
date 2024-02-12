import numpy as np


class Gauge:
    def __init__(self, name: str, profile: dict):
        self.name = name
        self.x = np.array([float(i) for i in profile.keys()])
        self.y = np.array([float(i) for i in profile.values()])

    def value_at(self, v):
        # if v > self.x_max or v < self.x_min:
        #     return 0
        # TODO: Smoother curve
        return np.interp(v, self.x, self.y)
